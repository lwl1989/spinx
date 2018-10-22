package http

import (
	"bufio"
	"net"
	"strconv"
)

func Handler(conn net.Conn) (req *Request) {
	req = &Request{
		KeepConn:false,
		Rwc: bufio.NewReader(conn),
	}

	err := req.Parse()
	if err != nil {
		Response(conn, err)
	}
	return req
}
