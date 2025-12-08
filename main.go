package main

import (
	"awsql/models"
	"awsql/storage"
	"bufio"
	"fmt"
	"net"
	"strings"
	"awsql/parser"
)

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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	buffer := "" // temporary storage for incomplete queries

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected")
			return
		}

		buffer += line

		// Check if query ends with semicolon
		if strings.ContainsRune(buffer, ';') {
			query := strings.TrimSpace(buffer)
			fmt.Println("Full query received:", query)

			result, err := parser.ParseSQL(query)
			if err != nil {
				conn.Write([]byte("ERR: " + err.Error() + "\n"))
				buffer = ""
				continue
			}

			switch t := result.(type) {
				case models.CreateTableQuery:
					err := storage.CreateTableOnDisk(t)
					if err != nil {
						fmt.Printf("Error creating table: %v\n", err)
						conn.Write([]byte("ERR: " + err.Error() + "\n"))
					} else {
						fmt.Printf("Table %s created with columns %v\n", t.Table, t.Columns)
						conn.Write([]byte("Table created\n"))
					}
				case models.InsertQuery:
					err := storage.InsertIntoTable(t)
					if err != nil {
						fmt.Printf("Error inserting into table: %v\n", err)
						conn.Write([]byte("ERR: " + err.Error() + "\n"))
					} else {
						fmt.Printf("Inserted into table %s values %v\n", t.Table, t.Values)
						conn.Write([]byte("Row inserted\n"))
					}
				case models.SelectQuery:
					conn.Write([]byte(fmt.Sprintf("PARSED: %+v\n", t)))
			}

			// reset buffer for next query
			buffer = ""
		}
	}
}


