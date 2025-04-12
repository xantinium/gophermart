package models

import "time"

type OrderStatus uint8

const (
	OrderStatusNew OrderStatus = iota
	OrderStatusProcessing
	OrderStatusInvalid
	OrderStatusProcessed
)

func (status OrderStatus) String() string {
	switch status {
	case OrderStatusNew:
		return "NEW"
	case OrderStatusProcessing:
		return "PROCESSING"
	case OrderStatusInvalid:
		return "INVALID"
	case OrderStatusProcessed:
		return "PROCESSED"
	default:
		return ""
	}
}

func NewUser(id int, login, passwordHash string, created, updated time.Time) User {
	return User{
		id:           id,
		login:        login,
		passwordHash: passwordHash,
		created:      created,
		updated:      updated,
	}
}

type User struct {
	id           int
	login        string
	passwordHash string
	created      time.Time
	updated      time.Time
}

func (user User) Id() int {
	return user.id
}

func (user User) Login() string {
	return user.login
}

func (user User) PasswordHash() string {
	return user.passwordHash
}

func (user User) Created() time.Time {
	return user.created
}

func (user User) Updated() time.Time {
	return user.updated
}
