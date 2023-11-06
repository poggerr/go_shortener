package handlers

import "fmt"

type BucketItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// MapToBucket создает корзину ссылок из `map[string]string`
func MapToBucket(baseURL string, m map[string]string) *[]BucketItem {
	bucket := make([]BucketItem, 0, len(m))
	for k, v := range m {
		bucket = append(bucket, BucketItem{
			ShortURL:    fmt.Sprintf("%s/%s", baseURL, k),
			OriginalURL: v,
		})
	}
	return &bucket
}
