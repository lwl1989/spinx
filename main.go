package main

import (
	"github.com/lwl1989/spinx/core"
	"log"
	"github.com/lwl1989/spinx/cmd"
	"github.com/lwl1989/spinx/daemon"
	"github.com/lwl1989/spinx/http"
)

func main(){
	http.Do()
	return
	c := cmd.GetCommand()

	h := c.Get("help")
	path := c.Get("config")

	if h != "" {
		cmd.ShowHelp()
		return
	}

	if path == "" {
		cmd.ShowHelp()
		return
	}

	config := core.GetVhosts(path)
	if config == nil {
		log.Panic("config not load")
	}
 	isDaemon := c.Get("daemon")

	if isDaemon == "0" {
		daemon.Run(config)
	}else{
		daemon.RunDaemon(config)
	}
}




