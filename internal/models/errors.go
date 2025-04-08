package models

import "errors"

var (
	// ErrNotFound ошибка: объект не найден.
	ErrNotFound = errors.New("not found")
)
