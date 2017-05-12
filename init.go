package main

import (
	"fmt"
	"os"
	"strings"
)

func initSalias() {
	path, err := getPath()
	if err != nil {
		showError(err)
		os.Exit(1)
	}

	cmds, err := getCmds(path)
	if err != nil {
		showError(err)
		os.Exit(1)
	}

	command := "alias"
	if strings.Contains(os.Getenv("SHELL"), "fish") {
		command = "abbr"
	}
	var aliases string
	for key := range cmds {
		aliases += fmt.Sprintf("%s %s='salias %s'\n", command, key, key)
	}
	fmt.Print(aliases)
}
