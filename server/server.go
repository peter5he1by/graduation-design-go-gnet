package server

import (
	"encoding/json"
	_ "github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pkg/pool/goroutine"
	log "github.com/sirupsen/logrus"
	"go-gnet/constant"
	"go-gnet/database"
	"go-gnet/database/mysql/model"
	"go-gnet/server/pool"
	"go-gnet/server/protocol"
	"go-gnet/server_grpc"
	"time"
)

var (
	EdgeLoginCheckDelay uint
)

type CustomServer struct {
	*gnet.EventServer
	Addr       string
	Multicore  bool
	Codec      gnet.ICodec
	WorkerPool *goroutine.Pool
}

func (server *CustomServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Infof(
		"gnet server is listening on %s (multi-cores: %t, loops: %d)",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop,
	)
	go keepRefreshDeviceStatus()
	go server_grpc.StartGrpcServer(9001)
	return
}

func (server *CustomServer) React(payload []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	// !!! store customize protocol header
	// item := protocol.PicqPacket{Version: protocol.DefaultProtocolVersion}
	// c.SetContext(item)

	// 解析Context
	_ctx := c.Context()
	var ctx *model.Device
	if _ctx != nil {
		ctx = _ctx.(*model.Device)
	} else {
		ctx = nil
	}

	// 解析请求的data
	data := protocol.PacketDataReceived{}
	err := json.Unmarshal(payload, &data)
	if err != nil {
		log.Error(err)
		return
	}

	// 调用处理器，异步响应
	if handler, ok := Router[data.Code]; ok {
		_ = server.WorkerPool.Submit(func() {
			// log.Debug("↓ ========================================================")
			err := handler(ctx, &data, c)
			if err != nil {
				log.Error(err)
			}
			// log.Debug("↑ ========================================================")
		})
	} else {
		log.WithFields(
			log.Fields{"code": data.Code, "msg": data.Msg},
		).Warn("unhandled message.")
	}

	return
}

func (server CustomServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.WithField("remote", c.RemoteAddr()).Info("new connection established")
	checkContext := func() {
		// 轮询
		mutex := 5000 // 单位ms
		for mutex > 0 {
			if c != nil && c.Context() != nil {
				return
			}
			mutex -= 200
			time.Sleep(200 * time.Millisecond)
		}
		// 超过5s未登录
		err := c.Close()
		if err != nil {
			log.WithField("error", err).Error("an error occurred while closing a connection")
		}
	}
	go checkContext() // 监测连接信息
	return
}

func (server CustomServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	if c.Context() == nil {
		return
	}
	// 断开连接之后
	pool.PoolLock.Lock()
	id := c.Context().(*model.Device).ID
	delete(pool.LoginConnectionPool, id)
	// 非主动断开连接，即为停止状态
	if pool.CachedDevicesStatus[id] != constant.DeviceStatusOffline {
		pool.CachedDevicesStatus[id] = constant.DeviceStatusStopped
		pool.SaveStatusChanged(id, constant.DeviceStatusStopped)
	} else {
		// 正常下线
		delete(pool.CachedDevicesStatus, id)
		pool.SaveStatusChanged(id, constant.DeviceStatusOffline)
	}
	pool.PoolLock.Unlock()
	log.WithField("remote", c.RemoteAddr()).Warning("connection closed")
	return
}

func keepRefreshDeviceStatus() {
	log.Info("goroutine: monitoring device status")
	for {
		for id := range pool.CachedDevicesStatus {
			// TCP连接存在，设备已登录
			pool.PoolLock.Lock()
			if c, ok := pool.LoginConnectionPool[id]; ok && c.Conn.Context() != nil {
				heartbeat := database.Redis.GetDeviceLastHeartbeat(id)
				// 心跳记录存在且距今超过5s，或心跳记录不存在，且距离连接时间超过5s
				if (time.Now().Sub(c.When) > 5*time.Second) && (heartbeat == nil || (heartbeat != nil && time.Now().Sub(*heartbeat) > 5*time.Second)) {
					// 记录状态变更
					if pool.CachedDevicesStatus[id] != constant.DeviceStatusMalfunc {
						pool.SaveStatusChanged(id, constant.DeviceStatusMalfunc)
					}
					pool.CachedDevicesStatus[id] = constant.DeviceStatusMalfunc
				} else {
					// 心跳正常，在可接受范围内
					if pool.CachedDevicesStatus[id] != constant.DeviceStatusRunning {
						pool.SaveStatusChanged(id, constant.DeviceStatusRunning)
					}
					pool.CachedDevicesStatus[id] = constant.DeviceStatusRunning
				}
			} else {
				// TCP未连接
				if pool.CachedDevicesStatus[id] != constant.DeviceStatusStopped {
					// 脱机状态，连接不存在，或连接上下文不包含实际数据
					pool.CachedDevicesStatus[id] = constant.DeviceStatusOffline
				}
				// 停止状态，保留
			}
			pool.PoolLock.Unlock()
		}
		time.Sleep(1 * time.Second)
	}
}
