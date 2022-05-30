package model

// DeviceInfo 设备维护信息（记录，旧信息保留以供查询）
type DeviceInfo struct {
	BaseModel
	DeviceID     uint   `json:"deviceId"`     // 所属设备id
	SoftwareInfo string `json:"softwareInfo"` // 软件版本信息
	HardwareInfo string `json:"hardwareInfo"` // 硬件版本信息
	Remark       string `json:"remark"`       // 备注，用来保存所有额外信息

	Device Device `json:"device" gorm:"foreignKey:ID;references:DeviceID"` // 所属设备（many-to-one）
}

func (i DeviceInfo) TableName() string {
	return "device_info"
}
