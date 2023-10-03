package storage

import (
	"encoding/json"
	"os"

	"github.com/poggerr/go_shortener/internal/logger"
)

// Save базовое сохранение
func (strg *Storage) Save(key, value string) {
	strg.data[key] = value
	if strg.path != "" {
		strg.SaveToFile()
	}
}

// SaveToFile сохранение в файл
func (strg *Storage) SaveToFile() {
	file, err := os.OpenFile(strg.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			logger.Initialize().Error(err)
		}
	}(file)
	if err != nil {
		logger.Initialize().Error(err)
	}

	data, _ := json.Marshal(strg.data)

	data = append(data, '\n')

	_, err = file.Write(data)
	if err != nil {
		logger.Log.Error("Ошибка при сохранении файла ", err)
	}
}
