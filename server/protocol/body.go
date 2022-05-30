package protocol

import (
	"encoding/json"
	"github.com/panjf2000/gnet"
	log "github.com/sirupsen/logrus"
)

type PacketDataReceived struct {
	Code uint16 `json:"code"`
	Msg  string `json:"msg"`
	Data []byte `json:"data"` // 方便服务端读取
}

type PacketDataSent struct {
	Code uint16 `json:"code"`
	Msg  string `json:"msg"`
	Data []byte `json:"data"` // js端无所谓
}

func MakeResponse(code uint16, msg string, data interface{}) []byte {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	packet := PacketDataSent{
		Code: code,
		Msg:  msg,
		Data: marshaledData,
	}
	ret, err := json.Marshal(packet)
	if err != nil {
		return nil
	}
	return ret
}

func WriteResponse(c gnet.Conn, data *PacketDataSent) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		return err
	}
	err = c.AsyncWrite(bytes)
	if err != nil {
		log.Error(err)
		return err
	}
	return err
}
