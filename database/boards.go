package database

// Board is single dashboard table/column.
type Board struct {
	Title string `json:"title"`
	Query string `json:"query"`
}
