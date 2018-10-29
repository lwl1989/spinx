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

	pro := &Procotol{
		Cf:cf,
		req:req,
		res:make(chan *Response),
	}
	err = pro.Do()
	if err != nil {
		Error(conn, err)
		return
	}

	for{
		select {
			case res := <-pro.res:
				fmt.Println(res)
		}
	}


	//todo:这里监听协程对象 返回数据
	/**
	for {
		select
	      case: <- error
			Error(conn, error)
	 	  case: <- IResponse
			Success(conn, error)
	}
	 */
}
