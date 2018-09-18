package fcgi

import (
	"net"
	"fmt"
	"os"
)

/**
	主要实现对后端的请求(PHP PYTHON)
 */

type CgiClient struct {
	req         *Request
	res 		[]byte
	response 	net.Conn

}

func GetCgiClient(response net.Conn) *CgiClient {
	return &CgiClient{
		req:GetRequest(),
		res:make([]byte,0),
		response:response,
	}
}

func (fpm *CgiClient) Read(buf []byte) (n int, err error) {
	//todo:数据写往RES
	fpm.res = buf
	return len(buf), nil
}

func (fpm *CgiClient) Write(buf []byte) (n int, err error)  {
	//todo:数据写往REQ
	fpm.req.buf = buf
	return len(buf), nil
}

//请求后端
func (fpm *CgiClient) Proxy() {
	//config conf.Conf
	//todo:获取到server
	server := "127.0.0.1:8888"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	//fmt.Println("connect success")
	conn.Write(fpm.req.buf)

	fpm.res = Read(conn)
}

//发送结果接请求
func (fpm *CgiClient) Response()  {
	defer fpm.response.Close()
	fpm.response.Write(fpm.res)
}