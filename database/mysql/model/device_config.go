package model

// DeviceConfig 设备配置
type DeviceConfig struct {
	BaseModel
	DeviceID uint   `json:"deviceId"` // 所属设备id
	Content  string `json:"content"`
	Type     string `json:"type"`

	Device Device `json:"device" gorm:"foreignKey:ID;references:DeviceID"` // 所属设备（many-to-one）
}

func (c DeviceConfig) TableName() string {
	return "device_config"
}
