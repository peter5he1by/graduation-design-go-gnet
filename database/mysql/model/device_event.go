package model

// DeviceEvent 设备工作信息
type DeviceEvent struct {
	BaseModel
	DeviceID uint   `json:"deviceId"` // 所属设备ID
	Type     string `json:"eventType"`

	Device Device `gorm:"foreignKey:ID;references:DeviceID"`
}

func (e DeviceEvent) TableName() string {
	return "device_event"
}
