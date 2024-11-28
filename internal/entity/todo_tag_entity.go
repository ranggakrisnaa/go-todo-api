package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TodoTag struct {
	ID        uint           `gorm:"column:id;primaryKey"`
	UUID      uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid()"`
	TodoID    uint           `gorm:"not null"`
	TagID     uint           `gorm:"not null"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;autoDeleteTime:milli"`
}

func (t *TodoTag) TableName() string {
	return "todo_tags"
}
