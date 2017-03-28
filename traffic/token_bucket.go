//A simple token bucket implementation.

package traffic

import (
	"errors"
	"time"
)

type TokenBucket struct {
	Capacity   int
	Bucket     chan int
	GetTimeout time.Duration
	PutTimeout time.Duration
}

func NewTokenBucket(capacity int, getTimeout time.Duration, putTimeout time.Duration) *TokenBucket {
	return &TokenBucket{
		Capacity:   capacity,
		Bucket:     make(chan int, capacity),
		GetTimeout: getTimeout,
		PutTimeout: putTimeout,
	}
}

func (self *TokenBucket) Get() error {
	timeoutChan := make(chan int, 1)

	go func(duration time.Duration, c chan int) {
		time.Sleep(duration)
		timeoutChan <- 1
	}(self.GetTimeout, timeoutChan)

	select {
	case <-self.Bucket:
		return nil
	case <-timeoutChan:
		return errors.New("get token timeout")
	}

	return nil
}

func (self *TokenBucket) Put() error {
	timeoutChan := make(chan int, 1)

	go func(duration time.Duration, c chan int) {
		time.Sleep(duration)
		timeoutChan <- 1
	}(self.PutTimeout, timeoutChan)

	select {
	case self.Bucket <- 1:
		return nil
	case <-timeoutChan:
		return errors.New("put token timeout")
	}

	return nil
}
