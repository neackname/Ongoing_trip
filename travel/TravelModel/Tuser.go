package TravelModel

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type TraUser struct {
	gorm.Model
	OpenID     string `gorm:"type:varchar(256); not null"`
	Telephone  string `gorm:"type:varchar(128);"`
	NickName   string `gorm:"type:varchar(128)"`
	Motto      string `gorm:"type:varchar(128)"`
	Gender     int    `gorm:"type:int"` //0表示男、1表示女
	City       string `gorm:"type:varchar(128)"`
	Province   string `gorm:"type:varchar(128)"`
	Country    string `gorm:"type:varchar(128)"`
	AvatarURL  string `gorm:"type:varchar(128)"`
	UnionID    string `gorm:"type:varchar(128)"`
	SessionKey string `gorm:"type:varchar(128)"`

	UserFoot      []TraUserFoot      `gorm:"foreignKey:UserID"`
	UserPostStart []TraUserPostStart `gorm:"foreignKey:UserID"`
	UserFootStart []TraUserFootStart `gorm:"foreignKey:UserID"`
}

type TraUserPostStart struct {
	ID     uint      `gorm:"primaryKey"`
	UserID uint      `gorm:"not null"` // 外键，关联到User
	PostID uuid.UUID `gorm:"not null"` // 收藏项ID
	Name   string    // 收藏项名称
}

type TraUserFootStart struct {
	ID     uint      `gorm:"primaryKey"`
	UserID uint      `gorm:"not null"` // 外键，关联到User
	FootID uuid.UUID `gorm:"not null"` // 收藏项ID
	Name   string    // 收藏项名称
}

type TraUserFoot struct {
	ID     uint      `gorm:"primaryKey"`
	UserID uint      `gorm:"not null"` // 外键，关联到User
	FootID uuid.UUID `gorm:"not null"` // 收藏项ID
	Name   string    // 收藏项名称
}
