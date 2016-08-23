package conf

type LogConf struct {
	MaxFileCnt     int
	MaxFileSize    uint64
	FileNamePrefix string
	Dir            string
	Lvl            int
}