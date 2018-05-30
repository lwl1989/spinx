package daemon

import (
	"net/http"
	"os"
	"log"
	"github.com/lwl1989/spinx/core"
)

func Run(config core.Vhosts) {

	daemonCommandArrays := []string{"start", "stop", "install", "remove"}
	for _, value := range os.Args {
		for _, command := range daemonCommandArrays {
			if command == value {
				RunDaemon(config)
				return
			}
		}
	}
	handler := core.GetHandler()
	handler.Vhosts = config

	ports := config.GetPorts()
	namePortsListen := make([]string,0)

	handlerMap := make(map[string]*http.ServeMux)
	for _,port := range ports {
		names := config.GetNames(port)
		//fmt.Println(names)
		for _,name := range names{
			hMap,err := config.GetHostMap(port,name)
			if err != nil {
				log.Fatal("port:"+port+" name:"+name+" config is not found")
				continue
			}
			staticHandler := http.NewServeMux()
			staticHandler.Handle("/", http.FileServer(http.Dir(hMap.DocumentRoot)))
			handlerMap[name+port] = staticHandler
			listen := "localhost"
			if port != "80" {
				listen = listen+":"+port
			}
			namePortsListen = append(namePortsListen, listen)
		}
	}


	handler.HandlerMap = handlerMap
	for _,listen := range namePortsListen {
		err := http.ListenAndServe(listen, handler)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
