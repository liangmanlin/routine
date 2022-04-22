package routine

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
)

func TestCreateInheritedMap(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		thd := currentThread(true)
		assert.NotNil(t, thd)
		assert.Nil(t, thd.inheritableThreadLocals)
		thd.inheritableThreadLocals = &threadLocalMap{}
		assert.Nil(t, thd.inheritableThreadLocals.table)
		assert.Nil(t, createInheritedMap())
		//
		wg.Done()
	}()
	wg.Wait()
}

func TestCreateInheritedMapNil(t *testing.T) {
	tls := NewInheritableThreadLocal[string]()
	tls.Set("")
	srcValue := tls.Get()
	assert.Equal(t, "", srcValue)
	assert.True(t, srcValue == "")

	mp := createInheritedMap()
	assert.NotNil(t, mp)
	getValue := mp.get(tls.(*inheritableThreadLocal[string]).id)
	assert.Equal(t, "", getValue)
	assert.True(t, getValue == "")

	mp2 := createInheritedMap()
	assert.NotNil(t, mp2)
	assert.NotSame(t, mp, mp2)
	getValue2 := mp2.get(tls.(*inheritableThreadLocal[string]).id)
	assert.Equal(t, "", getValue2)
	assert.True(t, getValue2 == "")
}

func TestCreateInheritedMapValue(t *testing.T) {
	tls := NewInheritableThreadLocal[uint64]()
	value := rand.Uint64()
	tls.Set(value)
	srcValue := tls.Get()
	assert.NotSame(t, &value, &srcValue)
	assert.Equal(t, value, srcValue)

	mp := createInheritedMap()
	assert.NotNil(t, mp)
	getValue := mp.get(tls.(*inheritableThreadLocal[uint64]).id)
	assert.NotSame(t, &value, &getValue)
	assert.Equal(t, value, getValue)

	mp2 := createInheritedMap()
	assert.NotNil(t, mp2)
	assert.NotSame(t, mp, mp2)
	getValue2 := mp2.get(tls.(*inheritableThreadLocal[uint64]).id)
	assert.NotSame(t, &value, &getValue2)
	assert.Equal(t, value, getValue2)
}

func TestCreateInheritedMapStruct(t *testing.T) {
	tls := NewInheritableThreadLocal[personCloneable]()
	value := personCloneable{Id: 1, Name: "Hello"}
	tls.Set(value)
	srcValue := tls.Get()
	assert.NotSame(t, &value, &srcValue)
	assert.Equal(t, value, srcValue)

	mp := createInheritedMap()
	assert.NotNil(t, mp)
	getValue := mp.get(tls.(*inheritableThreadLocal[personCloneable]).id)
	assert.NotSame(t, &value, &getValue)
	assert.Equal(t, value, getValue)

	mp2 := createInheritedMap()
	assert.NotNil(t, mp2)
	assert.NotSame(t, mp, mp2)
	getValue2 := mp2.get(tls.(*inheritableThreadLocal[personCloneable]).id)
	assert.NotSame(t, &value, &getValue2)
	assert.Equal(t, value, getValue2)
}

func TestCreateInheritedMapPointer(t *testing.T) {
	tls := NewInheritableThreadLocal[*person]()
	value := &person{Id: 1, Name: "Hello"}
	tls.Set(value)
	srcValue := tls.Get()
	assert.Same(t, value, srcValue)
	assert.Equal(t, *value, *srcValue)

	mp := createInheritedMap()
	assert.NotNil(t, mp)
	getValue := mp.get(tls.(*inheritableThreadLocal[*person]).id).(*person)
	assert.Same(t, value, getValue)
	assert.Equal(t, *value, *getValue)

	mp2 := createInheritedMap()
	assert.NotNil(t, mp2)
	assert.NotSame(t, mp, mp2)
	getValue2 := mp2.get(tls.(*inheritableThreadLocal[*person]).id).(*person)
	assert.Same(t, value, getValue2)
	assert.Equal(t, *value, *getValue2)
}

func TestCreateInheritedMapCloneable(t *testing.T) {
	tls := NewInheritableThreadLocal[*personCloneable]()
	value := &personCloneable{Id: 1, Name: "Hello"}
	tls.Set(value)
	srcValue := tls.Get()
	assert.Same(t, value, srcValue)
	assert.Equal(t, *value, *srcValue)

	mp := createInheritedMap()
	assert.NotNil(t, mp)
	getValue := mp.get(tls.(*inheritableThreadLocal[*personCloneable]).id).(*personCloneable)
	assert.NotSame(t, value, getValue)
	assert.Equal(t, *value, *getValue)

	mp2 := createInheritedMap()
	assert.NotNil(t, mp2)
	assert.NotSame(t, mp, mp2)
	getValue2 := mp2.get(tls.(*inheritableThreadLocal[*personCloneable]).id).(*personCloneable)
	assert.NotSame(t, value, getValue2)
	assert.Equal(t, *value, *getValue2)
}

func TestFill(t *testing.T) {
	a := make([]Any, 6)
	fill(a, 4, 5, unset)
	for i := 0; i < 6; i++ {
		if i == 4 {
			assert.True(t, a[i] == unset)
		} else {
			assert.Nil(t, a[i])
			assert.True(t, a[i] != unset)
		}
	}
}
