package main

import "testing"

func TestBuild(t *testing.T) {
	if true == false {
		t.Errorf("expected true != false")
	}
}
