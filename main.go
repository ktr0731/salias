package main

import (
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
	if len(os.Args) < 3 {
		showError("invalid arguments")
		os.Exit(1)
	}

	command, subCommand, args := os.Args[1], os.Args[2], os.Args[3:]

	dir, err := homedir.Dir()
	if err != nil {
		showError("cannot get home dir: %s", err)
		os.Exit(1)
	}
	path := filepath.Join(dir, ".config", "salias", "salias.toml")
	if envPath := os.Getenv("SALIAS_PATH"); envPath != "" {
		if envPathAbs, err := filepath.Abs(envPath); err != nil {
			showError("passed salias path is invalid")
			os.Exit(1)
		} else if envPath != "" {
			path = envPathAbs
		}
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		showError("config path: %s not found", path)
		os.Exit(1)
	} else if err != nil {
		showError("file status error: %s", err)
		os.Exit(1)
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		showError("cannot read salias.toml: %s", err)
		os.Exit(1)
	}

	var commands interface{}
	err = toml.Unmarshal(bytes, &commands)
	if err != nil {
		showError("cannot unmarshal toml: %s", err)
		os.Exit(1)
	}

	var ok bool
	if commands, ok = commands.(map[string]interface{})[command]; !ok {
		showError("no such command in commands managed by salias")
		os.Exit(1)
	}

	var aliases map[string]interface{}
	if aliases, ok = commands.(map[string]interface{}); !ok {
		showError("no such sub-command in sub-commands by salias")
		os.Exit(1)
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
			os.Exit(1)
		}
		return
	}

	// Normal command
	cmd := exec.Command(command, os.Args[2:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		// TODO: exit code 取得
		os.Exit(1)
	}
}

func showError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("\x1b[31msalias: %s\x1b[0m\n", format), args...)
}
