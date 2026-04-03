package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	List []Task
}

func HeapSize(s *Scheduler) int {
	return len(s.List)

}
func (s *Scheduler) HeapUp(i int) {
	parent := (i - 1) / 2

	for i > 0 && s.List[parent].Priority < s.List[i].Priority {
		// Обмен значений
		s.List[i], s.List[parent] = s.List[parent], s.List[i]

		// Перемещаемся вверх по дереву
		i = parent
		parent = (i - 1) / 2
	}
}
func (s *Scheduler) HeapDown(i int) {
	for {
		left := 2*i + 1
		right := 2*i + 2
		largest := i

		if left < HeapSize(s) && s.List[left].Priority > s.List[largest].Priority {
			largest = left
		}
		if right < HeapSize(s) && s.List[right].Priority > s.List[largest].Priority {
			largest = right
		}

		if largest == i {
			break
		}

		s.List[i], s.List[largest] = s.List[largest], s.List[i]
		i = largest
	}
}

func NewScheduler() Scheduler {
	return Scheduler{List: make([]Task, 0)}
}

func (s *Scheduler) AddTask(task Task) {
	s.List = append(s.List, task)
	i := HeapSize(s) - 1
	parent := (i - 1) / 2
	for i > 0 && s.List[parent].Priority < s.List[i].Priority {
		s.List[i], s.List[parent] = s.List[parent], s.List[i]
		i = parent
		parent = (i - 1) / 2
	}
	s.HeapUp(HeapSize(s) - 1)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	if taskID < 0 || taskID >= HeapSize(s) {
		panic(fmt.Errorf("индекс %d вне диапазона [0, %d]", taskID, HeapSize(s)))
	}

	oldPriority := s.List[taskID].Priority
	s.List[taskID].Priority = newPriority

	// Если значение увеличилось, нужно поднять элемент вверх
	if newPriority > oldPriority {
		s.HeapUp(taskID)
	} else if newPriority < oldPriority {
		// Если значение уменьшилось, нужно опустить элемент вниз
		s.HeapDown(taskID)
	}
}

// ChangePriorityByIdentifier изменяет приоритет задачи по идентификатору
func (s *Scheduler) ChangePriorityByIdentifier(identifier, newPriority int) bool {
	index := s.findByIdentifier(identifier)
	if index == -1 {
		return false
	}

	s.ChangeTaskPriority(index, newPriority)
	return true
}

// findIdentifier ищет индекс задачи по идентификатору
func (s *Scheduler) findByIdentifier(identifier int) int {
	for i, task := range s.List {
		if task.Identifier == identifier {
			return i
		}
	}
	return -1
}

func (s *Scheduler) GetTask() Task {
	if len(s.List) == 0 {
		panic("heap is empty")
	}

	result := s.List[0]
	lastIndex := HeapSize(s) - 1
	s.List[0] = s.List[lastIndex]
	s.List = s.List[:lastIndex]
	if HeapSize(s) > 0 {
		s.HeapDown(0)
	}
	return result
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangePriorityByIdentifier(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
