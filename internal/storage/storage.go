// Package storage содержит код для взаимодействия с хранилищем
package storage

import (
	"database/sql"
)

// Storage хранилище
type Storage struct {
	data map[string]string
	path string
	DB   *sql.DB
}

// NewStorage создание нового хранилища
func NewStorage(p string, db *sql.DB) *Storage {
	return &Storage{
		data: make(map[string]string),
		path: p,
		DB:   db,
	}
}
