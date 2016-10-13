//Basic configure struct definition of the framework.
//The corresponding configure files are:
//your_project_dir/conf/framework.conf.environment_name.xml.

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

type MonConf struct {
	Enabled bool
	Host    string
	Port    int
}

type FrameworkConf struct {
	LogConf LogConf
	SrvConf SrvConf
	MonConf MonConf
}
