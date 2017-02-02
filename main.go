package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("invalid arguments")
	}

	command := os.Args[1]
	subCommand := os.Args[2]
	args := os.Args[3:]

	bytes, _ := ioutil.ReadFile("./salias.toml")

	var commands interface{}
	err := toml.Unmarshal(bytes, &commands)
	if err != nil {
		log.Fatal(err)
	}

	if commands, ok := commands.(map[string]interface{})[command]; ok {
		if aliases, ok := commands.(map[string]interface{}); ok {
			for k, alias := range aliases {
				if k == subCommand {
					newArgs := make([]string, 1+len(args))
					newArgs[0] = alias.(string)
					for i := range args {
						newArgs[i+1] = args[i]
					}
					fmt.Println(command, strings.Join(newArgs, " "))
					b, err := exec.Command(command, newArgs...).Output()
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println(string(b))
				}
			}
		}
	}

}
