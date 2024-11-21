package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TodoTag struct {
	ID        uint           `gorm:"column:id;primaryKey"`
	UUID      uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	TodoID    uint           `gorm:"not null" json:"todo_id" validate:"required"`
	TagID     uint           `gorm:"not null" json:"tag_id"`
	CreatedAt int64          `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64          `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;autoDeleteTime:milli"`
	Tag       Tag            `gorm:"foreignKey:tag_id;references:id"`
	Todo      Todo           `gorm:"foreignKey:todo_id;references:id"`
}

func (t *TodoTag) TableName() string {
	return "todo_tags"
}
