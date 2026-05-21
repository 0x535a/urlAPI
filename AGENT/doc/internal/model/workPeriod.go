package model

type WorkPeriod struct {
	ID                  string          `gorm:"primaryKey;column:id" json:"id"`
	WorkerID            string          `gorm:"column:workerID" json:"workerID" form:"workerID"`
	OrderID             string          `gorm:"column:orderID" json:"orderID" form:"orderID"`
	MachineID           string          `gorm:"column:machineID" json:"machineID" form:"machineID"`
	StartTime           int64           `gorm:"column:startTime" json:"startTime" form:"startTime"`
	EndTime             int64           `gorm:"column:endTime" json:"endTime" form:"endTime"`
	IsMachineDataMerged bool            `gorm:"column:isMachineDataMerged" json:"isMachineDataMerged"`
	ValidTimeSeconds    int64           `gorm:"column:validTimeSeconds" json:"validTimeSeconds"`
	MachineName         string          `gorm:"column:machineName" json:"machineName"`
	OrderName           string          `gorm:"column:orderName" json:"orderName"`
	UnitPriceYuan       float64         `gorm:"column:unitPriceYuan" json:"unitPriceYuan"`
	MachineONPeriods    []MachinePeriod `gorm:"foreignKey:WorkPeriodID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"machineONPeriods"`
}

type MachinePeriod struct {
	WorkPeriodID string `gorm:"column:workPeriodID" json:"workPeriodID"` // 外键字段
	StartTime    int64  `gorm:"column:startTime" json:"startTime"`
	EndTime      int64  `gorm:"column:endTime" json:"endTime"`
}
