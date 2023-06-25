package article

import (
	"strings"

	"github.com/helloferdie/golib/libdb"
	"github.com/helloferdie/golib/libslice"
)

// DBMode -
var DBMode = libdb.Mode{
	Skip:          []string{"created_at", "updated_at"},
	AutoTimestamp: true,
}

// TConfig -
var TConfig = libdb.Config{
	Table:      "article",
	Module:     "article",
	Fields:     strings.Join(libslice.GetTagSlice(Model{}, "db"), ","),
	SoftDelete: true,
}
