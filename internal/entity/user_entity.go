package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	UUID      uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid()"`
	Name      string         `gorm:"type:varchar(100);not null"`
	Email     string         `gorm:"type:varchar(255);unique;not null"`
	Password  string         `gorm:"type:text;not null"`
	CreatedAt time.Time      `gorm:"default:current_timestamp"`
	UpdatedAt time.Time      `gorm:"default:current_timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Todo      []Todo         `gorm:"foreignKey:user_id;references:id"`
}

func (u *User) TableName() string {
	return "users"
}
