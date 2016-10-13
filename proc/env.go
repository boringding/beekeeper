//Environment variable operations.

package proc

import (
	"os"
)

func GetEnv(key string) string {
	return os.Getenv(key)
}

func SetEnv(key string, val string) error {
	return os.Setenv(key, val)
}
