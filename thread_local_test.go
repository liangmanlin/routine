package routine

import (
	"github.com/stretchr/testify/assert"
	"math"
	"sync"
	"testing"
)

func TestThreadLocalId(t *testing.T) {
	tls := NewThreadLocal[string]()
	assert.GreaterOrEqual(t, tls.(*threadLocal[string]).id, 0)
	tls2 := NewThreadLocalWithInitial[string](func() string {
		return "Hello"
	})
	assert.Greater(t, tls2.(*threadLocal[string]).id, tls.(*threadLocal[string]).id)
}

func TestThreadLocalNextId(t *testing.T) {
	backup := threadLocalIndex
	defer func() {
		threadLocalIndex = backup
	}()
	//
	threadLocalIndex = math.MaxInt32
	assert.Panics(t, func() {
		nextThreadLocalId()
	})
}

func TestThreadLocal(t *testing.T) {
	tls := NewThreadLocal[int]()
	tls2 := NewThreadLocal[string]()
	tls.Remove()
	tls2.Remove()
	assert.Equal(t, 0, tls.Get())
	assert.Equal(t, "", tls2.Get())
	//
	tls.Set(1)
	tls2.Set("World")
	assert.Equal(t, 1, tls.Get())
	assert.Equal(t, "World", tls2.Get())
	//
	tls.Set(0)
	tls2.Set("")
	assert.Equal(t, 0, tls.Get())
	assert.Equal(t, "", tls2.Get())
	//
	tls.Set(2)
	tls2.Set("!")
	assert.Equal(t, 2, tls.Get())
	assert.Equal(t, "!", tls2.Get())
	//
	tls.Remove()
	tls2.Remove()
	assert.Equal(t, 0, tls.Get())
	assert.Equal(t, "", tls2.Get())
	//
	tls.Set(2)
	tls2.Set("!")
	assert.Equal(t, 2, tls.Get())
	assert.Equal(t, "!", tls2.Get())
	wg := &sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		Go(func() {
			assert.Equal(t, 0, tls.Get())
			assert.Equal(t, "", tls2.Get())
			wg.Done()
		})
	}
	wg.Wait()
	assert.Equal(t, 2, tls.Get())
	assert.Equal(t, "!", tls2.Get())
}

func TestThreadLocalMixed(t *testing.T) {
	tls := NewThreadLocal[int]()
	tls2 := NewThreadLocalWithInitial[string](func() string {
		return "Hello"
	})
	assert.Equal(t, 0, tls.Get())
	assert.Equal(t, "Hello", tls2.Get())
	//
	tls.Set(1)
	tls2.Set("World")
	assert.Equal(t, 1, tls.Get())
	assert.Equal(t, "World", tls2.Get())
	//
	tls.Set(0)
	tls2.Set("")
	assert.Equal(t, 0, tls.Get())
	assert.Equal(t, "", tls2.Get())
	//
	tls.Set(2)
	tls2.Set("!")
	assert.Equal(t, 2, tls.Get())
	assert.Equal(t, "!", tls2.Get())
	//
	tls.Remove()
	tls2.Remove()
	assert.Equal(t, 0, tls.Get())
	assert.Equal(t, "Hello", tls2.Get())
	//
	tls.Set(2)
	tls2.Set("!")
	assert.Equal(t, 2, tls.Get())
	assert.Equal(t, "!", tls2.Get())
	wg := &sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		Go(func() {
			assert.Equal(t, 0, tls.Get())
			assert.Equal(t, "Hello", tls2.Get())
			wg.Done()
		})
	}
	wg.Wait()
	assert.Equal(t, 2, tls.Get())
	assert.Equal(t, "!", tls2.Get())
}

func TestThreadLocalWithInitial(t *testing.T) {
	src := &person{Id: 1, Name: "Tim"}
	tls := NewThreadLocalWithInitial[*person](nil)
	tls2 := NewThreadLocalWithInitial[*person](func() *person {
		var value *person
		return value
	})
	tls3 := NewThreadLocalWithInitial[*person](func() *person {
		return src
	})
	tls4 := NewThreadLocalWithInitial[person](func() person {
		return *src
	})

	for i := 0; i < 100; i++ {
		p := tls.Get()
		assert.Nil(t, p)
		//
		p2 := tls2.Get()
		assert.Nil(t, p2)
		//
		p3 := tls3.Get()
		assert.Same(t, src, p3)

		p4 := tls4.Get()
		assert.NotSame(t, src, &p4)
		assert.Equal(t, *src, p4)

		wg := &sync.WaitGroup{}
		wg.Add(1)
		Go(func() {
			assert.Same(t, src, tls3.Get())
			p5 := tls4.Get()
			assert.NotSame(t, src, &p5)
			assert.Equal(t, *src, p5)
			//
			wg.Done()
		})
		wg.Wait()
	}

	tls3.Set(nil)
	tls4.Set(person{})
	assert.Nil(t, tls3.Get())
	assert.Equal(t, person{}, tls4.Get())

	tls3.Remove()
	tls4.Remove()
	assert.Same(t, src, tls3.Get())
	p6 := tls4.Get()
	assert.NotSame(t, src, &p6)
	assert.Equal(t, *src, p6)
}

func TestThreadLocalCrossCoroutine(t *testing.T) {
	tls := NewThreadLocal[string]()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	subWait := &sync.WaitGroup{}
	subWait.Add(2)
	finishWait := &sync.WaitGroup{}
	finishWait.Add(2)
	go func() {
		subWait.Wait()
		assert.Equal(t, "", tls.Get())
		finishWait.Done()
	}()
	Go(func() {
		subWait.Wait()
		assert.Equal(t, "", tls.Get())
		finishWait.Done()
	})
	tls.Remove()      //remove in parent goroutine should not affect child goroutine
	subWait.Done()    //allow sub goroutine run
	subWait.Done()    //allow sub goroutine run
	finishWait.Wait() //wait sub goroutine done
	finishWait.Wait() //wait sub goroutine done
}

func TestThreadLocalCreateBatch(t *testing.T) {
	const count = 128
	tlsList := make([]ThreadLocal[int], count)
	for i := 0; i < count; i++ {
		value := i
		tlsList[i] = NewThreadLocalWithInitial[int](func() int { return value })
	}
	for i := 0; i < count; i++ {
		assert.Equal(t, i, tlsList[i].Get())
	}
}

type person struct {
	Id   int
	Name string
}
