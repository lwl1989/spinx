package main

import (
	"github.com/lwl1989/spinx/core"
	"log"
	"github.com/lwl1989/spinx/cmd"
	"github.com/lwl1989/spinx/daemon"
	"github.com/lwl1989/spinx/http"
	"github.com/lwl1989/spinx/conf"
)

func main(){

	//str := "\n"
	//b := bytes.NewBufferString(str).Bytes()
	//var b1 byte = 10
	//
	//for _,v := range b {
	//	if v == b1 {
	//		fmt.Println("equals")
	//	}
	//}
	//fmt.Println(b)
	//fmt.Println(b1)
	//return
	//str := "www.test.com:18000"
	//fmt.Println(strings.Index(str,":"))
	conf.GetVhosts("config/server.json")
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




