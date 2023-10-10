// Package models содержит все модели проекта
package models

// Keys структура для загрузки коротких ссылок для удаления в бд
type Keys []struct {
	Key string `json:"key"`
}
