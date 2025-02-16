package TravelModel

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID      uuid.UUID `gorm:"primary_key"`
	UserID  uint64    `gorm:"not null"`
	Title   string    `gorm:"type:varchar(30); not null"`
	HeadImg string    `gorm:"type:text"`
	Content string    `gorm:"type:text; not null"`

	CreatedAt CustomTime
	UpdatedAt CustomTime
}

func (post *Post) BeforeCreate(db *gorm.DB) error {
	db.Statement.SetColumn("id", uuid.NewV4())
	return nil
}
