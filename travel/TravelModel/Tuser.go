package TravelModel

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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
	UserID uint64    `gorm:"not null"`               // 外键，关联到User
	PostID uuid.UUID `gorm:"type:char(36);not null"` // 收藏项
}

type TraUserPostLike struct {
	ID     uint64    `gorm:"primaryKey"`
	UserID uint64    `gorm:"not null;uniqueIndex:idx_user_post_like"`
	PostID uuid.UUID `gorm:"type:char(36);not null;uniqueIndex:idx_user_post_like"`
}

// 足迹收藏
type TraUserFootStart struct {
	ID     uint64    `gorm:"primaryKey"`
	UserID uint64    `gorm:"not null"`               // 外键，关联到User
	FootID uuid.UUID `gorm:"type:char(36);not null"` // 收藏项ID
}

// 用户足迹
type TraUserFoot struct {
	ID     uint64    `gorm:"primaryKey"`
	UserID uint64    `gorm:"not null"`               // 外键，关联到User
	FootID uuid.UUID `gorm:"type:char(36);not null"` // 收藏项ID
}

type TraFoot struct {
	ID               uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID           uint64    `gorm:"not null;index"`
	Title            string    `gorm:"type:varchar(128)"`
	Origin           string    `gorm:"type:varchar(64);not null"`
	OriginName       string    `gorm:"type:varchar(128)"`
	Destinations     string    `gorm:"type:text;not null"`
	DestinationNames string    `gorm:"type:text"`
	Mode             string    `gorm:"type:varchar(32);not null"`
	RouteResult      string    `gorm:"type:longtext;not null"`

	CreatedAt CustomTime
	UpdatedAt CustomTime
}

func (foot *TraFoot) BeforeCreate(db *gorm.DB) error {
	db.Statement.SetColumn("id", uuid.NewV4())
	return nil
}

type ChatMessage struct {
	ID         uuid.UUID `gorm:"type:char(36);primary_key"`
	FromUserID uint64    `gorm:"not null;index"`
	ToUserID   uint64    `gorm:"not null;index"`
	Content    string    `gorm:"type:text;not null"`

	CreatedAt CustomTime
	UpdatedAt CustomTime
}

func (m *ChatMessage) BeforeCreate(db *gorm.DB) error {
	db.Statement.SetColumn("id", uuid.NewV4())
	return nil
}
