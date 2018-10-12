package http

import (
	"net"
	"time"
	"log"
	"github.com/lwl1989/spinx/http/fcgi"
)



//func HandleConn(conn net.Conn) []byte {
//	defer conn.Close()
//
//	//conn.SetDeadline(time.Unix(time.Now().Unix()+5,0))
//	//conn.SetReadDeadline(time.Unix(time.Now().Unix()+5,0))
//
//	return fcgi.Read(conn)
//}


func Do()  {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Println("error listen:", err)
		return
	}
	defer l.Close()
	log.Println("listen ok")

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
