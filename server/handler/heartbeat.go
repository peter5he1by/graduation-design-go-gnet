package handler

import (
	"github.com/panjf2000/gnet"
	"go-gnet/database"
	"go-gnet/database/mysql/model"
	"go-gnet/server/protocol"
	"time"
)

func HeartbeatHandler(ctx *model.Device, data *protocol.PacketDataReceived, c gnet.Conn) error {
	if ctx == nil {
		return nil
	}
	id := ctx.ID
	t := time.Now()
	err := database.Redis.SetDeviceHeartbeat(id, &t)
	if err != nil {
		return err
	}
	return nil
}
