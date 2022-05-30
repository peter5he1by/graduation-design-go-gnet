package protocol

type Packet struct {
	Version    uint16
	BodyLength uint32
	Body       []byte
}

const (
	DefaultHeadLength             = 6      // bytes
	DefaultProtocolVersion uint16 = 0x0100 // default protocol version
)
