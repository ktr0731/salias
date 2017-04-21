package main

import (
	"fmt"
	"os"
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

	var aliases string
	for key := range cmds {
		aliases += fmt.Sprintf("alias %s='salias %s'\n", key, key)
	}
	fmt.Print(aliases)
}
