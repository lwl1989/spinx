package main

import (
	"github.com/spinx/core"
	"net/http"
	"fmt"
//	"log"
	"log"
)

var host core.Vhosts
var ports []string

var handlerMap map[string]*http.ServeMux
func main(){
	path,real := "","/usr/etc/spinx/server.json"
	fmt.Println("Please input your config path(default /usr/etc/spinx/server.json): ")
	fmt.Scanln(&path)
	if path != "" {
		real = path
	}
	real = "/www/spinx/src/config/server.json"

	host = core.GetVhosts(real)
	ports = host.GetPorts()
	handler := core.GetHandler()
	handlerMap = make(map[string]*http.ServeMux)
	handler.Vhosts = host

	handler.HandlerMap = handlerMap

	namePortsListen := make([]string,0)
	for _,port := range ports {
		names := host.GetNames(port)
		//fmt.Println(names)
		for _,name := range names{
			hMap,err := host.GetHostMap(port,name)
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


	for _,listen := range namePortsListen {
		err := http.ListenAndServe(listen, handler)
		fmt.Println(err)
	}

	fmt.Printf("Press Ctrl-C to quit.\n")
}




