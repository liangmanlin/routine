package routine

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"runtime/pprof"
	"sync"
	"testing"
	"time"
	"unsafe"
)

func TestCurrentThread(t *testing.T) {
	assert.NotNil(t, currentThread(true))
	assert.Same(t, currentThread(true), currentThread(true))
}

func TestPProf(t *testing.T) {
	const concurrency = 10
	const loopTimes = 10
	val := 1
	Set(unsafe.Pointer(&val))
	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		tmp := i
		go func() {
			for j := 0; j < loopTimes; j++ {
				time.Sleep(100 * time.Millisecond)
				Set(unsafe.Pointer(&tmp))
				assert.Equal(t, tmp, *Get[int]())
				pprof.Do(context.Background(), pprof.Labels("key", "value"), func(ctx context.Context) {
					assert.Nil(t, currentThread(false))
					assert.Nil(t, Get[int]())
					Set(unsafe.Pointer(&tmp))
					//
					label, find := pprof.Label(ctx, "key")
					assert.True(t, find)
					assert.Equal(t, "value", label)
					//
					assert.Equal(t, tmp, *Get[int]())
					//
					label2, find2 := pprof.Label(ctx, "key")
					assert.True(t, find2)
					assert.Equal(t, "value", label2)
				})
				assert.Nil(t, Get[int]())
			}
			wg.Done()
		}()
	}
	assert.Nil(t, pprof.StartCPUProfile(&bytes.Buffer{}))
	wg.Wait()
	pprof.StopCPUProfile()
	assert.Equal(t, val, *Get[int]())
}

func BenchmarkGet(b *testing.B) {
	val := []int{1, 23, 4, 1, 4, 4, 5}
	Set(unsafe.Pointer(&val))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Get[[]int]()
	}
}
