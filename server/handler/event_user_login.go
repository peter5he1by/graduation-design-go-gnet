package handler

import (
	"encoding/json"
	"github.com/panjf2000/gnet"
	log "github.com/sirupsen/logrus"
	"go-gnet/database"
	"go-gnet/database/mysql/model"
	"go-gnet/server/protocol"
)

func UserLoginEvent(ctx *model.Device, packet *protocol.PacketDataReceived, c gnet.Conn) error {
	if ctx == nil {
		return nil
	}
	event := model.DeviceEventLogin{}
	err := json.Unmarshal(packet.Data, &event)
	if err != nil {
		return err
	}
	err = database.Mysql.InsertDeviceEventLogin(ctx.ID, &event)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"username": event.Username,
		"purpose":  event.LoginPurpose,
	}).Debugf("device %s user login", ctx.Uuid)
	return nil
}
