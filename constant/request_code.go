package constant

const (
	ClientConnect    uint16 = 0x0001 // 边缘计算设备使用uuid和key登录通讯网关
	ClientDisconnect uint16 = 0x0002 // 边缘计算设备与通讯网关断开连接

	/* 信息 */

	ClientUpdateDeviceInfo   uint16 = 0x0012 // 更新设备维护信息
	ClientUploadDeviceLog    uint16 = 0x0015 // 上传日志
	ClientUploadDeviceConfig uint16 = 0x0017 // 上传配置

	/* 事件 */

	ClientEventUserLogin  uint16 = 0x0021 // 工作人员登录
	ClientEventUserLogout uint16 = 0x0022 // 工作人员注销
	ClientEventOperation  uint16 = 0x0025 // 工作人员操作记录

	/* 实时数据上传 */

	ClientDataUploadTemperature uint16 = 0x0101 // 温度

	ClientHeartbeat uint16 = 0xffff
)
