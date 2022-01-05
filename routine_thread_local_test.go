package routine

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func init() {
	gCInterval = time.Millisecond * 50 // for faster test
}

func TestThreadLocal(t *testing.T) {
	threadLocal := NewThreadLocal()
	threadLocal2 := NewThreadLocal()

	for i := 0; i < 100; i++ {
		src := "hello"
		threadLocal.Set(src)
		p := threadLocal.Get()
		assert.Equal(t, src, p.(string))
		//
		src2 := "world"
		threadLocal2.Set(src2)
		p2 := threadLocal2.Get()
		assert.Equal(t, src2, p2.(string))
	}

	for i := 0; i < 1000; i++ {
		num := rand.Int()
		threadLocal.Set(num)
		num2 := threadLocal.Get()
		assert.Equal(t, num, num2.(int))
	}

	v := threadLocal.Remove()
	assert.NotNil(t, v)

	Clear()
	vv1 := threadLocal.Get()
	assert.Nil(t, vv1)
	//
	vv2 := threadLocal2.Get()
	assert.Nil(t, vv2)
}

func TestThreadLocalConcurrency(t *testing.T) {
	const concurrency = 1000
	const loopTimes = 1000

	threadLocal := NewThreadLocal()
	threadLocal2 := NewThreadLocal()

	waiter := &sync.WaitGroup{}
	waiter.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			v := rand.Uint64()
			v2 := rand.Uint64()
			assert.True(t, v != 0)
			assert.True(t, v2 != 0)
			for i := 0; i < loopTimes; i++ {
				threadLocal.Set(v)
				tmp := threadLocal.Get()
				assert.True(t, tmp.(uint64) == v)
				//
				threadLocal2.Set(v2)
				tmp2 := threadLocal2.Get()
				assert.True(t, tmp2.(uint64) == v2)
			}
			waiter.Done()
		}()
	}
	waiter.Wait()
}

func TestThreadLocalGC(t *testing.T) {
	threadLocal := NewThreadLocal()
	threadLocal2 := NewThreadLocal()
	threadLocal3 := NewThreadLocal()
	threadLocal4 := NewThreadLocal()
	threadLocal5 := NewThreadLocal()

	// use ThreadLocal in multi goroutines
	gcRunningCnt := 0
	for i := 0; i < 10; i++ {
		waiter := &sync.WaitGroup{}
		for i := 0; i < 1000; i++ {
			waiter.Add(1)
			go func() {
				threadLocal.Set("hello world")
				threadLocal2.Set(true)
				threadLocal3.Set(&threadLocal3)
				threadLocal4.Set(rand.Int())
				threadLocal5.Set(time.Now())
				waiter.Done()
			}()
		}
		if gcRunning() {
			// record gc running count
			gcRunningCnt++
		}
		waiter.Wait()
		// wait for a while
		time.Sleep(gCInterval + time.Second)
		assert.False(t, gcRunning(), "#%v, timer not stoped?", i)
		storeMap := globalMap.Load().(map[int64]*threadLocalMap)
		assert.Equal(t, 0, len(storeMap), "#%v, storeMap not empty - %d", i, len(storeMap))
	}
	assert.True(t, gcRunningCnt > 0, "gc timer may not running!")
}

// BenchmarkThreadLocal-8   	10183446	        99.28 ns/op	       8 B/op	       0 allocs/op
func BenchmarkThreadLocal(b *testing.B) {
	threadLocalCount := 100
	threadLocals := make([]ThreadLocal, threadLocalCount)
	for i := 0; i < threadLocalCount; i++ {
		threadLocals[i] = NewThreadLocal()
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % threadLocalCount
		threadLocal := threadLocals[index]
		threadLocal.Set(i)
		if threadLocal.Get() != i {
			b.Fail()
		}
		threadLocal.Remove()
	}
}
