package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/panjf2000/gnet"
	log "github.com/sirupsen/logrus"
)

type PacketCodec struct{}

// Encode ...
func (cc *PacketCodec) Encode(c gnet.Conn, data []byte) ([]byte, error) {
	result := make([]byte, 0)

	buffer := bytes.NewBuffer(result)

	// version field
	if err := binary.Write(buffer, binary.BigEndian, DefaultProtocolVersion); err != nil {
		s := fmt.Sprintf("Pack version error , %v", err)
		return nil, errors.New(s)
	}

	// body length field
	dataLen := uint32(len(data))
	if err := binary.Write(buffer, binary.BigEndian, dataLen); err != nil {
		s := fmt.Sprintf("Pack data length error , %v", err)
		return nil, errors.New(s)
	}

	// body content
	if dataLen > 0 {
		if err := binary.Write(buffer, binary.BigEndian, data); err != nil {
			s := fmt.Sprintf("Pack data error , %v", err)
			return nil, errors.New(s)
		}
	}

	return buffer.Bytes(), nil
}

// Decode ...
func (cc *PacketCodec) Decode(c gnet.Conn) ([]byte, error) {
	// parse header
	headerLen := DefaultHeadLength
	if size, header := c.ReadN(headerLen); size == headerLen {
		byteBuffer := bytes.NewBuffer(header)
		var version uint16
		var dataLength uint32
		_ = binary.Read(byteBuffer, binary.BigEndian, &version)
		_ = binary.Read(byteBuffer, binary.BigEndian, &dataLength)
		// to check the protocol version,
		// reset buffer if the version is not correct
		if version != DefaultProtocolVersion {
			c.ResetBuffer()
			log.Errorln("Packet parsing error:", header)
			return nil, errors.New("packet parsing error")
		}
		// parse payload
		dataLen := int(dataLength) // max int32 can contain 210MB payload
		totalLen := headerLen + dataLen
		if dataSize, data := c.ReadN(totalLen); dataSize == totalLen {
			c.ShiftN(totalLen)
			// log.Println("parse success:", data, dataSize)

			// return the payload of the data
			return data[headerLen:], nil
		}
		// log.Println("not enough payload data:", dataLen, totalLen, dataSize)
		return nil, errors.New("not enough payload data")
	}
	// log.Println("not enough header data:", size)
	return nil, errors.New("not enough header data")
}
