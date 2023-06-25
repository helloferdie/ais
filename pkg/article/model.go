package article

import "github.com/helloferdie/golib/libdb"

// Model -
type Model struct {
	ID     int64  `db:"id" json:"id"`
	Author string `db:"author" json:"author"`
	Title  string `db:"title" json:"title"`
	Body   string `db:"body" json:"body"`
	libdb.ModelTimestamp
}
