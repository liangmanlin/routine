package routine

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestRunnable(t *testing.T) {
	count := 0
	var runnable Runnable
	runnable = func() {
		count++
	}
	runnable()
	assert.Equal(t, 1, count)
	//
	var fun func()
	fun = runnable
	fun()
	assert.Equal(t, 2, count)
}

func TestCallable(t *testing.T) {
	var callable Callable[string]
	callable = func() string {
		return "Hello"
	}
	assert.Equal(t, "Hello", callable())
	//
	var fun func() string
	fun = callable
	assert.Equal(t, "Hello", fun())
}

func TestGo_Error(t *testing.T) {
	run := false
	assert.NotPanics(t, func() {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		Go(func() {
			run = true
			wg.Done()
			panic("error")
		})
		wg.Wait()
	})
	assert.True(t, run)
}

func TestGo_Nil(t *testing.T) {
	assert.Nil(t, createInheritedMap())
	//
	run := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	Go(func() {
		assert.Nil(t, createInheritedMap())
		run = true
		wg.Done()
	})
	wg.Wait()
	assert.True(t, run)
}

func TestGo_Value(t *testing.T) {
	tls := NewThreadLocal[string]()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	inheritableTls := NewInheritableThreadLocal[string]()
	inheritableTls.Set("World")
	assert.Equal(t, "World", inheritableTls.Get())
	//
	assert.NotNil(t, createInheritedMap())
	//
	run := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	Go(func() {
		assert.NotNil(t, createInheritedMap())
		//
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, "World", inheritableTls.Get())
		//
		tls.Set("Hello2")
		assert.Equal(t, "Hello2", tls.Get())
		//
		inheritableTls.Remove()
		assert.Equal(t, "", inheritableTls.Get())
		//
		run = true
		wg.Done()
	})
	wg.Wait()
	assert.True(t, run)
	//
	assert.Equal(t, "Hello", tls.Get())
	assert.Equal(t, "World", inheritableTls.Get())
}

func TestGo_Cross(t *testing.T) {
	tls := NewThreadLocal[string]()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	wg := &sync.WaitGroup{}
	wg.Add(1)
	Go(func() {
		assert.Equal(t, "", tls.Get())
		wg.Done()
	})
	wg.Wait()
}

func TestGoWait_Error(t *testing.T) {
	run := false
	assert.Panics(t, func() {
		fea := GoWait(func() {
			run = true
			panic("error")
		})
		fea.Get()
	})
	assert.True(t, run)
}

func TestGoWait_Nil(t *testing.T) {
	assert.Nil(t, createInheritedMap())
	//
	run := false
	fea := GoWait(func() {
		assert.Nil(t, createInheritedMap())
		run = true
	})
	assert.Nil(t, fea.Get())
	assert.True(t, run)
}

func TestGoWait_Value(t *testing.T) {
	tls := NewThreadLocal[string]()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	inheritableTls := NewInheritableThreadLocal[string]()
	inheritableTls.Set("World")
	assert.Equal(t, "World", inheritableTls.Get())
	//
	assert.NotNil(t, createInheritedMap())
	//
	run := false
	fea := GoWait(func() {
		assert.NotNil(t, createInheritedMap())
		//
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, "World", inheritableTls.Get())
		//
		tls.Set("Hello2")
		assert.Equal(t, "Hello2", tls.Get())
		//
		inheritableTls.Remove()
		assert.Equal(t, "", inheritableTls.Get())
		//
		run = true
	})
	assert.Nil(t, fea.Get())
	assert.True(t, run)
	//
	assert.Equal(t, "Hello", tls.Get())
	assert.Equal(t, "World", inheritableTls.Get())
}

func TestGoWait_Cross(t *testing.T) {
	tls := NewThreadLocal[string]()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	GoWait(func() {
		assert.Equal(t, "", tls.Get())
	}).Get()
}

func TestGoWaitResult_Error(t *testing.T) {
	run := false
	assert.Panics(t, func() {
		fea := GoWaitResult(func() int {
			run = true
			if run {
				panic("error")
			}
			return 1
		})
		fea.Get()
	})
	assert.True(t, run)
}

func TestGoWaitResult_Nil(t *testing.T) {
	assert.Nil(t, createInheritedMap())
	//
	run := false
	fea := GoWaitResult(func() bool {
		assert.Nil(t, createInheritedMap())
		run = true
		return true
	})
	assert.True(t, fea.Get())
	assert.True(t, run)
}

func TestGoWaitResult_Value(t *testing.T) {
	tls := NewThreadLocal[string]()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	inheritableTls := NewInheritableThreadLocal[string]()
	inheritableTls.Set("World")
	assert.Equal(t, "World", inheritableTls.Get())
	//
	assert.NotNil(t, createInheritedMap())
	//
	run := false
	fea := GoWaitResult(func() bool {
		assert.NotNil(t, createInheritedMap())
		//
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, "World", inheritableTls.Get())
		//
		tls.Set("Hello2")
		assert.Equal(t, "Hello2", tls.Get())
		//
		inheritableTls.Remove()
		assert.Equal(t, "", inheritableTls.Get())
		//
		run = true
		return true
	})
	assert.True(t, fea.Get())
	assert.True(t, run)
	//
	assert.Equal(t, "Hello", tls.Get())
	assert.Equal(t, "World", inheritableTls.Get())
}

func TestGoWaitResult_Cross(t *testing.T) {
	tls := NewThreadLocal[string]()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	result := GoWaitResult(func() string {
		assert.Equal(t, "", tls.Get())
		return tls.Get()
	}).Get()
	assert.Equal(t, "", result)
}
