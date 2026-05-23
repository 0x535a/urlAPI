package model

type Machine struct {
	ID              string  `gorm:"primaryKey;column:id" json:"id"`
	Name            string  `gorm:"column:name;default:''" json:"name"`
	IsAvailable     bool    `gorm:"column:isAvailable;default:true" json:"isAvailable"`
	UnitPriceYuan   float64 `gorm:"column:unitPriceYuan;default:0" json:"unitPriceYuan"`
	Notes           string  `gorm:"column:notes;default:''" json:"notes"`
	WorkingPeriodID string  `gorm:"column:workingPeriodID;default:''" json:"workingPeriodID"`
	OrderID         string  `gorm:"column:orderID;default:''" json:"orderID"`
	LastSyncTime    int64   `gorm:"column:lastSyncTime;default:0" json:"lastSyncTime"`
}
