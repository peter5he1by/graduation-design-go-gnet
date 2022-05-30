package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func heartbeatKey(id uint) string {
	return fmt.Sprintf("device-heartbeat[%d]", id)
}

func (h Handle) GetDeviceLastHeartbeat(id uint) *time.Time {
	if h.CheckConnection() != nil {
		return nil
	}
	key := heartbeatKey(id)
	// 取值
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	val, err := h.DB.Get(ctx, key).Result()
	if err != redis.Nil && val != " " {
		// 反序列化
		// dec := gob.NewDecoder(bytes.NewBuffer([]byte(val)))
		// var t time.Time
		// if err := dec.Decode(&t); err == nil {
		//	return &t
		// }
		timestamp, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			log.Error(err)
			return nil
		}
		t := time.Unix(timestamp, 0)
		if err != nil {
			log.Error(err)
			return nil
		}
		return &t
	}
	return nil
}

func (h Handle) SetDeviceHeartbeat(id uint, t *time.Time) error {
	if err := h.CheckConnection(); err != nil {
		return err
	}
	key := heartbeatKey(id)
	// 序列化 time.Time 对象
	// buf := new(bytes.Buffer)
	// enc := gob.NewEncoder(buf)
	// if err := enc.Encode(t); err != nil {
	//	logrus.Error("Cannot encode time.Time object with gob !")
	// }
	// 保存到 redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := h.DB.Set(ctx, key, t.Unix(), 0).Result()
	if err != nil {
		log.Error(err)
	}
	return nil
}
