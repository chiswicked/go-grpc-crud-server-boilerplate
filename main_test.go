package main

import (
	"testing"
)

func TestStartMsg(t *testing.T) {
	if startMsg("test") != "Initializing test server" {
		t.Fail()
	}
}
