package models

import "github.com/google/uuid"

// URL структура для обозначения короткой и длинной ссылки
type URL struct {
	LongURL  string `json:"url"`
	ShortURL string `json:"result"`
}

// User структура для обозначения базового пользователя в проекте
type User struct {
	ID       *uuid.UUID `json:"id"`
	UserName string     `json:"username"`
	Pass     string     `json:"pass"`
}

// Urls структура для получения урлов пользователя
type Urls struct {
	UserID      string `db:"user_id"`
	LongURL     string `json:"original_url"`
	ShortURL    string `json:"short_url"`
	DeletedFlag bool   `db:"is_deleted"`
}

// Storage массив структур Urls для получения списка ссылок пользователя
type Storage []Urls
