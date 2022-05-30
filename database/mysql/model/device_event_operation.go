package model

type DeviceEventOperation struct {
	EventID       uint   `json:"eventId"`
	Username      string `json:"username"`
	OperationType string `json:"operationType"`
	Detail        string `json:"detail"`
}

func (receiver DeviceEventOperation) TableName() string {
	return "device_event_operation"
}
