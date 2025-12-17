package page

import (
	"os"
	"io"
)

// PageManager manages low-level page I/O to a file
type PageManager struct {
	file *os.File
}

// OpenPageManager opens a page manager for the given file path
func OpenPageManager(path string) (*PageManager, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return &PageManager{file: f}, nil
}

// ReadPage reads a page by its ID
func (pm *PageManager) ReadPage(id PageID) (*Page, error) {
	p := NewPage(id)
	offset := int64(id) * PageSize

	n, err := pm.file.ReadAt(p.Data, offset) // read bytes from file into page

	// If we hit EOF, page is new/empty
	if err != nil && err != io.EOF {
		return nil, err
	}

	// If fewer bytes read than PageSize (EOF or partial), zero-fill remaining bytes
	if n < len(p.Data) {
		for i := n; i < len(p.Data); i++ {
			p.Data[i] = 0
		}
	}

	return p, nil
}

// WritePage writes a page to disk
func (pm *PageManager) WritePage(p *Page) error {
	offset := int64(p.ID) * PageSize
	_, err := pm.file.WriteAt(p.Data, offset)
	return err
}

// AllocatePage allocates a new page and returns its ID
func (pm *PageManager) AllocatePage() (PageID, error) {
	info, err := pm.file.Stat()
	if err != nil {
		return 0, err
	}
	return PageID(info.Size() / PageSize), nil
}
