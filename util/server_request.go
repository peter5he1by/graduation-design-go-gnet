package util

import (
	"github.com/panjf2000/gnet"
	log "github.com/sirupsen/logrus"
	"go-gnet/server/protocol"
)

var codec = protocol.PacketCodec{}

func Send(conn gnet.Conn, code uint16, data interface{}) error {
	p := protocol.MakeResponse(code, "", data)
	// encode, err := codec.Encode(nil, p)
	// if err != nil {
	// 	return err
	// }
	log.Debug("server sent a packet")
	err := conn.AsyncWrite(p)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
