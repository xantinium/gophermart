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

type baseStruct struct {
	id      int
	created time.Time
	updated time.Time
}

func (base baseStruct) ID() int {
	return base.id
}

func (base baseStruct) Created() time.Time {
	return base.created
}

func (base baseStruct) Updated() time.Time {
	return base.updated
}

func NewUser(id int, login, passwordHash string, created, updated time.Time) User {
	return User{
		baseStruct: baseStruct{
			id:      id,
			created: created,
			updated: updated,
		},
		login:        login,
		passwordHash: passwordHash,
	}
}

type User struct {
	baseStruct

	login        string
	passwordHash string
}

func (user User) Login() string {
	return user.login
}

func (user User) PasswordHash() string {
	return user.passwordHash
}

func NewOrder(id int, number string, userID int, status OrderStatus, accrual *int, created, updated time.Time) Order {
	return Order{
		baseStruct: baseStruct{
			id:      id,
			created: created,
			updated: updated,
		},
		number:  number,
		userID:  userID,
		status:  status,
		accrual: accrual,
	}
}

type Order struct {
	baseStruct

	number  string
	userID  int
	status  OrderStatus
	accrual *int
}

func (order Order) Number() string {
	return order.number
}

func (order Order) UserID() int {
	return order.userID
}

func (order Order) Status() OrderStatus {
	return order.status
}

func (order Order) Accrual() *int {
	return order.accrual
}
