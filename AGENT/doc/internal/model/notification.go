package model

const (
	OrderNew = iota + 100
	OrderConfirm
	OrderComplete
	OrderUpdate
	OrderDelete
)

type Notification struct {
	ID        string `gorm:"primaryKey;column:id" json:"id"`
	UserID    string `gorm:"column:userID" json:"userID"`
	Time      int64  `gorm:"column:time" json:"time"`
	TypeCode  int    `gorm:"column:type" json:"typeCode"`
	ContentID string `gorm:"column:content" json:"contentID"`
	ExtraInfo string `gorm:"column:extraInfo" json:"extraInfo"`
}
