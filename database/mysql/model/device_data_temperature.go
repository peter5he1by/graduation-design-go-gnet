package model

import "time"

type DeviceDataTemperature struct {
	BaseModel
	DeviceID uint      `json:"deviceId"` // 所属设备id
	Time     time.Time `json:"time"`
	Data     float64   `json:"data"`

	Device Device `json:"device" gorm:"foreignKey:ID;references:DeviceID"` // 所属设备（many-to-one）
}

func (receiver DeviceDataTemperature) TableName() string {
	return "device_data_temperature"
}
