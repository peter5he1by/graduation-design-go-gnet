package model

// DeviceInstallationSite 设备安装场所
type DeviceInstallationSite struct {
	ID          uint     `gorm:"primarykey"`
	ParentId    *uint    `json:"parentId"`                                     // 树形结构
	Name        string   `json:"name"`                                         // 名称
	Description string   `json:"description"`                                  // 描述
	Devices     []Device `json:"devices" gorm:"foreignKey:InstallationSiteId"` // one-to-many
	// self-referential
	ParentInstallationSite *DeviceInstallationSite `json:"parentInstallationSite" gorm:"foreignKey:ID;references:ParentId"`
}

// type DIS DeviceInstallationSite

func (l DeviceInstallationSite) TableName() string {
	return "device_installation_site"
}
