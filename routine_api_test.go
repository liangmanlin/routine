package routine

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewThreadLocal(t *testing.T) {
	s := NewThreadLocal()
	s.Set("hello")
	assert.Equal(t, "hello", s.Get())
	//
	s2 := NewThreadLocal()
	assert.Equal(t, "hello", s.Get())
	s2.Set(22)
	assert.Equal(t, 22, s2.Get())
}

func TestMultiThreadLocal(t *testing.T) {
	s := NewThreadLocal()
	s2 := NewThreadLocal()
	s.Set("hello")
	s2.Set(22)
	assert.Equal(t, 22, s2.Get())
	assert.Equal(t, "hello", s.Get())
}

func TestBackupContext(t *testing.T) {
	s := NewThreadLocal()
	ic := BackupContext()

	waiter := &sync.WaitGroup{}
	waiter.Add(1)
	go func() {
		s.Set("hello")
		assert.Equal(t, "hello", s.Get())
		icLocalBackup := BackupContext()
		//
		RestoreContext(ic)
		assert.Nil(t, s.Get())
		//
		RestoreContext(icLocalBackup)
		assert.Equal(t, "hello", s.Get())
		//
		waiter.Done()
	}()
	waiter.Wait()
}

func TestGoid(t *testing.T) {
	assert.NotEqual(t, 0, Goid())
}

func TestAllGoid(t *testing.T) {
	const num = 10
	for i := 0; i < num; i++ {
		go func() {
			time.Sleep(time.Second)
		}()
	}
	time.Sleep(time.Millisecond)

	ids := AllGoids()
	t.Log("all gids: ", len(ids), ids)
}

func TestGoThreadLocal(t *testing.T) {
	waiter := &sync.WaitGroup{}
	waiter.Add(1)
	variable := "hello world"
	stg := NewThreadLocal()
	stg.Set(variable)
	Go(func() {
		v := stg.Get()
		assert.Equal(t, variable, v.(string))
		waiter.Done()
	})
	waiter.Wait()
}

// BenchmarkGoid-12    	278801190	         4.586 ns/op	       0 B/op	       0 allocs/op
func BenchmarkGoid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Goid()
	}
}

// BenchmarkAllGoid-12    	 5949680	       228.3 ns/op	     896 B/op	       1 allocs/op
func BenchmarkAllGoid(b *testing.B) {
	const num = 16
	for i := 0; i < num; i++ {
		go func() {
			time.Sleep(time.Second)
		}()
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AllGoids()
	}
}
