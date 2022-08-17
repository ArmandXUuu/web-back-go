package structs

type Todo struct {
	ID      int    `db:"id" json:"id" ignore:"true"`
	Active  bool   `db:"active" json:"active" ignore:"true"`
	Done    bool   `db:"done" json:"done" ignore:"true"`
	Content string `db:"content" json:"content"`
}
