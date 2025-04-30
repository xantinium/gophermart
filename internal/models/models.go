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

func NewOrder(id int, number string, userID int, status OrderStatus, accrual float32, created, updated time.Time) Order {
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
	accrual float32
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

func (order Order) Accrual() float32 {
	return order.accrual
}

func NewWithdrawal(id int, order string, sum float32, userID int, created, updated time.Time) Withdrawal {
	return Withdrawal{
		baseStruct: baseStruct{
			id:      id,
			created: created,
			updated: updated,
		},
		order:  order,
		sum:    sum,
		userID: userID,
	}
}

type Withdrawal struct {
	baseStruct

	order  string
	sum    float32
	userID int
}

func (withdrawal Withdrawal) Order() string {
	return withdrawal.order
}

func (withdrawal Withdrawal) Sum() float32 {
	return withdrawal.sum
}

func (withdrawal Withdrawal) UserID() int {
	return withdrawal.userID
}
