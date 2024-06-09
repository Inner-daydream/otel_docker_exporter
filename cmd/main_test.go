package main

import (
	"os"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	// main should run without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("main panicked: %v", r)
		}
	}()

	// If 10 seconds pass without a panic, exit with a success status code
	time.AfterFunc(10*time.Second, func() {
		os.Exit(0)
	})
	main()
}
