package database

import (
	"context"
	"fmt"
	redis2 "github.com/go-redis/redis/v8"
	"go-gnet/database/mysql"
	"go-gnet/database/redis"
	mysqldriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	Mysql *mysql.Handle
	Redis *redis.Handle
)

var (
	RedisAddr     string
	RedisDb       int
	RedisPassword string
	MySQLAddr     string
	MySQLDb       string
	MySQLUsername string
	MySQLPassword string
)

func InitDatabase() {
	// 连接Redis
	Redis = &redis.Handle{DB: redis2.NewClient(&redis2.Options{
		Addr:     RedisAddr,
		DB:       RedisDb,
		Password: RedisPassword,
	})}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := Redis.DB.Ping(ctx).Result(); err != nil {
		panic(err)
	}
	// 连接MySQL
	db, err := gorm.Open(
		mysqldriver.Open(
			fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", MySQLUsername, MySQLPassword, MySQLAddr, MySQLDb),
		),
		&gorm.Config{
			PrepareStmt: true,
		},
	)
	if err != nil {
		panic(err)
	}
	_db, err := db.DB()
	if err != nil {
		panic(err)
	}
	_db.SetConnMaxLifetime(59 * time.Second)
	Mysql = &mysql.Handle{DB: db}
}
