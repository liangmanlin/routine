package routine

import (
	"runtime"
	"unsafe"
)

const threadMagic = int64('r')<<48 |
	int64('o')<<40 |
	int64('u')<<32 |
	int64('t')<<24 |
	int64('i')<<16 |
	int64('n')<<8 |
	int64('e')

type thread struct {
	labels map[string]string //pprof
	magic  int64             //mark
	id     int64             //goid
	val    unsafe.Pointer
}

// currentThread returns a pointer to the currently executing goroutine's thread struct.
//go:norace
//go:nocheckptr
func currentThread(create bool) *thread {
	gp := getg()
	goid := gp.goid
	label := gp.getLabels()
	//nothing inherited
	if label == nil {
		if create {
			newt := &thread{labels: nil, magic: threadMagic, id: goid}
			gp.setLabels(unsafe.Pointer(newt))
			runtime.SetFinalizer(newt, threadFinalize)
			return newt
		}
		return nil
	}
	//inherited map then create
	t, magic, id := extractThread(gp, label)
	if magic != threadMagic {
		if create {
			mp := *(*map[string]string)(label)
			newt := &thread{labels: mp, magic: threadMagic, id: goid}
			gp.setLabels(unsafe.Pointer(newt))
			runtime.SetFinalizer(newt, threadFinalize)
			return newt
		}
		return nil
	}
	//inherited thread then recreate
	if id != goid {
		if create || t.labels != nil {
			newt := &thread{labels: t.labels, magic: threadMagic, id: goid}
			gp.setLabels(unsafe.Pointer(newt))
			runtime.SetFinalizer(newt, threadFinalize)
			return newt
		}
		gp.setLabels(nil)
		return nil
	}
	//all is ok
	return t
}

// extractThread extract thread from unsafe.Pointer and catch fault error.
//go:norace
//go:nocheckptr
func extractThread(gp g, label unsafe.Pointer) (t *thread, magic int64, id int64) {
	if thread_safe {
		old := gp.setPanicOnFault(true)
		defer func() {
			gp.setPanicOnFault(old)
			recover()
		}()
	}
	t = (*thread)(label)
	return t, t.magic, t.id
}

func threadFinalize(t *thread) {
	t.magic = 0
}
