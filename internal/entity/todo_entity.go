package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Todo struct {
	ID          uint           `gorm:"column:id;primaryKey"`
	UUID        uuid.UUID      `gorm:"column:uuid;type:uuid;default:gen_random_uuid()"`
	UserID      uint           `gorm:"column:user_id"`
	Title       string         `gorm:"column:title"`
	Description string         `gorm:"column:description"`
	IsCompleted bool           `gorm:"column:is_completed"`
	DueTime     time.Time      `gorm:"column:due_time"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;autoDeleteTime:milli"`
	User        User           `gorm:"foreignKey:user_id;references:id"`
	Tag         []Tag          `gorm:"many2many:todo_tags;foreignKey:ID;joinForeignKey:TodoID;References:ID;joinReferences:TagID"`
}

func (t *Todo) TableName() string {
	return "todos"
}
