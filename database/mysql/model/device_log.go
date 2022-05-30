package model

import (
	"time"
)

// DeviceLog 设备日志
type DeviceLog struct {
	BaseModel
	DeviceID  uint      `json:"deviceId"`
	Content   string    `json:"content"`   // 日志内容，应该会比较大
	StartTime time.Time `json:"startTime"` // 日志时间（或者叫起始时间，具体业务具体使用）
	EndTime   time.Time `json:"endTime"`   // 日志时间（或者叫结束时间，具体业务具体使用）

	Device Device `gorm:"foreignKey:ID;references:DeviceID"`
}

func (e DeviceLog) TableName() string {
	return "device_log"
}
