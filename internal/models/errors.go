package models

import "errors"

var (
	// ErrNotFound ошибка: объект не найден.
	ErrNotFound = errors.New("not found")
	// ErrNotFound ошибка: объект не является уникальным.
	ErrAlreadyExists = errors.New("already exists")
	// ErrAlreadyExists ошибка: не удалось сравнить значения.
	ErrFailedToMatch = errors.New("failed to match")
	// ErrInvalidOrderNum ошибка: невалидный номер заказа.
	ErrInvalidOrderNum = errors.New("invalid order number")
)
