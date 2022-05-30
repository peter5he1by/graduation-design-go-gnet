package model

import "go-gnet/constant"

type DeviceEventStatusChange struct {
	EventID uint                  `json:"eventId"`
	Status  constant.DeviceStatus `json:"status"`
}

func (receiver DeviceEventStatusChange) TableName() string {
	return "device_event_status_change"
}
