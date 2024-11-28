package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"column:id;primaryKey"`
	UUID      uuid.UUID      `gorm:"column:uuid;type:uuid;default:gen_random_uuid()"`
	Name      string         `gorm:"column:name"`
	Email     string         `gorm:"column:email"`
	Password  string         `gorm:"column:password"`
	Token     string         `gorm:"column:token"`
	Role      string         `gorm:"column:role;default:'user'"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;autoDeleteTime:milli"`
	Todo      []Todo         `gorm:"foreignKey:user_id;references:id"`
}

func (u *User) TableName() string {
	return "users"
}
