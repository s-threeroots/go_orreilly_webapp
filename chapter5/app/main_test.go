package main

import (
	"testing"
)

func Test_dialdb(t *testing.T) {

	err := dialdb()
	if err != nil {
		t.Error(err)
	}
}
