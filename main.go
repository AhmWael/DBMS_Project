package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type SelectQuery struct {
	Table string
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
				conn.Write([]byte("ERR: " + err.Error() + "\n"))
			} else {
				conn.Write([]byte(fmt.Sprintf("PARSED: %+v\n", result)))
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
	// Example: SELECT * FROM users
	parts := strings.Fields(query) // splits by spaces
	for _, tok := range parts {
		fmt.Println("\"", tok, "\"")
	}

	// Expecting: ["SELECT", "*", "FROM", "users"]
	if len(parts) != 4 {
		return SelectQuery{}, nil // basic error handling for now
	}

	return SelectQuery{
		Table: parts[3],
	}, nil
}
