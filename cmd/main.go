package main

import (
	"context"
	"fmt"
	"github.com/PCloud63514/goat"
)

func main() {
	goat.New(
		goat.Configuration(FooConfiguration{}, BarConfiguration{}),
		goat.Provide(NewFoo, NewBar),
	).Run()
}

func NewFoo() *Foo {
	return &Foo{
		OnStop: func(ctx context.Context) error {
			fmt.Sprintf("Foo OnStop!")
			return nil
		},
	}
}

func NewBar(foo *Foo) *Bar {
	return &Bar{
		foo: foo,
		OnStart: func(ctx context.Context) error {
			fmt.Sprintf("Bar OnStart!")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Sprintf("Bar OnStop!")
			return nil
		},
	}
}

type Foo struct {
	OnStop func(ctx context.Context) error
}
type Bar struct {
	foo     *Foo
	OnStart func(ctx context.Context) error
	OnStop  func(ctx context.Context) error
}

type FooConfiguration struct {
	AppName string `properties:"app.ame"`
	Version string `properties:"app.version"`
}

type BarConfiguration struct {
	ServerHost string `properties:"server.host"`
	ServerPort int    `properties:"server.port"`
}
