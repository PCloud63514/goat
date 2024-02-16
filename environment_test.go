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

func TestEnvironment_GetRequiredPropertyString(t *testing.T) {
	env := NewEnvironment()

	msg, err := env.GetRequiredPropertyString("HELLO_WORLD_MSG")
	if err != nil {
		t.Fatal("Test Failed")
	}
	fmt.Println(msg)
}

func TestEnvironment_SetProperty(t *testing.T) {
	env := NewEnvironment()
	env.SetProperty("TEST", "ABCDEFG")
	msg, err := env.GetRequiredPropertyString("TEST")
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(msg)
}
