package models

import "github.com/google/uuid"

type BatchList []struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
	ShortURL      string `json:"short_url"`
}

type URL struct {
	LongURL  string `json:"url"`
	ShortURL string `json:"result"`
}

type User struct {
	Id       *uuid.UUID `json:"id"`
	UserName string     `json:"username"`
	Pass     string     `json:"pass"`
}

type Urls struct {
	UserId      string `db:"user_id"`
	LongURL     string `json:"original_url"`
	ShortURL    string `json:"short_url"`
	DeletedFlag bool   `db:"is_deleted"`
}

type Storage []Urls
