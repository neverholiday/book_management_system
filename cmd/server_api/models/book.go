package models

import "time"

type Book struct {
	ID                string     `gorm:"column:id"`
	Title             string     `gorm:"column:title"`
	Author            string     `gorm:"column:author"`
	ISBN              *string    `gorm:"column:isbn"`
	Publisher         *string    `gorm:"column:publisher"`
	PublicationYear   *int       `gorm:"column:publication_year"`
	Genre             *string    `gorm:"column:genre"`
	Description       *string    `gorm:"column:description"`
	Pages             *int       `gorm:"column:pages"`
	Language          string     `gorm:"column:language"`
	Price             *float64   `gorm:"column:price"`
	Quantity          int        `gorm:"column:quantity"`
	AvailableQuantity int        `gorm:"column:available_quantity"`
	Location          *string    `gorm:"column:location"`
	Status            string     `gorm:"column:status"`
	CreatedDate       time.Time  `gorm:"column:created_date"`
	UpdatedDate       time.Time  `gorm:"column:updated_date"`
	DeletedDate       *time.Time `gorm:"column:deleted_date"`
}