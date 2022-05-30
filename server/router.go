package server

import (
	"go-gnet/constant"
	"go-gnet/server/handler"
)

var Router = map[uint16]handler.AdviceHandleFunc{
	constant.ClientHeartbeat:             handler.HeartbeatHandler,
	constant.ClientConnect:               handler.Connect,
	constant.ClientUpdateDeviceInfo:      handler.UpdateDeviceInfo,
	constant.ClientEventUserLogin:        handler.UserLoginEvent,
	constant.ClientEventUserLogout:       handler.UserLogoutEvent,
	constant.ClientEventOperation:        handler.OperationEvent,
	constant.ClientUploadDeviceLog:       handler.UploadDeviceLog,
	constant.ClientDisconnect:            handler.Disconnect,
	constant.ClientUploadDeviceConfig:    handler.UpdateDeviceConfig,
	constant.ClientDataUploadTemperature: handler.UploadDeviceDataTemperature,
}

func InitRouter() {
	for k, f := range Router {
		Router[k] = handler.Advice(f)
	}
}
