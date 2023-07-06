package models

import (
	"github.com/google/uuid"
)

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
	ID       *uuid.UUID `json:"id"`
	UserName string     `json:"username"`
	Pass     string     `json:"pass"`
}

type Urls struct {
	UserID      string `db:"user_id"`
	LongURL     string `json:"original_url"`
	ShortURL    string `json:"short_url"`
	DeletedFlag bool   `db:"is_deleted"`
}

type Keys []struct {
	Key string `json:"key"`
}

type Storage []Urls
