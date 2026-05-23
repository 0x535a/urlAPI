package model

type Order struct {
	ID        string `gorm:"column:id;primaryKey" json:"id" cn:"系统ID"`
	VersionID string `gorm:"column:versionID;primaryKey" json:"-" cn:"版本ID"`
	FactoryID string `gorm:"column:factoryID" json:"factoryID" cn:"订单ID"`
	//
	IsLatest bool `gorm:"column:isLatest" json:"isLatest" cn:"是否最新"`

	// has
	Name  string `gorm:"column:name" json:"name" cn:"订单名称"`
	Notes string `gorm:"column:notes;default:''" json:"notes" cn:"订单备注"`

	// eq
	Status int    `gorm:"column:status" json:"status" cn:"订单状态"`
	Type   string `gorm:"column:type" json:"type" cn:"订单类型"`

	// ge & le
	Priority              int   `gorm:"column:priority" json:"priority" cn:"优先级"`
	PlannedStartTime      int64 `gorm:"column:plannedStartTime" json:"plannedStartTime" cn:"计划开始时间"`
	PlannedEndTime        int64 `gorm:"column:plannedEndTime" json:"plannedEndTime" cn:"计划结束时间"`
	AssignedTime          int64 `gorm:"column:assignedTime" json:"assignedTime" cn:"分配时间"`
	MarkedFinishedTime    int64 `gorm:"column:markedFinishedTime" json:"markedFinishedTime" cn:"标记完成时间"`
	ConfirmedFinishedTime int64 `gorm:"column:confirmedFinishedTime" json:"confirmedFinishedTime" cn:"确认完成时间"`

	AssignedToes       string             `gorm:"column:assignedToes" json:"-" cn:"分配对象"`
	AssignedMachinesDB string             `gorm:"column:assignedMachines;default:''" json:"-" cn:"分配的机器"`
	FinishedWorkersDB  string             `gorm:"column:finishedWorkers" json:"-" cn:"完成的工人"`
	WorkPeriodsDB      string             `gorm:"column:workPeriods" json:"-" cn:"工作周期"`
	ColorAndQuantities []ColorAndQuantity `gorm:"foreignKey:OrderID,VersionID;references:ID,VersionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"colorsAndQuantities" cn:"颜色与数量"`
}

type ColorAndQuantity struct {
	OrderID   string `gorm:"column:orderID" json:"orderID"`
	VersionID string `gorm:"column:versionID" json:"versionID"`
	Color     string `gorm:"column:color" json:"color"`
	Quantity  string `gorm:"column:quantity" json:"quantity"`
}
