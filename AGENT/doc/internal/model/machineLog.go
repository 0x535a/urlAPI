package model

import (
	"time"
)

type MachineLog struct {
	ID               int       `gorm:"primaryKey;column:LogMachineryID" json:"id"`
	MacID            int       `gorm:"column:MacID" json:"macID"`
	BeginTime        time.Time `gorm:"column:BeginTime" json:"beginTime"`
	EndTime          time.Time `gorm:"column:EndTime" json:"endTime"`
	LogMachineryType int       `gorm:"column:LogMachineryType" json:"logMachineryType"`
}

type MachinePasID struct {
	PasID int `gorm:"column:PasID" json:"pasID"`
	MacID int `gorm:"column:MacID" json:"macID"`
}

func (MachineLog) TableName() string {
	return "LogMachinery"
}

func (MachinePasID) TableName() string {
	return "Passageway"
}
