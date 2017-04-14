package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func setTestEnv(key, val string) func() {
	preVal := os.Getenv(key)
	os.Setenv(key, val)

	return func() {
		os.Setenv(key, preVal)
	}
}

func TestMain(t *testing.T) {
	resetEnv := setTestEnv("SALIAS_PATH", "./salias_test.toml")
	defer resetEnv()

	b, err := exec.Command("go", "run", "main.go", "git", "l").Output()
	if err != nil {
		t.Error(err)
	}

	if len(b) == 0 {
		t.Error(err)
	}

	fmt.Println(string(b))
}
