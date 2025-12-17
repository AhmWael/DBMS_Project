package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync/atomic"
	"strconv"

	"awsql/config"
	"awsql/models"
	"awsql/storage"
	"awsql/parser"
	"awsql/wal"
)

var txCounter uint64

// nextTxID generates a new transaction ID
func nextTxID() string {
	id := atomic.AddUint64(&txCounter, 1)
	return strconv.FormatUint(id, 10)
}


func main() {
	// Listen on TCP port 8888
	ln, err := net.Listen("tcp", ":8888")
	if err != nil {
		panic(err)
	}
	fmt.Println("DBMS server running on port 8888...")

	for {
		// Accept a connection
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

// handleConnection processes a client connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Greeting message
	conn.Write([]byte(fmt.Sprintf(
		"%s (%s)\r\nType \"help\" for help.\r\nawsql> ",
		config.ServerName,
		config.ServerVersion,
	)))

	reader := bufio.NewReader(conn)
	buffer := "" // temporary storage for incomplete queries

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected")
			return
		}

		buffer += line
		trimmed := strings.TrimSpace(buffer)
		if trimmed == "" {
			conn.Write([]byte("awsql> "))
			continue
		}

		// Handle help command
		if strings.EqualFold(trimmed, "help") {
			helpText := "" +
				"Supported commands:\r\n" +
				"  help\r\n" +
				"  exit\r\n" +
				"  CREATE TABLE table_name (col1, col2, ...);\r\n" +
				"  INSERT INTO table_name VALUES (val1, val2, ...);\r\n" +
				"  SELECT col1, col2 FROM table_name WHERE condition;\r\n"
			conn.Write([]byte(helpText + "\r\nawsql> "))
			buffer = "" // reset buffer
			continue
		}

		// Handle exit command
		if strings.EqualFold(trimmed, "exit") {
			conn.Write([]byte("Goodbye!\r\n"))
			return
		}


		// Check if query ends with semicolon
		if strings.ContainsRune(buffer, ';') {
			query := strings.TrimSpace(buffer)
			fmt.Println("Full query received:", query)

			result, err := parser.ParseSQL(query)
			if err != nil {
				conn.Write([]byte("ERR: " + err.Error() + "\r\n"))
				buffer = ""
				continue
			}

			switch t := result.(type) {
				case models.CreateTableQuery:
					txID := nextTxID()

					wal.LogBegin(txID)
					wal.LogCreate(txID, t.Table, strings.Join(t.Columns, ","))

					err := storage.CreateTableOnDisk(t)
					if err != nil {
						conn.Write([]byte("ERR: " + err.Error() + "\r\n"))
						return
					}

					wal.LogCommit(txID)
					conn.Write([]byte("Table created\r\n"))
					fmt.Printf("Table %s created with columns %v\r\n", t.Table, t.Columns)
				case models.InsertQuery:
					txID := nextTxID()

					wal.LogBegin(txID)
					wal.LogInsert(txID, t.Table, strings.Join(t.Values, ","))

					err := storage.InsertIntoTable(t)
					if err != nil {
						conn.Write([]byte("ERR: " + err.Error() + "\r\n"))
						return
					}

					wal.LogCommit(txID)
					conn.Write([]byte("Row inserted\r\n"))
					fmt.Printf("Inserted into %s values %v\r\n", t.Table, t.Values)
				case models.SelectQuery:
					rows, err := storage.SelectFromTable(t)
					if err != nil {
						conn.Write([]byte("ERR: " + err.Error() + "\r\n"))
						break
					}
					for _, r := range rows {
						conn.Write([]byte(strings.Join(r, ",") + "\r\n"))
					}
			}

			// reset buffer for next query
			buffer = ""
			conn.Write([]byte("awsql> "))
		}
	}
}


