package models

type Permission struct {
	ID          uint   `gorm:"primary_key" json:"-"`
	Description string `gorm:"not null"`
}
