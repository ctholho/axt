package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("TZ") != "UTC" {
		fmt.Println("Error: Tests must be run with the 'TZ=UTC' environment variable.")
		fmt.Println("E.g. run tests with 'TZ=UTC go test ./...'.")
		os.Exit(1)
	}

	os.Exit(m.Run())
}
