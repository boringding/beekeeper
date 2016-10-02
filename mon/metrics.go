package mon

import (
	"errors"
	"expvar"
	"fmt"
)

type Metrics struct {
	name string
	val  *expvar.Int
}

func NewMetrics(name string, initVal int64) (err error, m *Metrics) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
			m = nil
		}
	}()

	val := expvar.NewInt(name)
	val.Set(initVal)

	m = &Metrics{
		name: name,
		val:  val,
	}

	return
}

func (self *Metrics) Add(delta int64) error {
	if self.val == nil {
		return errors.New("empty expvar")
	}

	self.val.Add(delta)

	return nil
}

func (self *Metrics) Set(val int64) error {
	if self.val == nil {
		return errors.New("empty expvar")
	}

	self.val.Set(val)

	return nil
}
