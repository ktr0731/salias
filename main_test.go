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
	resetEnv := setTestEnv("SALIAS_PATH", "./salias_test.toml")
	defer resetEnv()

	tests := []struct {
		commands []string
		expect   string
	}{
		{
			[]string{},
			"invalid arguments, please set least one command as argument",
		},
		{
			[]string{"g", "makisekurisu"},
			"no such command in commands managed by salias",
		},
	}

	for _, test := range tests {
		outBuf, errBuf := new(bytes.Buffer), new(bytes.Buffer)
		exitCode, err := run(&commandIO{
			writer:    outBuf,
			errWriter: errBuf,
		}, test.commands)

		if err == nil && errBuf.Len() == 0 {
			t.Error("error not occurred")
		}

		if exitCode == 0 {
			t.Error("exit with 0")
		}

		if errBuf.Len() != 0 && !strings.Contains(errBuf.String(), test.expect) {
			t.Errorf("errBuf: expect: %s, actual: %s", test.expect, errBuf)
		}

		if err != nil && !strings.Contains(err.Error(), test.expect) {
			t.Errorf("err: expect: %s, actual: %s", test.expect, err.Error())
		}
	}
}

func do_getPath(path string) error {
	resetEnv := setTestEnv("SALIAS_PATH", path)
	defer resetEnv()

	_, err := getPath()

	return err
}

func Test_getPath_errors(t *testing.T) {
	resetEnv := setTestEnv("SALIAS_PATH", "./salias_test.toml")
	defer resetEnv()

	tests := []struct {
		path   string
		expect string
	}{
		{
			path:   "./hoge.toml",
			expect: "path specified by $SALIAS_PATH is not exists",
		},
		// {
		// 	path:   "",
		// 	expect: "config file salias.toml not found",
		// },
	}

	for _, test := range tests {
		err := do_getPath(test.path)
		if err == nil {
			t.Error("error not occurred")
		}
		if !strings.Contains(err.Error(), test.expect) {
			t.Errorf("expect: %s, actual: %s", test.expect, err.Error())
		}
	}
}
