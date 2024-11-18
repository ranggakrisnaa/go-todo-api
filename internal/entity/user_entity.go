package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID      uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	Email     string         `gorm:"type:varchar(255);unique;not null" json:"email" validate:"required"`
	Password  string         `gorm:"type:text;not null" json:"password" validate:"required"`
	CreatedAt time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time      `gorm:"default:current_timestamp" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Todo      []Todo         `gorm:"foreignKey:user_id;references:id"`
}

func (u *User) TableName() string {
	return "users"
}
