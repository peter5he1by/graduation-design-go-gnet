package server

import (
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pkg/pool/goroutine"
	"go-gnet/server/protocol"
	"time"
)

var (
	Addr string
)

func Start() {
	InitRouter()
	addr := fmt.Sprintf("tcp://%s", Addr)
	codec := &protocol.PacketCodec{}
	cs := &CustomServer{Addr: addr, Multicore: true, Codec: codec, WorkerPool: goroutine.Default()}
	err := gnet.Serve(cs, addr, gnet.WithMulticore(true), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(codec))
	if err != nil {
		panic(err)
	}
}
