package models

// BatchList структура для сохранения нескольких ссылок
type BatchList []struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
	ShortURL      string `json:"short_url"`
}
