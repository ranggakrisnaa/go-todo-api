package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID        uint           `gorm:"column:id;primaryKey"`
	UUID      uuid.UUID      `gorm:"column:uuid;type:uuid;default:gen_random_uuid()"`
	Name      string         `gorm:"column:name"`
	CreatedAt int64          `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64          `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;autoDeleteTime:milli"`
	TodoTag   []TodoTag      `gorm:"foreignKey:tag_id;references:id"`
}

func (t *Tag) TableName() string {
	return "tags"
}
