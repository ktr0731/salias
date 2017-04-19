package main

import (
	"bytes"
	"os"
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

	outBuf, errBuf := new(bytes.Buffer), new(bytes.Buffer)
	exitCode, err := run(&commandIO{
		writer:    outBuf,
		errWriter: errBuf,
	}, []string{"git", "l"})
	if err != nil {
		t.Errorf("error: %s, errBuf: %s", err, errBuf)
	}
	if exitCode != 0 {
		t.Errorf("exit with %d", exitCode)
	}

	if outBuf.Len() == 0 {
		t.Error("output is empty")
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
		outBuf, errBuf := new(bytes.Buffer), new(bytes.Buffer)
		exitCode, err := run(&commandIO{
			writer:    outBuf,
			errWriter: errBuf,
		}, test.commands)
		if err != nil {
			t.Errorf("error: %s, errBuf: %s", err, errBuf)
		}
		if exitCode != 0 {
			t.Errorf("exit with %d", exitCode)
		}

		if outBuf.Len() == 0 {
			t.Error("output is empty")
		}

		if !strings.Contains(outBuf.String(), test.expect) {
			t.Errorf("expect: %s, actual: %s", test.expect, outBuf)
		}
	}
}
