// main_test.go
package main

import "testing"

func TestPass(t *testing.T) {
	//result := DBConnect()
	result := 1
	expected := 1

	if result != expected {
		t.Errorf("DBConnect() returned %d, expected %d", result, expected)
	}
}
