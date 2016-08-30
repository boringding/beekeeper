package conf

type SrvConf struct {
	Name                string
	Host                string
	Port                int
	KeepAlive           bool
	KeepAliveSeconds    int64
	ReadTimeoutSeconds  int64
	WriteTimeoutSeconds int64
	MaxHeaderBytes      int
}
