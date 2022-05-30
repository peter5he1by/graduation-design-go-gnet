package handler

import (
	"encoding/json"
	"errors"
	"github.com/panjf2000/gnet"
	log "github.com/sirupsen/logrus"
	"go-gnet/constant"
	"go-gnet/database"
	"go-gnet/database/mysql/model"
	"go-gnet/server/pool"
	"go-gnet/server/protocol"
	"strings"
	"time"
)

type AuthPackage struct {
	Uuid      string `json:"uuid"`
	SecretKey string `json:"secretKey"`
}

// 设备登录
func Connect(ctx *model.Device, packet *protocol.PacketDataReceived, c gnet.Conn) error {
	if ctx != nil {
		return errors.New("duplicated login")
		// return nil
	}
	loginInfo := AuthPackage{}
	err := json.Unmarshal(packet.Data, &loginInfo)
	if err != nil {
		return err
	}
	existedDevice, err := database.Mysql.SelectDeviceByUuid(loginInfo.Uuid)
	if err != nil {
		return err
	}
	if existedDevice == nil || strings.Compare(loginInfo.SecretKey, existedDevice.SecretKey) != 0 {
		log.WithField("uuid", loginInfo.Uuid).Warning("Login failed.")
		return nil
	}
	pool.PoolLock.Lock()
	if pool.LoginConnectionPool[existedDevice.ID] == nil {
		// 登录成功
		c.SetContext(existedDevice)
		// 添加到连接池
		pool.LoginConnectionPool[existedDevice.ID] = &pool.ConnectionHandle{Conn: c, When: time.Now()}
		// 缓存设备状态
		pool.CachedDevicesStatus[existedDevice.ID] = constant.DeviceStatusRunning
		// 记录状态变更
		pool.SaveStatusChanged(existedDevice.ID, constant.DeviceStatusRunning)
		log.WithField("uuid", existedDevice.Uuid).Infof("new device logged in")
	}
	pool.PoolLock.Unlock()
	return nil
}
