package utils

import (
	"fmt"
	"testing"
)

func TestGetTime(t *testing.T) {

}

func TestRunPanicless(t *testing.T) {
	RunPanicless(func() {
		panic(1)
	})
	RunPanicless(func() {
		panic(fmt.Errorf("bad"))
	})
}
