package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Todo struct {
	ID          uint           `gorm:"column:id;primaryKey"`
	UUID        uuid.UUID      `gorm:"column:uuid"`
	UserID      uint           `gorm:"column:user_id"`
	Title       string         `gorm:"column:title"`
	Description string         `gorm:"column:description"`
	IsCompleted bool           `gorm:"column:is_completed"`
	CreatedAt   int64          `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt   int64          `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;autoDeleteTime:milli"`
	User        User           `gorm:"foreignKey:user_id;references:id"`
	TodoTag     []TodoTag      `gorm:"foreignKey:todo_id;references:id"`
}

func (t *Todo) TableName() string {
	return "todos"
}
