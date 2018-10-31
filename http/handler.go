package http

import (
	"bufio"
	"net"
	"fmt"
)

func Handler(conn net.Conn)  {
	req := &Request{
		KeepConn:false,
		Rwc: bufio.NewReader(conn),
	}

	cf,err := req.Parse()
	if err != nil {
		Error(conn, err)
		return
	}

	ctx := &Context{
		Cf:cf,
		req:req,
		res:make(chan *Response),
		err:make(chan error),
	}
	ctx.Do()

	for{
		select {
			case res := <-ctx.res:
				fmt.Println(res)
			case err := <-ctx.err:
				Error(conn, err)
		}
	}
}
