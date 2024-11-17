package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID      uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	Name      string         `gorm:"not null" json:"name"`
	CreatedAt time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time      `gorm:"default:current_timestamp" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	TodoTag   []TodoTag      `gorm:"foreignKey:tag_id;references:id"`
}

func (t *Tag) TableName() string {
	return "tags"
}
