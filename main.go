package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
)

func main() {
	exitCode, err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\x1b[31msalias: %s\x1b[0m\n", err)
	}
	os.Exit(exitCode)
}

func run() (int, error) {
	if len(os.Args) < 3 {
		return 1, errors.New("invalid arguments")
	}

	command, subCommand, args := os.Args[1], os.Args[2], os.Args[3:]

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
		newArgs := make([]string, 0, 1+len(args)+len(subArgs))
		if splitted := strings.Split(subArgs, " "); len(splitted) != 1 {
			newArgs = append(splitted, newArgs...)
		} else {
			newArgs[0] = splitted[0]
		}

		for i := range args {
			newArgs[i+1] = args[i]
		}

		cmd := exec.Command(command, newArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			// コマンド自体のエラーは拾わない
			// TODO: exit code 取得
			return 1, nil
		}
		return 0, nil
	}

	// Normal command
	cmd := exec.Command(command, os.Args[2:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		// TODO: exit code 取得
		return 1, nil
	}
	return 0, nil
}
