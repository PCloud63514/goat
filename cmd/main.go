package main

import (
	"fmt"
)

type Foo struct {
}

func (foo *Foo) Hello() {
	fmt.Println("Hello!")
}

func (foo *Foo) Handle() {
	fmt.Println("Foo Handle!")
	foo.Hello()
}

func (foo *Foo) Start() {
	fmt.Println("Foo Start!")
	foo.Hello()
}

type Bar struct {
}

func (bar *Bar) Start() {
	fmt.Println("Bar Start!")
}

func (bar *Bar) Stop() {
	fmt.Println("Bar Stop!")
}

func main() {
	var instances []any = []any{&Foo{}, &Bar{}}

	handlers := GetInstanceByType[Handler](instances)
	for _, handler := range handlers {
		handler.Handle()
	}

	startEventListeners := GetInstanceByType[StartEventListener](instances)
	for _, eventListener := range startEventListeners {
		eventListener.Start()
	}

	stopEventListeners := GetInstanceByType[StopEventListener](instances)
	for _, eventListener := range stopEventListeners {
		eventListener.Stop()
	}

	//GoatApplication().Run()
}

func GetInstanceByType[T any](instances []any) []T {
	result := []T{}

	for _, instance := range instances {
		if v, ok := instance.(T); ok {
			result = append(result, v)
		}
	}
	return result
}

type Handler interface {
	Handle()
}

type StartEventListener interface {
	Start()
}

type StopEventListener interface {
	Stop()
}
