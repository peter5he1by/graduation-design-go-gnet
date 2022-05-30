package handler

import (
	"encoding/json"
	"github.com/panjf2000/gnet"
	log "github.com/sirupsen/logrus"
	"go-gnet/database"
	"go-gnet/database/mysql/model"
	"go-gnet/server/protocol"
)

func OperationEvent(ctx *model.Device, packet *protocol.PacketDataReceived, c gnet.Conn) error {
	if ctx == nil {
		return nil
	}
	event := model.DeviceEventOperation{}
	err := json.Unmarshal(packet.Data, &event)
	if err != nil {
		return err
	}
	err = database.Mysql.InsertDeviceEventOperation(ctx.ID, &event)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"user": event.Username, "type": event.OperationType},
	).Debugf("device %s operation", ctx.Uuid)
	return nil
}
