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

func (d *DB) CreateAccount(name, email, host string) error {
	var acc Account
	if err := d.FindAccountByEmail(email, &acc); err != nil {
		return err
	}
	if acc.ID > 0 {
		return ErrEmailIsAlreadyUsed
	}

	if err := d.FindAccountByUserName(name, &acc); err != nil {
		return err
	}

	if acc.ID > 0 {
		return ErrAccountAlreadyExist
	}

	acc.UserName = name
	acc.UserEmail = email
	acc.Host = host
	acc.CreatedAt = time.Now().UTC()
	acc.LastActiveAt = time.Now().UTC()

	return d.Conn.Create(&acc).Error
}

func (d *DB) FindAccountByUserName(userName string, acc *Account) error {
	return d.Conn.Debug().Find(acc, "user_name = ?", userName).Error
}

func (d *DB) FindAccountByEmail(email string, acc *Account) error {
	return d.Conn.Debug().Find(acc, "user_email = ?", email).Error
}

func (d *DB) FindAccountByID(id int, acc *Account) error {
	return d.Conn.Debug().Find(acc, "id = ?", id).Error
}

func (d *DB) UpdateAccountHost(id int64, host string) error {
	return d.Conn.Debug().Model(&Account{}).Where("id = ?", id).Update("host", host).Error
}

func (d *DB) BanAccount(id int64) error {
	return d.Conn.Model(&Account{}).Where("id = ?", id).Update("banned", true).Error
}
