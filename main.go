package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Condition struct {
	Left     string
	Operator string
	Right    string
}

type SelectQuery struct {
	Columns []string
	Table   string
	Where   *Condition
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

			result, err := ParseSQL(query)
			if err != nil {
				panic(err)
				//conn.Write([]byte("ERR: " + err.Error() + "\n"))
			} else {
				fmt.Printf("PARSED: %+v\n", result)
				//conn.Write([]byte(fmt.Sprintf("PARSED: %+v\n", result)))
			}

			// reset buffer for next query
			buffer = ""
		}
	}
}

func ParseSQL(query string) (interface{}, error) {
	query = strings.TrimSpace(query)
	query = strings.TrimSuffix(query, ";")

	// Make everything uppercase for keyword matching
	upper := strings.ToUpper(query)

	if strings.HasPrefix(upper, "SELECT") {
		return parseSelect(query)
	}

	return nil, nil
}

func parseSelect(query string) (SelectQuery, error) {
	query = strings.TrimSpace(query)
	query = strings.TrimSuffix(query, ";")

	// Split by spaces first
	parts := strings.Fields(query)

	// Must start with SELECT
	if len(parts) < 4 || strings.ToUpper(parts[0]) != "SELECT" {
		return SelectQuery{}, fmt.Errorf("invalid SELECT syntax")
	}

	// Find FROM index
	fromIdx := -1
	for i, tok := range parts {
		if strings.ToUpper(tok) == "FROM" {
			fromIdx = i
			break
		}
	}

	if fromIdx == -1 || fromIdx == 1 {
		return SelectQuery{}, fmt.Errorf("missing FROM clause")
	}

	// Columns are everything between SELECT and FROM
	columnsStr := strings.Join(parts[1:fromIdx], " ")
	columns := strings.Split(columnsStr, ",")
	for i := range columns {
		columns[i] = strings.TrimSpace(columns[i])
	}

	// Table name is next token after FROM
	table := parts[fromIdx+1]

	// Check if WHERE clause exists
	var cond *Condition
	whereIdx := -1
	for i, tok := range parts {
		if strings.ToUpper(tok) == "WHERE" {
			whereIdx = i
			break
		}
	}

	if whereIdx != -1 {
		// Simple condition: left operator right
		if len(parts) < whereIdx+4 {
			return SelectQuery{}, fmt.Errorf("invalid WHERE clause")
		}
		cond = &Condition{
			Left:     parts[whereIdx+1],
			Operator: parts[whereIdx+2],
			Right:    parts[whereIdx+3],
		}
	}

	return SelectQuery{
		Columns: columns,
		Table:   table,
		Where:   cond,
	}, nil
}
