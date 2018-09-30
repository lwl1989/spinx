package http

import (
	"net"
	"time"
	"log"
	"github.com/lwl1989/spinx/http/fcgi"
	"fmt"
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

	for {
		Conn, err := l.Accept()
		//time.Sleep(time.Second * 10)
		if err != nil {
			log.Println("accept error:", err)
			time.Sleep(time.Second * 10)
			continue
		}
		buf := HandleConn(Conn)
		//fmt.Println(string(buf))
		var last byte
		for _,v := range buf {
			b := make([]byte,0)

			if last == byte(10) && last == v {
				fmt.Println("http头结束")
				fmt.Println(string(b[:]))
				fmt.Println("正文开始")
			}
			fmt.Println(last,v)
			last = v
		}
		return
		go fcgi.GetRequest(buf)
	}
}
