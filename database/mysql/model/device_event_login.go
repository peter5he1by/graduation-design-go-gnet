package model

type DeviceEventLogin struct {
	EventID      uint   `json:"eventId"`
	Username     string `json:"username"`
	LoginPurpose string `json:"loginPurpose"`
}

func (receiver DeviceEventLogin) TableName() string {
	return "device_event_login"
}
