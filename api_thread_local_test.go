package routine

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
)

const (
	concurrency = 1000
	loopTimes   = 200
)

func TestSupplier(t *testing.T) {
	var supplier Supplier[string]
	supplier = func() string {
		return "Hello"
	}
	assert.Equal(t, "Hello", supplier())
	//
	var fun func() string
	fun = supplier
	assert.Equal(t, "Hello", fun())
}

//===

func TestThreadLocal_New(t *testing.T) {
	tls := NewThreadLocal[string]()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2 := NewThreadLocal[int]()
	assert.Equal(t, "Hello", tls.Get())
	tls2.Set(22)
	assert.Equal(t, 22, tls2.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	fea := GoWait(func() {
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, 0, tls2.Get())
	})
	fea.Get()
}

func TestThreadLocal_Multi(t *testing.T) {
	tls := NewThreadLocal[string]()
	tls2 := NewThreadLocal[int]()
	tls.Set("Hello")
	tls2.Set(22)
	assert.Equal(t, 22, tls2.Get())
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	fea := GoWait(func() {
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, 0, tls2.Get())
	})
	fea.Get()
}

func TestThreadLocal_Concurrency(t *testing.T) {
	tls := NewThreadLocal[uint64]()
	tls2 := NewThreadLocal[uint64]()
	//
	tls2.Set(33)
	assert.Equal(t, uint64(33), tls2.Get())
	//
	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		assert.Equal(t, uint64(0), tls.Get())
		assert.Equal(t, uint64(33), tls2.Get())
		Go(func() {
			assert.Equal(t, uint64(0), tls.Get())
			assert.Equal(t, uint64(0), tls2.Get())
			v := rand.Uint64()
			v2 := rand.Uint64()
			for j := 0; j < loopTimes; j++ {
				tls.Set(v)
				tmp := tls.Get()
				assert.Equal(t, v, tmp)
				//
				tls2.Set(v2)
				tmp2 := tls2.Get()
				assert.Equal(t, v2, tmp2)
			}
			wg.Done()
		})
	}
	wg.Wait()
	//
	fea := GoWait(func() {
		assert.Equal(t, uint64(0), tls.Get())
		assert.Equal(t, uint64(0), tls2.Get())
	})
	fea.Get()
}

//===

func TestThreadLocalWithInitial_New(t *testing.T) {
	tls := NewThreadLocalWithInitial[string](func() string {
		return "Hello"
	})
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2 := NewThreadLocalWithInitial[int](func() int {
		return 22
	})
	assert.Equal(t, "Hello", tls.Get())
	assert.Equal(t, 22, tls2.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	fea := GoWait(func() {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 22, tls2.Get())
	})
	fea.Get()
}

func TestThreadLocalWithInitial_Multi(t *testing.T) {
	tls := NewThreadLocalWithInitial[string](func() string {
		return "Hello"
	})
	tls2 := NewThreadLocalWithInitial[int](func() int {
		return 22
	})
	tls.Set("Hello")
	tls2.Set(22)
	assert.Equal(t, 22, tls2.Get())
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	fea := GoWait(func() {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 22, tls2.Get())
	})
	fea.Get()
}

func TestThreadLocalWithInitial_Concurrency(t *testing.T) {
	tls := NewThreadLocalWithInitial[Any](func() Any {
		return "Hello"
	})
	tls2 := NewThreadLocalWithInitial[uint64](func() uint64 {
		return uint64(22)
	})
	//
	tls2.Set(33)
	assert.Equal(t, uint64(33), tls2.Get())
	//
	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, uint64(33), tls2.Get())
		Go(func() {
			assert.Equal(t, "Hello", tls.Get())
			assert.Equal(t, uint64(22), tls2.Get())
			v := rand.Uint64()
			v2 := rand.Uint64()
			for j := 0; j < loopTimes; j++ {
				tls.Set(v)
				tmp := tls.Get()
				assert.Equal(t, v, tmp.(uint64))
				//
				tls2.Set(v2)
				tmp2 := tls2.Get()
				assert.Equal(t, v2, tmp2)
			}
			wg.Done()
		})
	}
	wg.Wait()
	//
	fea := GoWait(func() {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, uint64(22), tls2.Get())
	})
	fea.Get()
}

//===

func TestInheritableThreadLocal_New(t *testing.T) {
	tls := NewInheritableThreadLocal[string]()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2 := NewInheritableThreadLocal[int]()
	assert.Equal(t, "Hello", tls.Get())
	tls2.Set(22)
	assert.Equal(t, 22, tls2.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	fea := GoWait(func() {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 33, tls2.Get())
	})
	fea.Get()
}

func TestInheritableThreadLocal_Multi(t *testing.T) {
	tls := NewInheritableThreadLocal[string]()
	tls2 := NewInheritableThreadLocal[int]()
	tls.Set("Hello")
	tls2.Set(22)
	assert.Equal(t, 22, tls2.Get())
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	fea := GoWait(func() {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 33, tls2.Get())
	})
	fea.Get()
}

