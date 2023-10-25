package memory

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/poggerr/go_shortener/internal/handlers"
	"github.com/poggerr/go_shortener/internal/models"
	"github.com/poggerr/go_shortener/internal/utils"
	"sync"
)

// Storage реализует хранение ссылок в памяти.
// Является потоко безопасной реализацией Repository
type Storage struct {
	storage map[string]map[string]string
	mx      sync.Mutex
}

var _ handlers.Repository = (*Storage)(nil)

// NewStorage cоздает и возвращает экземпляр Storage
func NewStorage() *Storage {
	s := Storage{}
	s.storage = make(map[string]map[string]string)
	return &s
}

// Store сохраняет ссылку в хранилище с указанным id
func (s *Storage) Store(ctx context.Context, user *uuid.UUID, link string) (id string, err error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	id, err = utils.CreateShortURL(ctx, s.isExist)
	if err != nil {
		return "", err
	}

	if _, ok := s.storage[user.String()]; !ok {
		s.storage[user.String()] = make(map[string]string)
	}

	s.storage[user.String()][id] = link
	return id, nil
}

// isExist проверяет наличие id в сторадже
func (s *Storage) isExist(_ context.Context, id string) bool {
	for _, user := range s.storage {
		_, ok := user[id]
		if ok {
			return true
		}
	}
	return false
}

// Restore возвращает исходную ссылку по переданному короткому ID
func (s *Storage) Restore(_ context.Context, id string) (link string, err error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	for _, user := range s.storage {
		l, ok := user[id]
		if ok {
			return l, nil
		}
	}

	return "", fmt.Errorf("ссылка не найдена: %s", id)
}

// Unstore - помечает список ранее сохраненных ссылок удаленными
// только тех ссылок, которые принадлежат пользователю
// Только для совместимости контракта
func (s *Storage) Unstore(_ context.Context, _ string, _ []string) {
	panic("not implemented for memory storage")
}

// GetUserStorage возвращает map[id]link ранее сокращенных ссылок указанным пользователем
func (s *Storage) GetUserStorage(_ context.Context, user *uuid.UUID, _ string) (map[string]string, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	ub, ok := s.storage[user.String()]
	if !ok {
		return nil, fmt.Errorf("пользователь не найден")
	}
	return ub, nil
}

// StoreBatch сохраняет пакет ссылок из map[correlation_id]original_link и возвращает map[correlation_id]short_link
func (s *Storage) StoreBatch(ctx context.Context, user *uuid.UUID, batchIn models.BatchList, _ string) (models.BatchList, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if _, ok := s.storage[user.String()]; !ok {
		s.storage[user.String()] = make(map[string]string)
	}

	for corrID, link := range batchIn {
		id, err := utils.CreateShortURL(ctx, s.isExist)
		if err != nil {
			return nil, err
		}
		s.storage[user.String()][id] = link.ShortURL
		batchIn[corrID].ShortURL = id
	}

	return batchIn, nil
}

func (s *Storage) Delete(_ context.Context, _ *uuid.UUID, _ []string) {
	return
}

// Ping проверяет, что экземпляр Storage создан корректно, например с помощью NewStorage()
func (s *Storage) Ping(_ context.Context) error {
	if s.storage == nil {
		return fmt.Errorf("storage is nil")
	}
	return nil
}

// Close ничего не делает, требуется только для совместимости с контрактом
func (s *Storage) Close() error {
	// Do nothing
	return nil
}
