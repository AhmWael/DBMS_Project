package storage

import (
    "os"
    "strings"
    "fmt"
	"awsql/models"
)

// Create table: write header to disk
func CreateTableOnDisk(q models.CreateTableQuery) error {
    filename := "data/" + q.Table + ".txt"

    // If file exists, error
    if _, err := os.Stat(filename); err == nil {
        return fmt.Errorf("table already exists")
    }

    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    // First line = header
    _, err = f.WriteString(strings.Join(q.Columns, ",") + "\n")
    return err
}

// Insert row into table
func InsertIntoTable(q models.InsertQuery) error {
    filename := "data/" + q.Table + ".txt"

    f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer f.Close()

    _, err = f.WriteString(strings.Join(q.Values, ",") + "\n")
    return err
}

