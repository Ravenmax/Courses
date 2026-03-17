package main

import (
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type COWBuffer struct {
	data     []byte
	refCount *int32
	mutex    sync.Mutex
}

func NewCOWBuffer(data []byte) COWBuffer {
	refCount := int32(1)
	return COWBuffer{
		data:     data,
		refCount: &refCount,
	}
}

func (b *COWBuffer) Clone() COWBuffer {
	atomic.AddInt32(b.refCount, 1)
	return COWBuffer{
		data:     b.data,
		refCount: b.refCount,
	}
}

func (b *COWBuffer) Close() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	newCount := atomic.AddInt32(b.refCount, -1)
	if newCount == 0 {
		b.data = nil
	}
}

func (b *COWBuffer) Update(index int, value byte) bool {
	//залочим объект во избежания ошибок двойного доступа
	b.mutex.Lock()
	//в конце разлочим в любом случае
	defer b.mutex.Unlock()
	//проверка на границы индекса, если выходим за границы сразу ошибка
	if index < 0 || index >= len(b.data) {
		return false
	} else {
		//смотрим количество ссылок на этот объект
		currentRefs := atomic.LoadInt32(b.refCount)
		//если кто-то уже ссылается делаем копию объекта и там меняем данные
		if currentRefs > 1 {
			// копируем данные и делаем новый указатель на количество ссылок
			newDataCopy := make([]byte, len(b.data))
			copy(newDataCopy, b.data)
			newDataCopy[index] = value
			atomic.AddInt32(b.refCount, -1)
			newRefCount := int32(1)
			b.refCount = &newRefCount
			b.data = newDataCopy
		} else {
			//просто меняем значение по индексу
			b.data[index] = value
		}
	}
	return true
}

func (b *COWBuffer) String() string {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if len(b.data) == 0 {
		return ""
	}
	//возвращаем строку с указателем на её массив взятый от указателя на массив байт
	str := unsafe.String(unsafe.SliceData(b.data), len(b.data))
	return str
}

func TestCOWBuffer(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()

	copy1 := buffer.Clone()
	copy2 := buffer.Clone()

	assert.Equal(t, unsafe.SliceData(data), unsafe.SliceData(buffer.data))
	assert.Equal(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	assert.True(t, (*byte)(unsafe.SliceData(data)) == unsafe.StringData(buffer.String()))
	assert.True(t, (*byte)(unsafe.StringData(buffer.String())) == unsafe.StringData(copy1.String()))
	assert.True(t, (*byte)(unsafe.StringData(copy1.String())) == unsafe.StringData(copy2.String()))

	assert.True(t, buffer.Update(0, 'g'))
	assert.False(t, buffer.Update(-1, 'g'))
	assert.False(t, buffer.Update(4, 'g'))

	assert.True(t, reflect.DeepEqual([]byte{'g', 'b', 'c', 'd'}, buffer.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy1.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy2.data))

	assert.NotEqual(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	copy1.Close()

	previous := copy2.data
	copy2.Update(0, 'f')
	current := copy2.data

	// 1 reference - don't need to copy buffer during update
	assert.Equal(t, unsafe.SliceData(previous), unsafe.SliceData(current))

	copy2.Close()
}
