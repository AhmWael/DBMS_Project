package table

import (
	"encoding/binary"
	"awsql/storage/page"
)

const NextPageOffset = 0
const PageUsedOffset = 8   // 2 bytes for used size
const PageDataOffset = 10  // row data starts after NextPageID + UsedBytes

// InsertRow inserts a row into the page data buffer
func InsertRow(p []byte, row []byte) bool {
    used := binary.LittleEndian.Uint16(p[PageUsedOffset:PageUsedOffset+2])

	// Check if there is enough space
    if int(used)+2+len(row) > len(p)-PageDataOffset {
        return false
    }

    // Write the row length
    binary.LittleEndian.PutUint16(p[PageDataOffset+used:], uint16(len(row)))
    copy(p[PageDataOffset+used+2:], row)

    // Update used size
    binary.LittleEndian.PutUint16(p[PageUsedOffset:PageUsedOffset+2], used+2+uint16(len(row)))
    return true
}

// SetNextPage sets the next page ID in the page header
func SetNextPage(p []byte, next page.PageID) {
    binary.LittleEndian.PutUint64(p[NextPageOffset:], uint64(next))
}

// GetNextPage gets the next page ID from the page header
func GetNextPage(p []byte) page.PageID {
    return page.PageID(binary.LittleEndian.Uint64(p[NextPageOffset:]))
}

