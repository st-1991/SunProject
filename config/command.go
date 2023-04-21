package config

import (
	"fmt"
	"os"
	"strings"
)

type Options []string

func (o Options) Params() map[string]string {
	commands := make(map[string]string)
	for _, arg := range os.Args[1:] {
		shell := strings.Split(arg, "=")
		if len(shell) != 2 {
			continue
		}
		if !o.isOptions(shell[0]) {
			continue
		}
		commands[shell[0]] = shell[1]
	}
	return commands
}

func (o Options) isOptions(key string) bool {
	fmt.Printf("%+v", o)
	for _, item := range o {
		if key == item {
			return true
		}
	}
	return false
}