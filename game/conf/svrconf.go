package conf

import "gnet/lib/log"

//本地redis配置
var (
	localRedisHost string = "192.168.88.235"
	localRedisPort int    = 6379
	localRedisAuth string = "Tuten123_"
)

//日志配置
var (
	LogFilePath   string = "./log"         //log's file path
	LogFileName   string = "game"          //log's file name
	LogFileLevel  int    = log.LEVEL_MAX   //log's file level
	LogShellLevel int    = log.DEBUG_LEVEL //log's shell level
	LogMaxLine    int    = 10000           //log's max line per file
	LogBufferSize int    = 2000            //log's max buffer size
)
