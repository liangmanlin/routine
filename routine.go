package routine

import "unsafe"

func Set(ptr unsafe.Pointer) {
	th := currentThread(true)
	th.val = ptr
}

func Get[T any]() *T {
	th := currentThread(true)
	return (*T)(th.val)
}

func GetPtr() unsafe.Pointer {
	th := currentThread(true)
	return th.val
}
