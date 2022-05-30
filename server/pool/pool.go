package pool

import (
	"github.com/panjf2000/gnet"
	log "github.com/sirupsen/logrus"
	"go-gnet/constant"
	"go-gnet/database"
	"go-gnet/database/mysql/model"
	"sync"
	"time"
)

type ConnectionHandle struct {
	Conn gnet.Conn
	When time.Time
}

var PoolLock sync.Mutex

var LoginConnectionPool = make(map[uint]*ConnectionHandle)
var CachedDevicesStatus = make(map[uint]constant.DeviceStatus)

var WaitingForConfigLock sync.Mutex
var WaitingForConfig = make(map[uint]uint)

func SaveStatusChanged(id uint, status constant.DeviceStatus) {
	// 记录状态变更
	// if pool.CachedDevicesStatus[id] != status {
	err := database.Mysql.InsertDeviceStatusChange(id, &model.DeviceEventStatusChange{
		Status: status,
	})
	if err != nil {
		log.Error("An error occurred while saving device status change:", err)
	}
	// }
}
