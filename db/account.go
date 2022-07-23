package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	ID           int64
	UserName     string
	UserEmail    string
	Host         string
	CreatedAt    time.Time
	LastActiveAt time.Time
	Banned       bool
}

var (
	ErrAccountAlreadyExist = errors.New("account with such user name already exists")
	ErrEmailIsAlreadyUsed  = errors.New("email is already used")
)

func (d *DB) CreateAccount(name, email string) error {
	var acc *Account
	d.FindAccountByEmail(email, acc)
	if acc.ID > 0 {
		return ErrEmailIsAlreadyUsed
	}

	d.FindAccountByUserName(name, acc)
	if acc.ID > 0 {
		return ErrAccountAlreadyExist
	}

	d.Conn.Create(Account{
		UserName:     name,
		UserEmail:    email,
		LastActiveAt: time.Now().UTC(),
		CreatedAt:    time.Now().UTC(),
		Banned:       false,
	})
	return nil
}

func (d *DB) FindAccountByUserName(userName string, acc *Account) {
	d.Conn.Find(acc).Where("user_name = ?", userName)
}

func (d *DB) FindAccountByEmail(email string, acc *Account) {
	d.Conn.Find(acc).Where("email = ?", email)
}

func (d *DB) FindAccountByID(id int, acc *Account) {
	d.Conn.Find(acc).Where("id = ?", id)
}

func (d *DB) UpdateAccountHost(id int, host string) {
	d.Conn.Model(&Account{}).Where("id = ?", id).Update("host", host)
}

func (d *DB) BanAccount(id int64) {
	d.Conn.Model(&Account{}).Where("id = ?", id).Update("banned", true)
}
