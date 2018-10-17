package http

import (
	"net"
	"time"
	"log"
	"github.com/lwl1989/spinx/http/fcgi"
	"github.com/lwl1989/spinx/conf"
)


func Do()  {
	multiDo()
	normalDo()
}

func normalDo() {
	listen("8888")
}

func multiDo()  {
	ports := conf.HostMaps.GetPorts()

	for _,port := range ports {
		go listen(port)
	}
}

func listen(port string)  {
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("error listen:", err)
		return
	}
	defer l.Close()
	//log.Println("listen ok")



	for {
		Conn, err := l.Accept()
		//time.Sleep(time.Second * 10)
		if err != nil {
			log.Println("accept error:", err)
			time.Sleep(time.Second * 10)
			continue
		}
		go fcgi.Handler(Conn)
	}
}