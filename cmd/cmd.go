package cmd

import (
	"flag"
)

type Command struct {
	list  map[string]string
	alias map[string]string
}

func GetCommand() *Command {
	c := &Command{
		make(map[string]string),
		make(map[string]string),
	}
	c.init()
	return c
}

func (command *Command) init() {
	path := ""
	help := ""
	daemon := false
	flag.StringVar(&help, "h", "", "")
	flag.StringVar(&path, "c", "", "")
	flag.BoolVar(&daemon, "d", false, "")
	flag.Parse()

	if path == "" {
		path = getCurrentDirectory() + "/conf/config.json"
	}

	if daemon {
		command.SetAlias("d", "daemon", "1")
	} else {
		command.SetAlias("d", "daemon", "0")
	}

	command.SetAlias("h", "help", help)
	command.SetAlias("c", "config", path)
}

func (command *Command) Get(key string) string {
	if _, ok := command.list[key]; ok {
		return command.list[key]
	} else {
		if _, ok := command.alias[key]; ok {
			key = command.alias[key]
		}
		return command.Get(key)
	}

	return ""
}

func (command *Command) Set(key string, value string) {
	command.list[key] = value
}

func (command *Command) SetAlias(key string, alias string, value string) {
	command.Set(key, value)
	command.alias[alias] = key
}
