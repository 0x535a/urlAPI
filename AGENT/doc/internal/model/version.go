package model

type Version struct {
	ID        string `gorm:"primaryKey;column:id" json:"id"`
	Type      string `gorm:"column:type" json:"type"`
	UserID    string `gorm:"column:userID" json:"userID"`
	ContentID string `gorm:"column:contentID" json:"contentID"`
	From      string `gorm:"column:from" json:"from"`
	To        string `gorm:"column:to" json:"to"`
	Time      int64  `gorm:"column:time" json:"time"`
}
