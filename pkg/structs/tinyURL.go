package structs

type TinyURL struct {
	ID        int    `db:"id" json:"id"`
	BaseURL   string `db:"origin_url" json:"base_url"`
	ShortCode string `db:"md5" json:"short_code"`
}
