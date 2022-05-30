package handler

import (
	"github.com/panjf2000/gnet"
	log "github.com/sirupsen/logrus"
	"go-gnet/database/mysql/model"
	"go-gnet/server/protocol"
)

type AdviceHandleFunc func(ctx *model.Device, data *protocol.PacketDataReceived, c gnet.Conn) error

func Advice(handlerFunc AdviceHandleFunc) AdviceHandleFunc {
	return func(ctx *model.Device, data *protocol.PacketDataReceived, c gnet.Conn) error {
		// log.Debug("↓ ----------------------------", data)
		err := handlerFunc(ctx, data, c)
		// log.Debug("↑ ============================")
		if err != nil {
			log.Error(err)
		}
		return err
	}
}
