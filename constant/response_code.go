package constant

const (
	ServerOK                        uint16 = 0x0000
	ServerRequestUpdateDeviceConfig uint16 = 0x0011 // 网关要求上传最新配置信息
	ServerIssueDeviceConfig         uint16 = 0x0012 // 网关下发配置信息
)
