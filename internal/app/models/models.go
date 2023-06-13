package models

type BatchList []struct {
	CorrelationId string `json:"correlation_id"`
	OriginalUrl   string `json:"original_url"`
	ShortUrl      string `json:"short_url"`
}

//type BatchListRes []struct {
//	CorrelationId string `json:"correlation_id"`
//	ShortUrl      string `json:"short_url"`
//}

type URL struct {
	LongURL  string `json:"url"`
	ShortURL string `json:"result"`
}
