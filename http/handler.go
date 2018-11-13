package http

import (
	"bufio"
	"net"
	"net/http"
	"errors"
	"reflect"
)

//accept user request with socket
//handler this request and parse http protocol text
//build context obj do any thing
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
		res:make(chan interface{}),
		err:make(chan error),
	}
	ctx.Do()

	for{
		select {
			case res := <-ctx.res:
				switch res.(type) {
					case error:
						Error(conn, err)
						break
					case *http.Response:
						rs := res.(*http.Response)
						rs.Write(conn)
						break
					case *Response:
						Success(conn, res.(*Response))
						break
					default:
						Error(conn, errors.New("not support with type "+reflect.TypeOf(res).String()))
				}
			case err := <-ctx.err:
				Error(conn, err)
		}
	}
}
