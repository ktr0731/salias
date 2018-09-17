package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	xdgbasedir "github.com/zchee/go-xdgbasedir"
)

type commandIO struct {
	reader            io.Reader
	writer, errWriter io.Writer
}

func showError(err error) {
	fmt.Fprintf(os.Stderr, "\x1b[31msalias: %s\x1b[0m\n", err)
}

func execCmd(cmdIO *commandIO, cmdName string, args ...string) int {
	path, err := exec.LookPath(cmdName)
	if err != nil {
		log.Println(err)
		return 1
	}

	if err := syscall.Exec(path, append([]string{cmdName}, args...), os.Environ()); err != nil {
		return 1
	}
	return 0
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getPath() (string, error) {
	var path string
	paths := []string{"salias.toml", ".salias.toml"}

	xdgConfigHome := xdgbasedir.ConfigHome()

	// first, check xdg dir
	for _, name := range paths {
		path = filepath.Join(xdgConfigHome, "salias", name)
		if isExist(path) {
			return path, nil
		}
	}
	dir, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("cannot get home dir: %s", err)
	}

	if envPath := os.Getenv("SALIAS_PATH"); envPath != "" {
		if envPathAbs, err := filepath.Abs(envPath); err != nil {
			return "", errors.New("passed salias path is invalid")
		} else if envPath != "" {
			path = envPathAbs
		}
		if isExist(path) {
			return path, nil
		}
		return "", errors.New("path specified by $SALIAS_PATH is not exists")
	}

	// if not found, check home dir
	for _, name := range paths {
		path = filepath.Join(dir, name)
		if isExist(path) {
			return path, nil
		}
	}

	return "", errors.New("config file salias.toml not found")
}

// map[go:map[i:install r:run] docker:map[i:image]]
type commands map[string]command

// map[i:install r:run]
type command map[string]string

func getCmds() (commands, error) {
	path, err := getPath()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find salias path")
	}

	var cmds commands
	_, err = toml.DecodeFile(path, &cmds)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read salias.toml")
	}

	return cmds, nil
}

func writeCmds(cmds commands) error {
	buf := new(bytes.Buffer)
	enc := toml.NewEncoder(buf)
	if err := enc.Encode(&cmds); err != nil {
		return errors.Wrap(err, "failed to encode")
	}

	path, perr := getPath()
	if perr != nil {
		return errors.Wrap(perr, "failed to find salias path")
	}

	file, ferr := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModeAppend)
	if ferr != nil {
		return errors.Wrap(ferr, "cannot open salias.toml")
	}
	defer file.Close()
	if _, err := file.WriteString(buf.String()); err != nil {
		return errors.Wrap(err, "write string error")
	}
	return nil
}

func run(cmdIO *commandIO, args []string) (int, error) {
	if len(args) < 1 {
		return 1, errors.New("invalid arguments, please set least one command as argument")
	}

	// just like: salias -r <command>
	if len(args) == 1 {
		return execCmd(cmdIO, args[0]), nil
	}

	cmd, subCmd, subCmdArgs := args[0], args[1], args[2:]

	cmds, err := getCmds()
	if err != nil {
		return 1, errors.Wrap(err, "failed to get commands from config file")
	}

	// if an executable "cmd", but not in salias config file
	aliases := cmds[cmd]
	if aliases == nil {
		return 1, errors.New("no such command in commands managed by salias")
	}

	for subCmdName, alias := range aliases {
		if subCmdName != subCmd {
			continue
		}

		// has "!" prefix for another command
		if strings.HasPrefix(alias, "!") {
			alias = alias[1:]
			subArgs := strings.Split(strings.TrimSpace(alias), " ")
			if len(subArgs) == 1 {
				return execCmd(cmdIO, subArgs[0]), nil
			}
			return execCmd(cmdIO, subArgs[0], subArgs[1:]...), nil
		}

		// args passed by alias + args passed by command-line
		subArgs := strings.Split(strings.TrimSpace(alias), " ")
		newArgs := make([]string, 0, 1+len(subCmdArgs)+len(subArgs))

		newArgs = append(subArgs, append(newArgs, subCmdArgs...)...)
		return execCmd(cmdIO, cmd, newArgs...), nil
	}

	// Normal command
	return execCmd(cmdIO, cmd, args[1:]...), nil
}

func initSalias() (int, error) {
	cmds, err := getCmds()
	if err != nil {
		return 1, errors.Wrap(err, "failed to generate init script")
	}

	var aliases string
	for key := range cmds {
		aliases += fmt.Sprintf("alias %s='salias --run %s'\n", key, key)
	}
	fmt.Print(aliases)
	return 0, nil
}

func setAlias(args []string) (int, error) {
	// just like: salias go i="install"
	cmds, cerr := getCmds()
	if cerr != nil {
		return 1, errors.Wrap(cerr, "cannot read salias.toml")
	}

	var cmd command
	// make section of a command if not exist
	if c := cmds[args[0]]; c != nil {
		cmd = c
	} else {
		cmd = make(command)
	}
	// alias[0]: name, alias[1]: value
	alias := strings.Split(args[1], "=")
	if len(alias) == 1 {
		if value := cmds[args[0]][alias[0]]; value != "" {
			fmt.Println(value)
			return 0, nil
		}
		return 1, nil
	}
	cmd[alias[0]] = alias[1]
	cmds[args[0]] = cmd

	if err := writeCmds(cmds); err != nil {
		return 1, errors.Wrap(err, "cannot write salias.toml")
	}
	return 0, nil
}

func controller(args []string) (int, error) {
	if len(args) == 1 {
		// verify and show defined sub alias
		cmds, err := getCmds()
		if err != nil {
			return 1, errors.Wrap(err, "verify salias.toml failed\n")
		}
		enc := toml.NewEncoder(os.Stdout)
		enc.Encode(&cmds)
		return 0, nil
	}
	switch args[1] {
	case "--init", "-i", "__init__":
		return initSalias()
	case "--run", "-r":
		return run(&commandIO{
			reader:    os.Stdin,
			writer:    os.Stdout,
			errWriter: os.Stderr,
		}, args[2:])
	default:
		return setAlias(args[1:])
	}
}

func main() {
	exitCode, err := controller(os.Args)
	if err != nil {
		showError(err)
	}
	os.Exit(exitCode)
}
