package file

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/poggerr/go_shortener/internal/models"
	"github.com/poggerr/go_shortener/internal/service"
	"github.com/poggerr/go_shortener/internal/utils"
	"io"
	"os"
	"strings"
	"sync"
)

type Storage struct {
	storageReader *Reader
	storageWriter *Writer
	mx            sync.Mutex
}

var _ service.URLShortenerService = (*Storage)(nil)

// NewStorage cоздаёт и возвращает экземпляр Storage
func NewStorage(filename string) (fs *Storage, err error) {
	if err = utils.CheckFilename(filename); err != nil {
		return nil, err
	}
	fs = &Storage{}
	fs.storageReader, err = NewReader(filename)
	if err != nil {
		return nil, err
	}
	fs.storageWriter, err = NewWriter(filename)
	if err != nil {
		return nil, err
	}
	return fs, nil
}

// isExist проверяет наличие в файле указанного ID
// Если такой ID входит как подстрока в ссылку, то результат будет такой же, как если бы был найден ID
func (strg *Storage) isExist(_ context.Context, id string) bool {
	_, err := strg.storageReader.file.Seek(0, io.SeekStart)
	if err != nil {
		return false
	}

	scanner := bufio.NewScanner(strg.storageReader.file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), id) {
			// Не обрабатывается ситуация, когда в одной из ссылок может быть подстрока равная ID
			// Для этого можно сделать decoding JSON или захардкодить `"Key:"id"`
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		return false
	}
	return false
}

func (strg *Storage) Store(ctx context.Context, user *uuid.UUID, link string) (id string, err error) {
	strg.mx.Lock()
	defer strg.mx.Unlock()

	id, err = utils.CreateShortURL(ctx, strg.isExist)
	if err != nil {
		return "", err
	}

	err = strg.store(user.String(), id, link)
	if err != nil {
		return "", err
	}

	return id, err
}

func (strg *Storage) store(user string, shortURL string, origURL string) error {
	a := Alias{User: user, ShortURL: shortURL, OriginalURL: origURL}
	err := strg.storageWriter.Write(&a)
	if err != nil {
		return err
	}
	return nil
}

// Restore - находит по ID ссылку во внешнем файле, где данные хранятся в формате JSON
func (strg *Storage) Restore(_ context.Context, id string) (link string, err error) {
	strg.mx.Lock()
	defer strg.mx.Unlock()

	_, err = strg.storageReader.file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	for {
		alias, err := strg.storageReader.Read()
		if err != nil {
			return "", fmt.Errorf("link not found: %s", id)
		}

		if alias.ShortURL == id {
			return alias.OriginalURL, nil
		}
	}
}

// Delete - помечает список ранее сохраненных ссылок удаленными
// только тех ссылок, которые принадлежат пользователю
// Только для совместимости контракта
func (strg *Storage) Delete(_ context.Context, _ *uuid.UUID, _ []string) {
	panic("not implemented for file storage")
}

// GetUserStorage возвращает map[id]link ранее сокращенных ссылок указанным пользователем
func (strg *Storage) GetUserStorage(_ context.Context, user *uuid.UUID, _ string) (map[string]string, error) {
	strg.mx.Lock()
	defer strg.mx.Unlock()

	m := make(map[string]string)
	_, err := strg.storageReader.file.Seek(0, io.SeekStart)
	if err != nil {
		return m, err
	}

	scanner := bufio.NewScanner(strg.storageReader.file)
	for scanner.Scan() {
		txt := scanner.Text()
		if strings.Contains(txt, user.String()) {
			alias := &Alias{}
			dec := json.NewDecoder(bytes.NewBufferString(txt))
			if err = dec.Decode(&alias); err != nil {
				return map[string]string{}, err
			}
			m[alias.ShortURL] = alias.OriginalURL
		}
	}

	if err := scanner.Err(); err != nil {
		return map[string]string{}, err
	}

	return m, nil
}

// StoreBatch сохраняет пакет ссылок из map[correlation_id]original_link и возвращает map[correlation_id]short_link
func (strg *Storage) StoreBatch(ctx context.Context, user *uuid.UUID, batchIn models.BatchList, defURL string) (models.BatchList, error) {
	strg.mx.Lock()
	defer strg.mx.Unlock()

	for corrID, link := range batchIn {
		shortURL, err := utils.CreateShortURL(ctx, strg.isExist)
		if err != nil {
			return nil, err
		}
		err = strg.store(user.String(), shortURL, link.OriginalURL)
		if err != nil {
			return nil, err
		}
		batchIn[corrID].ShortURL = defURL + "/" + shortURL
	}

	return batchIn, nil
}

// Ping проверяет, что файл хранения доступен и экземпляры инициализированы
func (strg *Storage) Ping(_ context.Context) error {
	_, err := strg.storageWriter.file.Stat()
	if err != nil {
		return err
	}
	_, err = strg.storageReader.file.Stat()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Statistics(ctx context.Context) (*models.Statistic, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	var stat models.Statistic
	_, err := s.storageReader.file.Seek(0, io.SeekStart)
	if err != nil {
		panic("file storage in failed state")
	}

	scanner := bufio.NewScanner(s.storageReader.file)
	for scanner.Scan() {
		txt := scanner.Text()
		alias := &Alias{}
		dec := json.NewDecoder(bytes.NewBufferString(txt))
		if err = dec.Decode(&alias); err != nil {
			panic("file storage is corrupted")
		}
		stat.LinksCount++
		stat.UsersCount++
	}

	if err = scanner.Err(); err != nil {
		panic("file storage is corrupted")
	}

	return &stat, nil
}

type Reader struct {
	file    *os.File
	decoder *json.Decoder
}

func NewReader(fileName string) (*Reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return &Reader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *Reader) Read() (*Alias, error) {
	alias := &Alias{}
	if err := c.decoder.Decode(&alias); err != nil {
		return nil, err
	}
	return alias, nil
}

func (c *Reader) Close() error {
	return c.file.Close()
}

type Writer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewWriter(fileName string) (*Writer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &Writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Writer) Write(event *Alias) error {
	return p.encoder.Encode(&event)
}

func (p *Writer) Close() error {
	return p.file.Close()
}

func (strg *Storage) Close() error {
	err1 := strg.storageReader.Close()
	err2 := strg.storageWriter.Close()
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

// Alias - структура хранения ShortURL и OriginalURL во внешнем файле
type Alias struct {
	User        string
	ShortURL    string
	OriginalURL string
}
