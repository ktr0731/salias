package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "invalid arguments")
		os.Exit(1)
	}

	command := os.Args[1]
	subCommand := os.Args[2]
	args := os.Args[3:]

	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(dir, ".config", "salias", "salias.toml")
	if envPath := os.Getenv("SALIAS_PATH"); envPath != "" {
		if envPathAbs, err := filepath.Abs(envPath); err != nil {
			log.Fatal(err)
		} else if envPath != "" {
			path = envPathAbs
		}
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		log.Fatalf("config path: %s not found", path)
	} else if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var commands interface{}
	err = toml.Unmarshal(bytes, &commands)
	if err != nil {
		log.Fatal(err)
	}

	var ok bool
	if commands, ok = commands.(map[string]interface{})[command]; !ok {
		log.Fatal("type assertion failed")
	}

	var aliases map[string]interface{}
	if aliases, ok = commands.(map[string]interface{}); !ok {
		log.Fatal("type assertion failed2")
	}

	for k, alias := range aliases {
		if k == subCommand {
			newArgs := make([]string, 1+len(args))
			newArgs[0] = alias.(string)
			for i := range args {
				newArgs[i+1] = args[i]
			}
			cmd := exec.Command(command, newArgs...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err = cmd.Run(); err != nil {
				os.Exit(1)
			}
			os.Exit(0)
		}
	}
}
