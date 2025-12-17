package table

import (
	"bytes"
	"encoding/binary"
	"awsql/storage/page"
)

// TableHeader represents the header of a table
type TableHeader struct {
	NumColumns uint16
	PageIDs   []page.PageID // list of data pages for the table
}

// WriteHeader writes the table header to the given buffer
func WriteHeader(buf []byte, columns []string, pageIDs []page.PageID) {
	// Use a bytes buffer
	b := bytes.NewBuffer(buf[:0])

	// Write number of columns
	binary.Write(b, binary.LittleEndian, uint16(len(columns)))
	// Write column names
	for _, c := range columns {
		binary.Write(b, binary.LittleEndian, uint16(len(c)))
		b.WriteString(c)
	}

	// Write page IDs
	binary.Write(b, binary.LittleEndian, uint16(len(pageIDs)))
	// Write each page ID
    for _, pid := range pageIDs {
        binary.Write(b, binary.LittleEndian, uint64(pid))
    }
}

// ReadHeader reads the table header from the given buffer
func ReadHeader(buf []byte) (columns []string, pageIDs []page.PageID, err error) {
    r := bytes.NewReader(buf)
    var numCols uint16
    if err := binary.Read(r, binary.LittleEndian, &numCols); err != nil {
        return nil, nil, err
    }

    columns = make([]string, numCols)
    for i := 0; i < int(numCols); i++ {
        var colLen uint16
        binary.Read(r, binary.LittleEndian, &colLen)
        colBytes := make([]byte, colLen)
        r.Read(colBytes)
        columns[i] = string(colBytes)
    }

    var numPages uint16
    binary.Read(r, binary.LittleEndian, &numPages)
    pageIDs = make([]page.PageID, numPages)
    for i := 0; i < int(numPages); i++ {
        var pid uint64
        binary.Read(r, binary.LittleEndian, &pid)
        pageIDs[i] = page.PageID(pid)
    }

    return columns, pageIDs, nil
}