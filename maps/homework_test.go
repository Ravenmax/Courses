package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Element struct {
	Key   int
	Value int
}
type Node struct {
	Element
	Left  *Node
	Right *Node
}

type OrderedMap struct {
	Root *Node
	size int
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{}
}

func (m *OrderedMap) Insert(key, value int) {
	newNode := &Node{
		Element: Element{Key: key, Value: value},
	}

	if m.Root == nil {
		m.Root = newNode
		m.size++
		return
	}

	inserted := m.InsertNode(m.Root, newNode)
	if inserted {
		m.size++
	}
}
func (m *OrderedMap) InsertNode(root, newNode *Node) bool {
	if newNode.Key < root.Key {
		if root.Left == nil {
			root.Left = newNode
			return true
		}
		return m.InsertNode(root.Left, newNode)
	} else if newNode.Key > root.Key {
		if root.Right == nil {
			root.Right = newNode
			return true
		}
		return m.InsertNode(root.Right, newNode)
	} else {
		// Ключ уже существует - обновляем значение
		root.Value = newNode.Value
		return false
	}
}

func (m *OrderedMap) Erase(key int) {
	var removed bool
	m.Root, removed = m.removeNode(m.Root, key)
	if removed {
		m.size--
	}

}

func (m *OrderedMap) removeNode(root *Node, key int) (*Node, bool) {
	if root == nil {
		return nil, false
	}

	removed := false

	if key < root.Key {
		root.Left, removed = m.removeNode(root.Left, key)
		return root, removed
	} else if key > root.Key {
		root.Right, removed = m.removeNode(root.Right, key)
		return root, removed
	}

	// Нашли узел для удаления
	removed = true

	// Случай 1: Нет дочерних узлов
	if root.Left == nil && root.Right == nil {
		return nil, removed
	}

	// Случай 2: Только один дочерний узел
	if root.Left == nil {
		return root.Right, removed
	}
	if root.Right == nil {
		return root.Left, removed
	}

	// Случай 3: Два дочерних узла
	// Находим минимальный узел в правом поддереве
	minRight := m.findMin(root.Right)
	root.Key = minRight.Key
	root.Value = minRight.Value
	root.Right, _ = m.removeNode(root.Right, minRight.Key)

	return root, removed
}
func (m *OrderedMap) findMin(root *Node) *Node {
	current := root
	for current.Left != nil {
		current = current.Left
	}
	return current
}

func (m *OrderedMap) Contains(key int) bool {
	return m.findNode(m.Root, key) != nil
}

func (m *OrderedMap) findNode(root *Node, key int) *Node {
	if root == nil {
		return nil
	}

	if key < root.Key {
		return m.findNode(root.Left, key)
	} else if key > root.Key {
		return m.findNode(root.Right, key)
	}
	return root
}
func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	m.ForEachNode(m.Root, action) // need to implement
}

func (m *OrderedMap) ForEachNode(root *Node, fn func(key int, value int)) {
	if root != nil {
		m.ForEachNode(root.Left, fn)
		fn(root.Key, root.Value)
		m.ForEachNode(root.Right, fn)
	}
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
