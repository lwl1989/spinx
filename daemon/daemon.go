package daemon

import (
	"github.com/lwl1989/spinx/cmd"
	"github.com/takama/daemon"
	"log"
	"github.com/lwl1989/spinx/core"
)

func RunDaemon(config core.Vhosts) {

	srv, err := daemon.New("spinx", "spinx fastcgi proxy service")
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	service := &cmd.Service{Daemon: srv, Vhosts: config}
	status, err := service.Manage()
	if err != nil {
		log.Fatalln(status, "\nError: ", err)
	}
	log.Println(status)
}
