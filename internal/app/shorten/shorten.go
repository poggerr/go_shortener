package shorten

import (
	"encoding/json"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/logger"
	"math/rand"
)

func Shorting(oldUrl string, strg *storage.Storage) string {
	if oldUrl == "" {
		logger.Initialize().Error("Введите ссылку")
		//err := errors.New("Введите ссылку")
		return ""
	}
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	shortUrl := string(b)

	_, err := strg.Save(shortUrl, oldUrl)
	if err != nil {
		logger.Initialize().Info(err)
		return ""
	}
	return shortUrl
}

func UnShoring(newUrl string, strg *storage.Storage) (string, error) {
	ans, err := strg.OldUrl(newUrl)
	return ans, err
}

type Url struct {
	LongUrl  string `json:"url"`
	ShortUrl string `json:"result"`
}

func JsonCreater(longUrl []uint8, strg *storage.Storage, defUrl string) ([]uint8, error) {
	var url Url
	err := json.Unmarshal(longUrl, &url)
	if err != nil {
		return nil, err
	}

	shortUrl := Shorting(url.LongUrl, strg)
	shortenMap := make(map[string]string)

	shortUrl = defUrl + "/" + shortUrl

	shortenMap["result"] = shortUrl

	marshal, err2 := json.Marshal(shortenMap)
	if err2 != nil {
		return nil, err
	}

	return marshal, nil
}
