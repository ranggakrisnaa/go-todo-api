package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TodoTag struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID      uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	TodoID    uint           `gorm:"not null" json:"todo_id" validate:"required"`
	TagID     uint           `gorm:"not null" json:"tag_id"`
	CreatedAt time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time      `gorm:"default:current_timestamp" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Tag       Tag            `gorm:"foreignKey:tag_id;references:id"`
	Todo      Todo           `gorm:"foreignKey:todo_id;references:id"`
}

func (t *TodoTag) TableName() string {
	return "todo_tags"
}
