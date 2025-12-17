package page

type PageID uint64

type Page struct {
	ID   PageID
	Data []byte // always len == PageSize
}

// NewPage creates a new page with the given ID
func NewPage(id PageID) *Page {
	return &Page{
		ID:   id,
		Data: make([]byte, PageSize),
	}
}
