package http

import (
	"net"
)

type IResponse interface {
	String() (string)
}

func Response(conn net.Conn, str string) {
	conn.Write([]byte(str))
}

func Error(conn net.Conn, err error) {
	Response(conn, err.Error())
}

func Success(conn net.Conn, response IResponse) {
	Response(conn, response.String())
}
