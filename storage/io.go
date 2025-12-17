package storage

import (
    "strings"
    "fmt"
	"awsql/models"
    "awsql/storage/page"
    "awsql/storage/table"
    "encoding/binary"
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

	table.WriteHeader(headerPage.Data, q.Columns, []page.PageID{})

	return pm.WritePage(headerPage)
}

// InsertIntoTable inserts a row into the table, allocating pages if needed.
func InsertIntoTable(q models.InsertQuery) error {
    path := "data/" + q.Table + ".tbl"

    // Open page manager
    pm, err := page.OpenPageManager(path)
    if err != nil {
        return fmt.Errorf("failed to open page manager: %v", err)
    }

    // --- Read header page ---
    headerPage, err := pm.ReadPage(0)
    if err != nil {
        return fmt.Errorf("failed to read header page: %v", err)
    }

    columns, pageIDs, err := table.ReadHeader(headerPage.Data)
    if err != nil {
        return fmt.Errorf("failed to read table header: %v", err)
    }

    // --- Allocate first data page if none exist ---
    if len(pageIDs) == 0 {
        pid, err := pm.AllocatePage()
        if err != nil {
            return fmt.Errorf("failed to allocate first data page: %v", err)
        }
        pageIDs = append(pageIDs, pid)
        table.WriteHeader(headerPage.Data, columns, pageIDs)

        if err := pm.WritePage(headerPage); err != nil {
            return fmt.Errorf("failed to write header page: %v", err)
        }
    }

    // Prepare row data
    row := []byte(strings.Join(q.Values, ","))

    // Insert into last page
    lastPageID := pageIDs[len(pageIDs)-1]
    p, err := pm.ReadPage(lastPageID)
    if err != nil {
        // If page missing, create new
        p = page.NewPage(lastPageID)
    }

    if ok := table.InsertRow(p.Data, row); ok {
        // Row fits in current page
        if err := pm.WritePage(p); err != nil {
            return fmt.Errorf("failed to write page %d: %v", p.ID, err)
        }
        return nil
    }

    // Current page full, allocate new page
    newPID, err := pm.AllocatePage()
    if err != nil {
        return fmt.Errorf("failed to allocate new page: %v", err)
    }
    newPage := page.NewPage(newPID)

    // Insert row into new page
    if ok := table.InsertRow(newPage.Data, row); !ok {
        return fmt.Errorf("row too large to fit in empty page")
    }

    // Write new page
    if err := pm.WritePage(newPage); err != nil {
        return fmt.Errorf("failed to write new page %d: %v", newPID, err)
    }

    // Update linked list
    table.SetNextPage(p.Data, newPID) // link old last page -> new page
    if err := pm.WritePage(p); err != nil {
        return fmt.Errorf("failed to update last page link: %v", lastPageID)
    }

    // Update header to include new page
    pageIDs = append(pageIDs, newPID)
    table.WriteHeader(headerPage.Data, columns, pageIDs)

    if err := pm.WritePage(headerPage); err != nil {
        return fmt.Errorf("failed to write updated header page: %v", err)
    }

    return nil
}

// SelectFromTable selects rows from the table based on the query
func SelectFromTable(q models.SelectQuery) ([][]string, error) {
    path := "data/" + q.Table + ".tbl"

    pm, err := page.OpenPageManager(path)
    if err != nil {
        return nil, err
    }

    // Read header
    headerPage, _ := pm.ReadPage(0)
    columns, pageIDs, _ := table.ReadHeader(headerPage.Data)

    // Determine which columns to return
    var colIndexes []int
    if len(q.Columns) == 1 && q.Columns[0] == "*" {
        // select all
        for i := range columns {
            colIndexes = append(colIndexes, i)
        }
    } else {
        for _, c := range q.Columns {
            found := false
            for i, hc := range columns {
                if c == hc {
                    colIndexes = append(colIndexes, i)
                    found = true
                    break
                }
            }
            if !found {
                return nil, fmt.Errorf("column %s not found", c)
            }
        }
    }

    var result [][]string

    // Loop through pages
    for _, pid := range pageIDs {
        p, _ := pm.ReadPage(pid)
        rows := readRowsFromPage(p.Data)
        for _, r := range rows {
            if q.Where != nil {
                if !matchesCondition(r, columns, q.Where) {
                    continue
                }
            }
            var selected []string
            for _, idx := range colIndexes {
                selected = append(selected, r[idx])
            }
            result = append(result, selected)
        }
        // Follow next page pointers
        next := table.GetNextPage(p.Data)
        for next != 0 {
            p, _ = pm.ReadPage(next)
            rows := readRowsFromPage(p.Data)
            for _, r := range rows {
                if q.Where != nil {
                    if !matchesCondition(r, columns, q.Where) {
                        continue
                    }
                }
                var selected []string
                for _, idx := range colIndexes {
                    selected = append(selected, r[idx])
                }
                result = append(result, selected)
            }
            next = table.GetNextPage(p.Data)
        }
    }

    return result, nil
}

// readRowsFromPage reads all rows from the given page data
func readRowsFromPage(data []byte) [][]string {
    var rows [][]string
    used := binary.LittleEndian.Uint16(data[table.PageUsedOffset : table.PageUsedOffset+2])
    offset := table.PageDataOffset

    // Read rows until used bytes
    for offset < int(table.PageDataOffset+used) {
        // if not enough bytes for length prefix, break
        if offset+2 > len(data) {
            break
        }
        rowLen := binary.LittleEndian.Uint16(data[offset:])
        offset += 2
        // if not enough bytes for row data, break
        if offset+int(rowLen) > len(data) {
            break
        }
        rowData := string(data[offset : offset+int(rowLen)])
        rows = append(rows, strings.Split(rowData, ","))
        offset += int(rowLen)
    }
    return rows
}

// matchesCondition checks if a row matches the given condition
func matchesCondition(row []string, columns []string, cond *models.Condition) bool {
    for i, c := range columns {
        if c == cond.Left {
            switch cond.Operator {
            case "=":
                return row[i] == cond.Right
            case "!=":
                return row[i] != cond.Right
            case ">":
                return row[i] > cond.Right
            case "<":
                return row[i] < cond.Right
            }
        }
    }
    return false
}


