package databases

import (
	"time"
)

type Time struct {
	BasicTimeFields
	DeletedAt *time.Time `gorm:"column:deletedAt;null"`
}

type BasicTimeFields struct {
	CreatedAt time.Time `gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"column:updatedAt;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}
