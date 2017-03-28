package traffic

import (
	"fmt"
	"testing"
	"time"
)

func Test_TokenBucket(t *testing.T) {
	bucket := NewTokenBucket(1, 10*time.Second, 10*time.Second)

	go func() {
		time.Sleep(3 * time.Second)
		bucket.Put()
		fmt.Println("put a token")
	}()

	bucket.Get()
	fmt.Println("get a token")
}
