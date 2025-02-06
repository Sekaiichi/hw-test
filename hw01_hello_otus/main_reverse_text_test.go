package main

import "testing"

func TestReverseText(t *testing.T) {
	text := "Hello, OTUS!"
	expected := "!SUTO ,olleH"

	result := reverseText(text)

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
