package models

import "time"

type User struct {
	ID           string     `gorm:"column:id"`
	Email        string     `gorm:"column:email"`
	PasswordHash string     `gorm:"column:password_hash"`
	FirstName    string     `gorm:"column:first_name"`
	LastName     string     `gorm:"column:last_name"`
	Role         string     `gorm:"column:role"`
	Status       string     `gorm:"column:status"`
	CreatedDate  time.Time  `gorm:"column:created_date"`
	UpdatedDate  time.Time  `gorm:"column:updated_date"`
	DeletedDate  *time.Time `gorm:"column:deleted_date"`
}

func (u *User) GetID() string {
	return u.ID
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRole() string {
	return u.Role
}