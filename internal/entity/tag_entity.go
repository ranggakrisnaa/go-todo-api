package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID        uint           `gorm:"column:id;primaryKey"`
	UUID      uuid.UUID      `gorm:"column:uuid;type:uuid;default:gen_random_uuid()"`
	Name      string         `gorm:"column:name"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;autoDeleteTime:milli"`
	Todo      []Todo         `gorm:"many2many:todo_tags;foreignKey:ID;joinForeignKey:TagID;References:ID;joinReferences:TodoID"`
}

func (t *Tag) TableName() string {
	return "tags"
}
