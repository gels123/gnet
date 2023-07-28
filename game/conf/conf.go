package conf

import (
	"gnet/lib/log"
)

// redis配置定义
type stRedisConf struct {
	Host string
	Port int
	Db   int
	Auth string
	Inst int
}

// 日志配置定义
type stLogsConf struct {
	FilePath   string
	FileName   string
	FileLevel  int
	ShellLevel int
	MaxLine    int
	BufSize    int
}

// 相关配置
var (
	// redis配置
	RedisConf = stRedisConf{
		Host: "172.16.10.200",
		Port: 6379,
		Db:   1,
		Auth: "",
		Inst: 8,
	}
	// 日志配置
	LogsConf = stLogsConf{
		FilePath:   "../../",
		FileName:   "game",
		FileLevel:  log.DEBUG_LEVEL,
		ShellLevel: log.DEBUG_LEVEL,
		MaxLine:    1000000,
		BufSize:    5000,
	}
)
