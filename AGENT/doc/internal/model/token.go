package model

type Token struct {
	Token      string `gorm:"primaryKey;column:token" json:"token"`
	UserID     string `gorm:"column:userID" json:"userID"`
	UserType   string `gorm:"column:userType" json:"userType"`
	ExpireTime int64  `gorm:"column:expireTime" json:"expireTime"`
}
