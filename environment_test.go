package goat

import (
	"fmt"
	"testing"
)

func TestNewEnvironment(t *testing.T) {
	fmt.Println("HHHHHH")
	env := newEnvironment()

	if env == nil {
		t.Fatal("env must not be null.")
	}
}

func TestEnvironment_getRequiredPropertyString(t *testing.T) {
	env := newEnvironment()

	msg, err := env.getRequiredPropertyString("HELLO_WORLD_MSG")
	if err != nil {
		t.Fatal("Test Failed")
	}
	fmt.Println(msg)
}

func TestEnvironment_setProperty(t *testing.T) {
	env := newEnvironment()
	env.setProperty("TEST", "ABCDEFG")
	msg, err := env.getRequiredPropertyString("TEST")
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(msg)
}
