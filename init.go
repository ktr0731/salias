package main

import (
	"fmt"
)

func initSalias() error {
	cmds, err := getCmds()
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
