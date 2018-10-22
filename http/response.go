package http

import (
	"bytes"
	"net"
)


func Response(conn net.Conn, err error) {
	conn.Write(bytes.NewBufferString(err.Error()).Bytes())
}
