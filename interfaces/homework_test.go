package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}

type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	services map[string]*ServiceDefinition
}

type ServiceDefinition struct {
	Name        string
	Constructor interface{}
	Instance    interface{}
}

func NewContainer() *Container {
	return &Container{
		services: make(map[string]*ServiceDefinition),
	}
}

func (c *Container) RegisterType(name string, constructor interface{}) {

	c.services[name] = &ServiceDefinition{
		Name:        name,
		Constructor: constructor,
	}
}

func (c *Container) Resolve(name string) (interface{}, error) {
	// Находим определение сервиса
	def, exists := c.services[name]
	if !exists {
		return nil, fmt.Errorf("service '%s' not found", name)
	}
	constructorValue := reflect.ValueOf(def.Constructor)

	// Вызываем конструктор
	results := constructorValue.Call(nil)
	instance := results[0].Interface()
	return instance, nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}
