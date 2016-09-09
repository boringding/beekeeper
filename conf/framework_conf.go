package conf

type LogConf struct {
	MaxFileCnt     int
	MaxFileSize    uint64
	FileNamePrefix string
	Dir            string
	Lvl            string
}

type SrvConf struct {
	Name                   string
	Host                   string
	Port                   int
	KeepAlive              bool
	KeepAliveSeconds       int64
	ReadTimeoutSeconds     int64
	WriteTimeoutSeconds    int64
	MaxHeaderBytes         int
	ShutdownTimeoutSeconds int
}

type FrameworkConf struct {
	LogConf LogConf
	SrvConf SrvConf
}
