package handler

import (
	"encoding/json"
	"github.com/panjf2000/gnet"
	log "github.com/sirupsen/logrus"
	"go-gnet/database"
	"go-gnet/database/mysql/model"
	"go-gnet/server/protocol"
)

func UpdateDeviceInfo(ctx *model.Device, packet *protocol.PacketDataReceived, c gnet.Conn) error {
	if ctx == nil {
		return nil
	}
	deviceInfo := model.DeviceInfo{}
	err := json.Unmarshal(packet.Data, &deviceInfo)
	if err != nil {
		return err
	}
	// insert
	deviceInfo.DeviceID = ctx.ID
	err = database.Mysql.InsertDeviceInfo(&deviceInfo)
	if err != nil {
		return err
	}
	log.Debugf("Device %s updated device info.", ctx.Uuid)
	return nil
}
