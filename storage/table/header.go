package table

import (
	"bytes"
	"encoding/binary"
)

type TableHeader struct {
	NumColumns uint16
}

// WriteHeader writes the table header to the given buffer
func WriteHeader(buf []byte, columns []string) {
	b := bytes.NewBuffer(buf[:0])

	binary.Write(b, binary.LittleEndian, uint16(len(columns)))
	for _, c := range columns {
		binary.Write(b, binary.LittleEndian, uint16(len(c)))
		b.WriteString(c)
	}
}
