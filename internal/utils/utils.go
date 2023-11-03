package utils

import (
	"context"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"net/url"
	"os"
)

func NewShortURL() string {
	return gonanoid.Must(8)
}

func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func CheckFilename(filename string) (err error) {
	// Check if file already exists
	if _, err = os.Stat(filename); err == nil {
		return nil
	}

	// Attempt to create it
	var d []byte
	if err = os.WriteFile(filename, d, 0644); err == nil {
		err = os.Remove(filename) // And delete it
		if err != nil {
			return err
		}
		return nil
	}

	return err
}

// CreateShortID создает короткий ID с проверкой на валидность
func CreateShortURL(ctx context.Context, isExist func(context.Context, string) bool) (shortURL string, err error) {
	for i := 0; i < 10; i++ {
		shortURL = NewShortURL()
		if !isExist(ctx, shortURL) {
			return shortURL, nil
		}
	}
	return "", err
}
