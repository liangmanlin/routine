package routine

import "sync/atomic"

var inheritableThreadLocalIndex int32 = -1

func nextInheritableThreadLocalId() int {
	index := atomic.AddInt32(&inheritableThreadLocalIndex, 1)
	if index < 0 {
		panic("too many inheritable-thread-local indexed variables")
	}
	return int(index)
}

type inheritableThreadLocal[T Any] struct {
	id       int
	supplier Supplier[T]
}

func (tls *inheritableThreadLocal[T]) Get() T {
	t := currentThread(true)
	mp := tls.getMap(t)
	if mp != nil {
		v := mp.get(tls.id)
		if v != unset {
			return v.(T)
		}
	}
	return tls.setInitialValue(t)
}

func (tls *inheritableThreadLocal[T]) Set(value T) {
	t := currentThread(true)
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls.id, value)
	} else {
		tls.createMap(t, value)
	}
}

func (tls *inheritableThreadLocal[T]) Remove() {
	t := currentThread(false)
	if t == nil {
		return
	}
	mp := tls.getMap(t)
	if mp != nil {
		mp.remove(tls.id)
	}
}

func (tls *inheritableThreadLocal[T]) getMap(t *thread) *threadLocalMap {
	return t.inheritableThreadLocals
}

func (tls *inheritableThreadLocal[T]) createMap(t *thread, firstValue T) {
	mp := &threadLocalMap{}
	mp.set(tls.id, firstValue)
	t.inheritableThreadLocals = mp
}

func (tls *inheritableThreadLocal[T]) setInitialValue(t *thread) T {
	value := tls.initialValue()
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls.id, value)
	} else {
		tls.createMap(t, value)
	}
	return value
}

func (tls *inheritableThreadLocal[T]) initialValue() T {
	if tls.supplier == nil {
		var defaultValue T
		return defaultValue
	}
	return tls.supplier()
}
