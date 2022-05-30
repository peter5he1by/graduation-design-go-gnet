package main

import (
	"flag"
	"go-gnet/database"
	"go-gnet/server"
	"go-gnet/util"
)

func ParseFlags() {
	// 数据库部分
	flag.StringVar(&database.RedisAddr, "redis-addr", "", "Redis连接地址，格式如：主机:端口")
	flag.IntVar(&database.RedisDb, "redis-db", 0, "Redis数据库编号，默认0")
	flag.StringVar(&database.RedisPassword, "redis-pass", "", "Redis密码")
	flag.StringVar(&database.MySQLAddr, "mysql-addr", "", "MySQL连接地址，格式如：主机:端口")
	flag.StringVar(&database.MySQLDb, "mysql-db", "", "MySQL数据库名")
	flag.StringVar(&database.MySQLUsername, "mysql-user", "root", "MySQL用户名，默认root")
	flag.StringVar(&database.MySQLPassword, "mysql-pass", "", "MySQL密码")
	// gnet服务器部分
	flag.StringVar(&server.Addr, "addr", "", "服务器监听地址，格式如：主机:端口")
	flag.Parse()
}

func main() {
	/* 初始化顺序是需要注意的 */
	// 解析参数
	ParseFlags()
	// 工具初始化（日志
	util.InitUtil()
	// 初始化数据库
	database.InitDatabase()
	// 启动tcp服务器
	server.Start()
}
