package main

import (
	"context"
	"github.com/PCloud63514/goat"
)

func main() {
	goat.New(
		goat.Provide(NewFoo, NewBar),
	).Run()
}

func NewFoo() *Foo {
	return &Foo{}
}

func NewBar(foo *Foo) *Bar {
	return &Bar{foo: foo}
}

type Foo struct {
	OnStop func(ctx context.Context) error
}
type Bar struct {
	foo    *Foo
	OnStop func(ctx context.Context) error
}
