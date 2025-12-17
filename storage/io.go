package storage

import (
    "strings"
    "fmt"
	"awsql/models"
    "awsql/storage/page"
    "awsql/storage/table"
)

// Create table on disk
func CreateTableOnDisk(q models.CreateTableQuery) error {
	path := "data/" + q.Table + ".tbl"

	pm, err := page.OpenPageManager(path)
	if err != nil {
		return err
	}

	headerPageID, _ := pm.AllocatePage() // allocate page 0 for header
	headerPage := page.NewPage(headerPageID)

	table.WriteHeader(headerPage.Data, q.Columns)

	return pm.WritePage(headerPage)
}

// Insert row into table
func InsertIntoTable(q models.InsertQuery) error {
	path := "data/" + q.Table + ".tbl"

	pm, err := page.OpenPageManager(path)
	if err != nil {
		return err
	}

	// TEMP: always insert into page 1
	dataPageID := page.PageID(1)

	p, err := pm.ReadPage(dataPageID)
	if err != nil {
		p = page.NewPage(dataPageID)
	}

	row := []byte(strings.Join(q.Values, ","))

	ok := table.InsertRow(p.Data, row)
	if !ok {
		return fmt.Errorf("page full (splitting not implemented yet)")
	}

	return pm.WritePage(p)
}

