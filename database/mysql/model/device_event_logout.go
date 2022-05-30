package model

type DeviceEventLogout struct {
	EventID      uint   `json:"eventId"`
	Username     string `json:"username"`
	WorkContents string `json:"workContents"`
}

func (receiver DeviceEventLogout) TableName() string {
	return "device_event_logout"
}
