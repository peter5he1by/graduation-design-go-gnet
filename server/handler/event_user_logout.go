package handler

import (
	"encoding/json"
	"github.com/panjf2000/gnet"
	log "github.com/sirupsen/logrus"
	"go-gnet/database"
	"go-gnet/database/mysql/model"
	"go-gnet/server/protocol"
)

func UserLogoutEvent(ctx *model.Device, packet *protocol.PacketDataReceived, c gnet.Conn) error {
	if ctx == nil {
		return nil
	}
	event := model.DeviceEventLogout{}
	err := json.Unmarshal(packet.Data, &event)
	if err != nil {
		return err
	}
	err = database.Mysql.InsertDeviceEventLogout(ctx.ID, &event)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"username":     event.Username,
		"workContents": event.WorkContents,
	}).Debugf("device %s user logout", ctx.Uuid)
	return nil
}
