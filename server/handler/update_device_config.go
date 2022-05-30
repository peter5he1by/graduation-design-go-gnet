package handler

import (
	"encoding/json"
	"github.com/panjf2000/gnet"
	"go-gnet/database"
	"go-gnet/database/mysql/model"
	"go-gnet/server/pool"
	"go-gnet/server/protocol"
	"strings"
	"time"
)

func UpdateDeviceConfig(ctx *model.Device, packet *protocol.PacketDataReceived, c gnet.Conn) error {
	if ctx == nil {
		return nil
	}
	deviceConfig := &model.DeviceConfig{}
	err := json.Unmarshal(packet.Data, deviceConfig)
	if err != nil {
		return err
	}
	// 对比
	oldDeviceConfig, err := database.Mysql.SelectLatestDeviceConfig()
	if err != nil {
		return err
	}
	var id uint
	if strings.Compare(oldDeviceConfig.Content, deviceConfig.Content) != 0 {
		deviceConfig.DeviceID = ctx.ID
		id, err = database.Mysql.InsertDeviceConfig(deviceConfig)
		if err != nil {
			return err
		}
	} else {
		oldDeviceConfig.UpdatedAt = time.Now()
		err = database.Mysql.UpdateDeviceConfig(oldDeviceConfig)
		if err != nil {
			return err
		}
		id = oldDeviceConfig.ID
	}
	// 检索是否存在等待返回配置的grpc调用
	pool.WaitingForConfigLock.Lock()
	if _, ok := pool.WaitingForConfig[ctx.ID]; ok {
		pool.WaitingForConfig[ctx.ID] = id
		pool.WaitingForConfigLock.Unlock()
	}
	return nil
}
