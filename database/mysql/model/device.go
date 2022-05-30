package model

// Device 设备连接信息（已连接/未连接），这个作为整个软件的工作主体
type Device struct {
	BaseModel
	Uuid               string                 `json:"uuid"`
	SecretKey          string                 `json:"secretKey"`
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	InstallationSiteId uint                   `json:"installationSiteId"`
	InstallationSite   DeviceInstallationSite `json:"installationSite" gorm:"foreignKey:ID;references:InstallationSiteId"`
}

func (d Device) TableName() string {
	return "device"
}
