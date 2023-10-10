// Package models содержит все модели проекта
package models

import (
	"github.com/google/uuid"
)

// BatchList структура для сохранения нескольких ссылок
type BatchList []struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
	ShortURL      string `json:"short_url"`
}

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

// Keys структура для загрузки коротких ссылок для удаления в бд
type Keys []struct {
	Key string `json:"key"`
}
