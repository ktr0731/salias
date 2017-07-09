package main

import (
	"fmt"
)

func initSalias() error {
	path, err := getPath()
	if err != nil {
		return err
	}

	cmds, err := getCmds(path)
	if err != nil {
		return err
	}

	var aliases string
	for key := range cmds {
		aliases += fmt.Sprintf("alias %s='salias %s'\n", key, key)
	}
	fmt.Print(aliases)
	return nil
}
