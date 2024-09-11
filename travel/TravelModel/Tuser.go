package TravelModel

import "gorm.io/gorm"

type TraUser struct {
	gorm.Model
	OpenID     string `json:"openId" gorm:"type:varchar(20); not null"`
	NickName   string `json:"nickName" gorm:"type:varchar(128)"`
	Gender     int    `json:"gender"`
	City       string `json:"city"`
	Province   string `json:"province"`
	Country    string `json:"country"`
	AvatarURL  string `json:"avatarUrl"`
	UnionID    string `json:"unionId,omitempty"`
	SessionKey string `gorm:"type:varchar(128)"`
}
