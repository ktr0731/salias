package main

import (
	"os"
	"os/exec"
	"strings"
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
}

func TestMain_errors(t *testing.T) {
	resetEnv := setTestEnv("SALIAS_PATH", "./salias_errors_test.toml")
	defer resetEnv()

	tests := []struct {
		commands []string
		expect   string
	}{
		{
			[]string{},
			"invalid arguments",
		},
		{
			[]string{"makisekurisu"},
			"no such command",
		},
		{
			[]string{"git", "makisekurisu"},
			"no such sub-command",
		},
	}

	for _, test := range tests {
		b, err := exec.Command("go", append([]string{"run", "main.go"}, test.commands...)...).Output()
		if err != nil {
			t.Error(err)
		}

		if len(b) == 0 {
			t.Error(err)
		}

		if !strings.Contains(string(b), test.expect) {
			t.Errorf("expect: %s, actual: %s", test.expect, string(b))
		}
	}
}
