package http

import (
	"net"
)

type IResponse interface {
	String() string
	Bytes() []byte
}

// http response interface
// getter setter content
// getter setter code
// String and bytes contents
type IHttpResponse interface {
	String() string
	Bytes() []byte
	SetCode(code uint)
	GetCode() uint
	SetContent(content []byte)
	GetContent() []byte
}

func DoResponseWithChannel(res *Response) {

}

func DoResponse(conn net.Conn, buf []byte) {
	conn.Write(buf)
}

func Error(conn net.Conn, err error) {
	DoResponse(conn, []byte(err.Error()))
}

func Success(conn net.Conn, response IResponse) {
	DoResponse(conn, response.Bytes())
}
