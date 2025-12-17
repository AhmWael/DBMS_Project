package table

import (
	"bytes"
	"encoding/binary"
)

// InsertRow inserts a row into the page data buffer
func InsertRow(p []byte, row []byte) bool {
	buf := bytes.NewBuffer(p)

	for buf.Len() < len(p) {
		if buf.Len()+2+len(row) > len(p) {
			return false
		}
		binary.Write(buf, binary.LittleEndian, uint16(len(row)))
		buf.Write(row)
	}
	return true
}
