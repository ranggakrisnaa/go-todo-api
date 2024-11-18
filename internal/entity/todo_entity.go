package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Todo struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text;not null" json:"description"`
	IsCompleted bool           `gorm:"not null;default:true" json:"is_completed"`
	DueTime     time.Time      `gorm:"not null" json:"due_time"`
	CreatedAt   time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"default:current_timestamp" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	User        User           `gorm:"foreignKey:user_id;references:id"`
	TodoTag     []TodoTag      `gorm:"foreignKey:todo_id;references:id"`
}

func (t *Todo) TableName() string {
	return "todos"
}
