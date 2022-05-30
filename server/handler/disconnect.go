package handler

import (
	"github.com/panjf2000/gnet"
	log "github.com/sirupsen/logrus"
	"go-gnet/constant"
	"go-gnet/database/mysql/model"
	"go-gnet/server/pool"
	"go-gnet/server/protocol"
)

func Disconnect(ctx *model.Device, packet *protocol.PacketDataReceived, c gnet.Conn) error {
	if ctx == nil {
		return nil
	}
	log.WithField("uuid", ctx.Uuid).Info("device request disconnection")
	pool.PoolLock.Lock()
	pool.CachedDevicesStatus[ctx.ID] = constant.DeviceStatusOffline
	pool.PoolLock.Unlock()
	pool.SaveStatusChanged(ctx.ID, constant.DeviceStatusOffline)
	err := c.Close()
	if err != nil {
		return err
	}
	return nil
}
