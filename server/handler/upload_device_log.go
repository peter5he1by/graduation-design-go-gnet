package handler

import (
	"encoding/json"
	"github.com/panjf2000/gnet"
	log "github.com/sirupsen/logrus"
	"go-gnet/database"
	"go-gnet/database/mysql/model"
	"go-gnet/server/protocol"
)

func UploadDeviceLog(ctx *model.Device, packet *protocol.PacketDataReceived, c gnet.Conn) error {
	if ctx == nil {
		return nil
	}
	deviceLog := model.DeviceLog{}
	err := json.Unmarshal(packet.Data, &deviceLog)
	if err != nil {
		return err
	}
	deviceLog.DeviceID = ctx.ID
	err = database.Mysql.InsertDeviceLog(&deviceLog)
	if err != nil {
		return err
	}
	log.Debugf("Device %s upload device log.", ctx.Uuid)
	return nil
}
