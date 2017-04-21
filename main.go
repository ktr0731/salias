package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
)

type commandIO struct {
	reader            io.Reader
	writer, errWriter io.Writer
}

func main() {
	exitCode, err := run(&commandIO{
		writer:    os.Stdout,
		errWriter: os.Stderr,
	}, os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "\x1b[31msalias: %s\x1b[0m\n", err)
	}
	os.Exit(exitCode)
}

func run(cmdIO *commandIO, args []string) (int, error) {
	if len(args) < 1 {
		return 1, errors.New("invalid arguments, please set least one command as argument")
	}

	// コマンド名だけ指定された場合
	if len(args) == 1 {
		cmd := exec.Command(args[0])
		cmd.Stdout = cmdIO.writer
		cmd.Stderr = cmdIO.errWriter
		if err := cmd.Run(); err != nil {
			// TODO: exit code 取得
			fmt.Fprintln(os.Stderr, err)
			return 1, nil
		}
		return 0, nil
	}

	command, subCommand, subCommandArgs := args[0], args[1], args[2:]

	dir, err := homedir.Dir()
	if err != nil {
		return 1, fmt.Errorf("cannot get home dir: %s", err)
	}
	path := filepath.Join(dir, ".config", "salias", "salias.toml")
	if envPath := os.Getenv("SALIAS_PATH"); envPath != "" {
		if envPathAbs, err := filepath.Abs(envPath); err != nil {
			return 1, errors.New("passed salias path is invalid")
		} else if envPath != "" {
			path = envPathAbs
		}
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return 1, fmt.Errorf("config path: %s not found", path)
	} else if err != nil {
		return 1, fmt.Errorf("file status error: %s", err)
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return 1, fmt.Errorf("cannot read salias.toml: %s", err)
	}

	var commands interface{}
	err = toml.Unmarshal(bytes, &commands)
	if err != nil {
		return 1, fmt.Errorf("cannot unmarshal toml: %s", err)
	}

	var ok bool
	if commands, ok = commands.(map[string]interface{})[command]; !ok {
		return 1, errors.New("no such command in commands managed by salias")
	}

	var aliases map[string]interface{}
	if aliases, ok = commands.(map[string]interface{}); !ok {
		return 1, errors.New("no such sub-command in sub-commands by salias")
	}

	for k, alias := range aliases {
		if k != subCommand {
			continue
		}

		// コマンドラインから渡された引数 + エイリアス先の引数
		subArgs := strings.TrimSpace(alias.(string))
		newArgs := make([]string, 0, 1+len(subCommandArgs)+len(subArgs))
		if splitted := strings.Split(subArgs, " "); len(splitted) != 1 {
			newArgs = append(splitted, newArgs...)
		} else {
			newArgs = append(newArgs, splitted[0])
		}

		for _, arg := range subCommandArgs {
			newArgs = append(newArgs, arg)
		}

		cmd := exec.Command(command, newArgs...)
		cmd.Stdout = cmdIO.writer
		cmd.Stderr = cmdIO.errWriter
		if err = cmd.Run(); err != nil {
			// コマンド自体のエラーは拾わない
			// TODO: exit code 取得
			return 1, nil
		}
		return 0, nil
	}

	// Normal command
	cmd := exec.Command(command, args[1:]...)
	cmd.Stdout = cmdIO.writer
	cmd.Stderr = cmdIO.errWriter
	if err = cmd.Run(); err != nil {
		// TODO: exit code 取得
		return 1, nil
	}
	return 0, nil
}
