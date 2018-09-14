package models

type Role struct {
	ID          uint         `gorm:"primary_key"`
	Description string       `gorm:"not null"`
	Permissions []Permission `gorm:"many2many:permission_x_role;"`
}

const (
	ADMIN = 1
	USER  = 2
)

