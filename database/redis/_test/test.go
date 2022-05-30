package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go-gnet/database/redis"
	"go-gnet/util"
	"time"
)

func main() {
	util.InitLogger()
	h := redis.HighFreqStorage{RedisDb: redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})}
	err := h.CheckConnection()
	if err != nil {
		panic(err)
	}
	// write
	t := time.Now()
	err = h.SetDeviceHeartbeat("test-uuid", &t)
	if err != nil {
		panic(err)
	}
	// read
	t2 := h.GetDeviceLastHeartbeat("test-uuid")
	logrus.Info(t2)
}
