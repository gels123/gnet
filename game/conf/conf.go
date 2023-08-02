package conf

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
	FilePath string // 日志路径(绝对路径如"/data/", 相对路径如"./log")
	FileName string // 日志名称
	MaxSize  int64  // 最大文件大小(字节B)
}

// 相关配置
var (
	// 是否debug模式
	Debug = true
	// 是否开启后门
	BackDoor = true
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
		FilePath: "./log",
		FileName: "gamelog",
		MaxSize:  512000000,
	}
)