func TestInheritableThreadLocal_Concurrency(t *testing.T) {
	tls := NewInheritableThreadLocal[uint64]()
	tls2 := NewInheritableThreadLocal[uint64]()
	//
	tls2.Set(33)
	assert.Equal(t, uint64(33), tls2.Get())
	//
	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		assert.Equal(t, uint64(0), tls.Get())
		assert.Equal(t, uint64(33), tls2.Get())
		Go(func() {
			assert.Equal(t, uint64(0), tls.Get())
			assert.Equal(t, uint64(33), tls2.Get())
			v := rand.Uint64()
			v2 := rand.Uint64()
			for j := 0; j < loopTimes; j++ {
				tls.Set(v)
				tmp := tls.Get()
				assert.Equal(t, v, tmp)
				//
				tls2.Set(v2)
				tmp2 := tls2.Get()
				assert.Equal(t, v2, tmp2)
			}
			wg.Done()
		})
	}
	wg.Wait()
	//
	fea := GoWait(func() {
		assert.Equal(t, uint64(0), tls.Get())
		assert.Equal(t, uint64(33), tls2.Get())
	})
	fea.Get()
}

//===

func TestInheritableThreadLocalWithInitial_New(t *testing.T) {
	tls := NewInheritableThreadLocalWithInitial[string](func() string {
		return "Hello"
	})
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2 := NewInheritableThreadLocalWithInitial[int](func() int {
		return 22
	})
	assert.Equal(t, "Hello", tls.Get())
	assert.Equal(t, 22, tls2.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	fea := GoWait(func() {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 33, tls2.Get())
	})
	fea.Get()
}

func TestInheritableThreadLocalWithInitial_Multi(t *testing.T) {
	tls := NewInheritableThreadLocalWithInitial[string](func() string {
		return "Hello"
	})
	tls2 := NewInheritableThreadLocalWithInitial[int](func() int {
		return 22
	})
	tls.Set("Hello")
	tls2.Set(22)
	assert.Equal(t, 22, tls2.Get())
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	fea := GoWait(func() {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 33, tls2.Get())
	})
	fea.Get()
}

func TestInheritableThreadLocalWithInitial_Concurrency(t *testing.T) {
	tls := NewInheritableThreadLocalWithInitial[Any](func() Any {
		return "Hello"
	})
	tls2 := NewInheritableThreadLocalWithInitial[uint64](func() uint64 {
		return uint64(22)
	})
	//
	tls2.Set(33)
	assert.Equal(t, uint64(33), tls2.Get())
	//
	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, uint64(33), tls2.Get())
		Go(func() {
			assert.Equal(t, "Hello", tls.Get())
			assert.Equal(t, uint64(33), tls2.Get())
			v := rand.Uint64()
			v2 := rand.Uint64()
			for j := 0; j < loopTimes; j++ {
				tls.Set(v)
				tmp := tls.Get()
				assert.Equal(t, v, tmp.(uint64))
				//
				tls2.Set(v2)
				tmp2 := tls2.Get()
				assert.Equal(t, v2, tmp2)
			}
			wg.Done()
		})
	}
	wg.Wait()
	//
	fea := GoWait(func() {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, uint64(33), tls2.Get())
	})
	fea.Get()
}

//===

// BenchmarkThreadLocal-4                          16088140                74.48 ns/op            7 B/op          0 allocs/op
func BenchmarkThreadLocal(b *testing.B) {
	tlsCount := 100
	tlsSlice := make([]ThreadLocal[int], tlsCount)
	for i := 0; i < tlsCount; i++ {
		tlsSlice[i] = NewThreadLocal[int]()
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % tlsCount
		tls := tlsSlice[index]
		initValue := tls.Get()
		if initValue != 0 {
			b.Fail()
		}
		tls.Set(i)
		if tls.Get() != i {
			b.Fail()
		}
		tls.Remove()
	}
}

// BenchmarkThreadLocalWithInitial-4               15618451                77.03 ns/op            7 B/op          0 allocs/op
func BenchmarkThreadLocalWithInitial(b *testing.B) {
	tlsCount := 100
	tlsSlice := make([]ThreadLocal[int], tlsCount)
	for i := 0; i < tlsCount; i++ {
		index := i
		tlsSlice[i] = NewThreadLocalWithInitial[int](func() int {
			return index
		})
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % tlsCount
		tls := tlsSlice[index]
		initValue := tls.Get()
		if initValue != index {
			b.Fail()
		}
		tls.Set(i)
		if tls.Get() != i {
			b.Fail()
		}
		tls.Remove()
	}
}

// BenchmarkInheritableThreadLocal-4               16109587                73.17 ns/op            7 B/op          0 allocs/op
func BenchmarkInheritableThreadLocal(b *testing.B) {
	tlsCount := 100
	tlsSlice := make([]ThreadLocal[int], tlsCount)
	for i := 0; i < tlsCount; i++ {
		tlsSlice[i] = NewInheritableThreadLocal[int]()
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % tlsCount
		tls := tlsSlice[index]
		initValue := tls.Get()
		if initValue != 0 {
			b.Fail()
		}
		tls.Set(i)
		if tls.Get() != i {
			b.Fail()
		}
		tls.Remove()
	}
}

// BenchmarkInheritableThreadLocalWithInitial-4    14862778                78.77 ns/op            7 B/op          0 allocs/op
func BenchmarkInheritableThreadLocalWithInitial(b *testing.B) {
	tlsCount := 100
	tlsSlice := make([]ThreadLocal[int], tlsCount)
	for i := 0; i < tlsCount; i++ {
		index := i
		tlsSlice[i] = NewInheritableThreadLocalWithInitial[int](func() int {
			return index
		})
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % tlsCount
		tls := tlsSlice[index]
		initValue := tls.Get()
		if initValue != index {
			b.Fail()
		}
		tls.Set(i)
		if tls.Get() != i {
			b.Fail()
		}
		tls.Remove()
	}
}
