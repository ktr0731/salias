package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("invalid arguments")
	}

	command := os.Args[1]
	subCommand := os.Args[2]

	bytes, _ := ioutil.ReadFile("./salias.toml")

	var commands interface{}
	err := toml.Unmarshal(bytes, &commands)
	if err != nil {
		log.Fatal(err)
	}

	for k, alias := range commands.(map[string]interface{})[command].(map[string]interface{}) {
		if k == subCommand {
			fmt.Printf("%s is %s\n", k, alias)
			break
		}
	}
}
