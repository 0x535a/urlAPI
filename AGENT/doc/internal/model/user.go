package model

type User struct {
	ID          string `gorm:"primaryKey;column:id" json:"id"`
	Type        string `gorm:"column:type" json:"type"`
	Name        string `gorm:"unique;column:name" json:"name"`
	Phone       string `gorm:"column:phone" json:"phone"`
	IsWorking   bool   `gorm:"column:isWorking" json:"isWorking"`
	WXID        string `gorm:"column:wxid" json:"wxid"`
	WorkPeriods string `gorm:"column:workPeriods" json:"-"`
}
