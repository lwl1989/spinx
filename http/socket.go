package http

import (
	"net"
	"time"
	"log"

	"io"
	"bytes"
)

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

func HandleConn(conn net.Conn) []byte {
	defer conn.Close()

	conn.SetDeadline(time.Unix(time.Now().Unix()+5,0))
	conn.SetReadDeadline(time.Unix(time.Now().Unix()+5,0))

	buffer := make([]byte, 0)
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if n < 1024 {
			buf = buf[0:n]
		}
		if err != nil {
			if err == io.EOF {
				log.Println("Eof", len(buffer), n)
				//conn.Write([]byte("<h1>helloWrold!</h1>"))
				return buffer
			}
			log.Println(conn.RemoteAddr().String(), " connection error: ", err)
		}
		if n < 1 {
			log.Println("read over")
			//conn.Write([]byte("<h1>helloWrold!</h1>"))
			return buffer
		}
		buffer = append(buffer, buf[:]...)

		//log.Println(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))
	}

	//conn.Write([]byte("<h1>helloWrold!</h1>"))
	return buffer
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
			log.Println(string(buf[:]))
		}
		i++
		log.Printf("%d: accept a new connection\n", i)
	}
}
