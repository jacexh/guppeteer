package guppeteer

import (
	"errors"
	"sync"
	"sync/atomic"
)

type (
	Callback struct {
		called int32
		result chan []byte
		err    chan error
	}

	CallbackPool struct {
		pool sync.Pool
	}
)

var (
	ErrNotConsumed = errors.New("not consumed yet")
	ErrNotNotified = errors.New("not notified yet")
	ErrNotified    = errors.New("notified")
)

func NewCallback() *Callback {
	return &Callback{
		result: make(chan []byte, 1),
		err:    make(chan error, 1),
	}
}

func (cb *Callback) SetResult(d []byte) error {
	if atomic.LoadInt32(&cb.called) > 0 {
		return ErrNotified
	}
	atomic.AddInt32(&cb.called, 1)
	cb.result <- d
	return nil
}

func (cb *Callback) SetError(e error) error {
	if atomic.LoadInt32(&cb.called) > 0 {
		return ErrNotified
	}
	atomic.AddInt32(&cb.called, 1)
	cb.err <- e
	return nil
}

func (cb *Callback) Destroy() error {
	if atomic.LoadInt32(&cb.called) == 0 {
		return ErrNotNotified
	}
	if len(cb.result) > 0 || len(cb.err) > 0 {
		return ErrNotConsumed
	}
	close(cb.result)
	close(cb.err)
	return nil
}

func (cb *Callback) Reset() error {
	if atomic.LoadInt32(&cb.called) == 0 {
		return ErrNotNotified
	}
	if len(cb.result) > 0 || len(cb.err) > 0 {
		return ErrNotConsumed
	}
	atomic.StoreInt32(&cb.called, 0)
	return nil
}

func (cb *Callback) WaitResult() <-chan []byte {
	return cb.result
}

func (cb *Callback) WaitError() <-chan error {
	return cb.err
}

func NewCallbackPool() *CallbackPool {
	return &CallbackPool{pool: sync.Pool{New: func() interface{} { return NewCallback() }}}
}

func (cp *CallbackPool) Get() *Callback {
	return cp.pool.Get().(*Callback)
}

func (cp *CallbackPool) Put(cb *Callback) error {
	if err := cb.Reset(); err != nil {
		return err
	}
	cp.pool.Put(cb)
	return nil
}
