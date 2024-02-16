package goat

import (
	"fmt"
	"testing"
)

func TestNewEnvironment(t *testing.T) {
	fmt.Println("HHHHHH")
	env := NewEnvironment()

	if env == nil {
		t.Fatal("env must not be null.")
	}
}
