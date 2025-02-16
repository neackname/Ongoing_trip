package TravelModel

import (
	uuid "github.com/satori/go.uuid"
)

type TraUser struct {
	ID         uint64 `gorm:"primaryKey"`
	OpenID     string `gorm:"type:varchar(256); not null"`
	Telephone  string `gorm:"type:varchar(128);"`
	NickName   string `gorm:"type:varchar(128)"`
	Motto      string `gorm:"type:varchar(128)"`
	Gender     int    `gorm:"type:int"` //0表示未知，1表示女， 2表示男
	City       string `gorm:"type:varchar(128)"`
	Province   string `gorm:"type:varchar(128)"`
	Country    string `gorm:"type:varchar(128)"`
	AvatarURL  string `gorm:"type:varchar(128)"`
	UnionID    string `gorm:"type:varchar(128)"`
	SessionKey string `gorm:"type:varchar(128)"`

	UserFoot      []TraUserFoot      `gorm:"foreignKey:UserID"`
	UserPostStart []TraUserPostStart `gorm:"foreignKey:UserID"`
	UserFootStart []TraUserFootStart `gorm:"foreignKey:UserID"`

	CreatedAt CustomTime
	UpdatedAt CustomTime
}

// 文章收藏
type TraUserPostStart struct {
	ID     uint64    `gorm:"primaryKey"`
	UserID uint64    `gorm:"not null"` // 外键，关联到User
	PostID uuid.UUID `gorm:"not null"` // 收藏项
}

// 足迹收藏
type TraUserFootStart struct {
	ID     uint64    `gorm:"primaryKey"`
	UserID uint64    `gorm:"not null"` // 外键，关联到User
	FootID uuid.UUID `gorm:"not null"` // 收藏项ID
}

// 用户足迹
type TraUserFoot struct {
	ID     uint64    `gorm:"primaryKey"`
	UserID uint64    `gorm:"not null"` // 外键，关联到User
	FootID uuid.UUID `gorm:"not null"` // 收藏项ID
}
