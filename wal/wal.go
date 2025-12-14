package wal

import (
	"os"
	"sync"
)

const WAL_FILE = "data/wal.log"

var mu sync.Mutex

func appendRecord(record string) error {
	mu.Lock()
	defer mu.Unlock()

	f, err := os.OpenFile(WAL_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(record + "\n")
	if err != nil {
		return err
	}

	return f.Sync() // durability guarantee
}

// LogBegin logs the beginning of a transaction
func LogBegin(txID string) error {
	return appendRecord("BEGIN|" + txID)
}

// LogCreate logs a CREATE TABLE operation
func LogCreate(txID, table, cols string) error {
	return appendRecord("CREATE|" + table + "|" + cols)
}

// LogInsert logs an INSERT operation
func LogInsert(txID, table, values string) error {
	return appendRecord("INSERT|" + table + "|" + values)
}

// LogCommit logs the commit of a transaction
func LogCommit(txID string) error {
	return appendRecord("COMMIT|" + txID)
}

