package http

import (
	"net"
)

type IResponse interface {
	String() string
	Bytes() []byte
}

func Response(conn net.Conn, buf []byte) {
	conn.Write(buf)
}

func Error(conn net.Conn, err error) {
	Response(conn, []byte(err.Error()))
}

func Success(conn net.Conn, response IResponse) {
	Response(conn, response.Bytes())
}
