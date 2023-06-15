package models

type BatchList []struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
	ShortURL      string `json:"short_url"`
}

type URL struct {
	LongURL  string `json:"url"`
	ShortURL string `json:"result"`
}
