package models

// URL структура для обозначения короткой и длинной ссылки
type URL struct {
	LongURL  string `json:"url"`
	ShortURL string `json:"result"`
}

// Urls структура для получения урлов пользователя
type Urls struct {
	UserID      string `db:"user_id" json:"-"`
	LongURL     string `json:"original_url"`
	ShortURL    string `json:"short_url"`
	DeletedFlag bool   `db:"is_deleted" json:"-"`
}
