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

func showError(err error) {
	fmt.Fprintf(os.Stderr, "\x1b[31msalias: %s\x1b[0m\n", err)
}

func execCmd(cmdIO *commandIO, cmdName string, args ...string) int {
	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = cmdIO.writer
	cmd.Stderr = cmdIO.errWriter
	if err := cmd.Run(); err != nil {
		// TODO: exit code 取得
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getPath() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		panic(fmt.Sprintf("cannot get home dir: %s", err))
	}
	// path := filepath.Join(dir, ".config", "salias", "salias.toml")
	var path string
	if envPath := os.Getenv("SALIAS_PATH"); envPath != "" {
		if envPathAbs, err := filepath.Abs(envPath); err != nil {
			panic("passed salias path is invalid")
		} else if envPath != "" {
			path = envPathAbs
		}
		if isExist(path) {
			return path, nil
		}
		return "", fmt.Errorf("Path specified by $SALIAS_PATH is not exists")
	}

	path = filepath.Join(dir, ".config", "salias", "salias.toml")
	if isExist(path) {
		return path, nil
	}

	path = filepath.Join(dir, "salias.toml")
	if isExist(path) {
		return path, nil
	}

	return "", errors.New("config file salias.toml not found")
}

func getCmds(path string) (interface{}, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read salias.toml: %s", err)
	}

	var cmds interface{}
	err = toml.Unmarshal(b, &cmds)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal toml: %s", err)
	}

	return cmds, nil
}

func run(cmdIO *commandIO, args []string) (int, error) {
	if len(args) < 1 {
		return 1, errors.New("invalid arguments, please set least one command as argument")
	}

	// Init
	if args[0] == "__init__" {
		initSalias()
		return 0, nil
	}

	// コマンド名だけ指定された場合
	if len(args) == 1 {
		return execCmd(cmdIO, args[0]), nil
	}

	cmd, subCmd, subCmdArgs := args[0], args[1], args[2:]

	path, err := getPath()
	if err != nil {
		return 1, err
	}

	cmds, err := getCmds(path)
	if err != nil {
		return 1, err
	}

	var ok bool
	if cmds, ok = cmds.(map[string]interface{})[cmd]; !ok {
		return 1, errors.New("no such command in commands managed by salias")
	}

	var aliases map[string]interface{}
	if aliases, ok = cmds.(map[string]interface{}); !ok {
		return 1, errors.New("no such sub-command in sub-commands by salias")
	}

	for k, alias := range aliases {
		if k != subCmd {
			continue
		}

		// コマンドラインから渡された引数 + エイリアス先の引数
		subArgs := strings.TrimSpace(alias.(string))
		newArgs := make([]string, 0, 1+len(subCmdArgs)+len(subArgs))
		if splitted := strings.Split(subArgs, " "); len(splitted) != 1 {
			newArgs = append(splitted, newArgs...)
		} else {
			newArgs = append(newArgs, splitted[0])
		}

		for _, arg := range subCmdArgs {
			newArgs = append(newArgs, arg)
		}

		return execCmd(cmdIO, cmd, newArgs...), nil
	}

	// Normal command
	return execCmd(cmdIO, cmd, args[1:]...), nil
}

func main() {
	exitCode, err := run(&commandIO{
		writer:    os.Stdout,
		errWriter: os.Stderr,
	}, os.Args[1:])
	if err != nil {
		showError(err)
	}
	os.Exit(exitCode)
}
