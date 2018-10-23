package http

import (
	"bufio"
	"net"
	"github.com/lwl1989/spinx/http/fcgi"
)

func Handler(conn net.Conn)  {
	req := &Request{
		KeepConn:false,
		Rwc: bufio.NewReader(conn),
	}

	err := req.Parse()
	if err != nil {
		Error(conn, err)
		return
	}

	if req.Cf.Proxy != "" {
		//do proxy

		return
	}

	cgi,err := fcgi.New(req)
	if err != nil {
		Error(conn, err)
		return
	}

	go cgi.DoRequest()
	//todo:这里监听协程对象 返回数据
}
