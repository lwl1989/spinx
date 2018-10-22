package http

import (
	"bytes"
	"net"
)

type ResponseContent struct {
	received chan bool
	err chan error
	buf []byte
}

func (res *ResponseContent) content() []byte  {
	return res.buf
}

func Response(conn net.Conn, err error) {
	conn.Write(bytes.NewBufferString(err.Error()).Bytes())
}
