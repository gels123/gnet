package conf

import (
	"gnet/lib/log"
)

// redis配置定义
type RedisConf struct {
	host string
	port int
	db   int
	auth string
	inst int
}

var (
	localRedisHost string = "172.16.10.200"
	localRedisPort int    = 6379
	localRedisAuth string = "Tuten123_"
)

// 日志配置定义
type LogsConf struct {
	path       string
	name       string
	fileLv     int
	shellLv    int
	maxLine    int
	bufferSize int
}

var (
	LogFilePath   string = "./log"         //log's file path
	LogFileName   string = "game"          //log's file name
	LogFileLevel  int    = log.LEVEL_MAX   //log's file level
	LogShellLevel int    = log.DEBUG_LEVEL //log's shell level
	LogMaxLine    int    = 10000           //log's max line per file
	LogBufferSize int    = 2000            //log's max buffer size
)

// 配置
var (
	redisConf RedisConf = RedisConf{
		host: "172.16.10.200",
		port: 6379,
		db:   1,
		auth: "",
		inst: 8,
	}
)
