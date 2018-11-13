package http

import (
	"net"
	"net/http"
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

type Response struct {
	content []byte
}

func (response *Response) Bytes() []byte {
	return nil
}

func (response *Response) String() string {
	return ""
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
