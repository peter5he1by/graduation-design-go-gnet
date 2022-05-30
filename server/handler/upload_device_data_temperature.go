package handler

import (
	"encoding/json"
	"github.com/panjf2000/gnet"
	"go-gnet/database"
	"go-gnet/database/mysql/model"
	"go-gnet/server/protocol"
)

func UploadDeviceDataTemperature(ctx *model.Device, packet *protocol.PacketDataReceived, c gnet.Conn) error {
	if ctx == nil {
		return nil
	}
	temperature := model.DeviceDataTemperature{}
	err := json.Unmarshal(packet.Data, &temperature)
	if err != nil {
		return err
	}
	temperature.DeviceID = ctx.ID
	err = database.Mysql.InsertDeviceDataTemperature(&temperature)
	if err != nil {
		return err
	}
	return nil
}
