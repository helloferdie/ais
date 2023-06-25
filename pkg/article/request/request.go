package request

// List -
type List struct {
	Page         int64  `query:"page" validate:"required,numeric,min=1"`
	ItemsPerPage int64  `query:"items_per_page" validate:"required,numeric,min=1,max=500"`
	Author       string `query:"author"`
	Query        string `query:"query"`
}

// View -
type View struct {
	ID int64 `param:"id" loc:"common." validate:"required"`
}

// Post -
type Post struct {
	Author string `json:"author" loc:"author" validate:"required"`
	Title  string `json:"title" loc:"title" validate:"required"`
	Body   string `json:"body" loc:"body" validate:"required"`
}

// Update -
type Update struct {
	ID int64 `json:"id" loc:"common." validate:"required"`
	Post
}

// Delete -
type Delete struct {
	ID int64 `json:"id" loc:"common." validate:"required"`
}
