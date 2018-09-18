package http

import (
	"net"
	"time"
	"log"
	"github.com/lwl1989/spinx/http/fcgi"
)



func HandleConn(conn net.Conn) []byte {
	defer conn.Close()

	conn.SetDeadline(time.Unix(time.Now().Unix()+5,0))
	conn.SetReadDeadline(time.Unix(time.Now().Unix()+5,0))

	return fcgi.Read(conn)
}

func Do()  {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Println("error listen:", err)
		return
	}
	defer l.Close()
	log.Println("listen ok")

	var i int
	for {
		time.Sleep(time.Second * 10)
		if Conn, err := l.Accept(); err != nil {
			log.Println("accept error:", err)
			break
		}else{
			buf := HandleConn(Conn)
			cgiClient := fcgi.GetCgiClient(Conn)
			cgiClient.Proxy()
			cgiClient.Response()
			log.Println(string(buf[:]))
		}
		i++
		log.Printf("%d: accept a new connection\n", i)
	}
}
