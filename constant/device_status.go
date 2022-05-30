package constant

type DeviceStatus int32

const (
	DeviceStatusUnknown DeviceStatus = 0b0000 // 未知 从未上线过
	DeviceStatusRunning DeviceStatus = 0b0001 // 运行 心跳正常
	DeviceStatusStopped DeviceStatus = 0b0010 // 停止 断开连接
	DeviceStatusMalfunc DeviceStatus = 0b0100 // 故障 心跳停顿
	DeviceStatusOffline DeviceStatus = 0b1000 // 脱机 设备主动断开（提前说明）
)
