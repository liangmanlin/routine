package routine

import "unsafe"

func Set(ptr unsafe.Pointer) {
	th := currentThread(true)
	th.val = ptr
}

func Get[T any]() *T {
	th := currentThread(false)
	if th == nil {
		return nil
	}
	return (*T)(th.val)
}

func GetPtr() unsafe.Pointer {
	th := currentThread(false)
	if th == nil {
		return nil
	}
	return th.val
}
