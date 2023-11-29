package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/poggerr/go_shortener/internal/models"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type URLShortenerService interface {
	// Store сохраняет оригинальную ссылку и возвращает id (токен) сокращенного варианта.
	Store(ctx context.Context, user *uuid.UUID, longURL string) (id string, err error)
	// Restore возвращает оригинальную ссылку по его id.
	Restore(ctx context.Context, id string) (link string, err error)
	// Delete - помечает ссылки удаленными.
	// Согласно заданию - результат работы пользователю не возвращается.
	Delete(ctx context.Context, user *uuid.UUID, ids []string)
	// GetUserStorage возвращает массив всех ранее сокращенных пользователей ссылок.
	GetUserStorage(ctx context.Context, user *uuid.UUID, defURL string) (map[string]string, error)
	// StoreBatch сохраняет пакет ссылок в хранилище и возвращает список пакет id.
	StoreBatch(ctx context.Context, user *uuid.UUID, batchIn models.BatchList, defURL string) (models.BatchList, error)
	// Ping проверяет готовность к работе репозитория.
	Ping(context.Context) error

	Statistics(ctx context.Context) (*models.Statistic, error)
	// Close завершает работу репозитория в стиле graceful shutdown.
	Close() error
}
