package http

import (
	"net"
	"time"
	"log"
	"github.com/lwl1989/spinx/http/fcgi"
	"fmt"
	"bufio"
	"io"
)



func HandleConn(conn net.Conn) []byte {
	defer conn.Close()

	//conn.SetDeadline(time.Unix(time.Now().Unix()+5,0))
	//conn.SetReadDeadline(time.Unix(time.Now().Unix()+5,0))

	return fcgi.Read(conn)
}

func newBufioReader(r io.Reader) *bufio.Reader {
	return bufio.NewReader(r)
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
		//get request
		bio := newBufioReader(Conn)
		l, _, err := bio.ReadLine()
		Method, RequestURI, Proto, ok := fcgi.ParseRequestLine(string(l[:]))
		fmt.Println(Method, RequestURI, Proto, ok)

		//Host: localhost:8888
		l, _, err = bio.ReadLine()
		fmt.Println(string(l[:]))
		return
		buf := HandleConn(Conn)
		//var last byte
		//fmt.Println(string(buf))
		last := make([]byte,0)
		var enter byte = 13
		var line byte = 10

		//拆解一个http头
		//第一行 额外处理 以空格为分割
		//GET /dsafddsf/gfdghfdhfghj/jghjhg?dafds HTTP/1.1
		//其他行 以第一个冒号为分割
		//Host: localhost:8888
		//Connection: keep-alive
		//Cache-Control: max-age=0
		//Upgrade-Insecure-Requests: 1
		//User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36
		//Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8
		//Accept-Encoding: gzip, deflate, br
		//Accept-Language: zh-CN,zh;q=0.9
		fmt.Println(string(buf[:]))
		b := make([]byte,0)
		for _,v := range buf {

			l := len(last)
			b = append(b, v)
			if v == line {
				if l == 1 {
					last = append(last, v)
				}
				if l == 3 {
					fmt.Println("http头结束")
					fmt.Println(string(b[:]))
					fmt.Println("正文开始")
					return
				}
			} else {
				if l > 3 {
					last = make([]byte,0)
				}
				if v == enter {
					last = append(last, v)
				}else{
					last = make([]byte,0)
				}
			}
		}
		fmt.Println(string(b[:]))
		return
		go fcgi.GetRequest(buf)
	}
}
