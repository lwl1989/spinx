package cmd

import (
	"log"
	"os"
	"net/http"
	"github.com/takama/daemon"
	"os/signal"
	"syscall"
	"github.com/lwl1989/spinx/logger"
	"strings"
	"github.com/lwl1989/spinx/core"
)

var stdLog *log.Logger

// Service has embedded daemon
type Service struct {
	Daemon daemon.Daemon
	Vhosts core.Vhosts
	Server *http.Server
	start  chan bool
	Cmd    *Command
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {
	service.start = make(chan bool)

	nowLogger := &logger.Log{
		Log: &logger.FileLog{
			FilePath: "daemon.log",
		},
	}
	stdLog = log.New(nowLogger, "", log.Ldate|log.Ltime)
	usage := "Usage: command install | remove | start | stop | status"
	// if received any kind of command, do it
	stdLog.Println(os.Args)
	stdLog.Println("config：", service.Vhosts)
	if len(os.Args) > 1 {
		command := ""
		other := ""
		for _, c := range os.Args {
			if strings.Index(c, "-") == -1 {
				command = c
			} else {
				other += c + " "
			}
		}

		switch command {
		case "install":
			return service.Daemon.Install(other)
		case "remove":
			return service.Daemon.Remove()
		case "start":
			return service.Daemon.Start()
		case "stop":
			return service.Daemon.Stop()
		case "status":
			return service.Daemon.Status()
		default:
			return usage, nil
		}
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	go service.DoStart()

	for {
		select {
		case killSignal := <-interrupt:
			stdLog.Println("Got signal:", killSignal)
			if service.Server != nil {
				service.Server.Close()
			}

			if killSignal == os.Interrupt {
				return "Server exit", nil
			}
			return "Daemon was killed", nil
		}
	}
	stdLog.Println(usage)
	// never happen, but need to complete code
	return usage, nil
}

func (service *Service) DoStart() {
	stdLog.Println("doStart")
	ports := service.Vhosts.GetPorts()

	namePortsListen := make([]string,0)
	handlerMap := make(map[string]*http.ServeMux)
	for _,port := range ports {
		names := service.Vhosts.GetNames(port)
		for _,name := range names{
			hMap,err := service.Vhosts.GetHostMap(port,name)
			if err != nil {
				log.Fatal("port:"+port+" name:"+name+" config is not found")
				continue
			}
			staticHandler := http.NewServeMux()
			staticHandler.Handle("/", http.FileServer(http.Dir(hMap.DocumentRoot)))
			handlerMap[name+port] = staticHandler
			//http.HandleFunc("/", handler)
			listen := "localhost"
			if port != "80" {
				listen = listen+":"+port
			}
			namePortsListen = append(namePortsListen, listen)
		}
	}
	handler := core.GetHandler()
	handler.Vhosts = service.Vhosts
	handler.HandlerMap = handlerMap
	for _,listen := range namePortsListen {
		err := http.ListenAndServe(listen, handler)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
