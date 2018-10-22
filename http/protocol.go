package http

import (
	"bufio"
	"net/textproto"
	"sync"
	"strings"
)

var textprotoReaderPool sync.Pool

func newTextprotoReader(br *bufio.Reader) *textproto.Reader {
	if v := textprotoReaderPool.Get(); v != nil {
		tr := v.(*textproto.Reader)
		tr.R = br
		return tr
	}
	return textproto.NewReader(br)
}

func readRequest(b *bufio.Reader, deleteHostHeader bool) (err error) {
	return nil
}

// parseRequestLine parses "GET /foo HTTP/1.1" into its three parts.
func ParseRequestLine(line string) (method, requestURI, proto string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	s2 += s1 + 1
	return line[:s1], line[s1+1 : s2], line[s2+1:], true
}

//parse host line "HOST: localhost:8000"
func ParseHostLine(line string) (host, port string, ok bool) {
	port = "80"
	s1 := strings.Index(line, " ")
	if s1 < 0 {
		return
	}
	s2 := strings.Index(line[s1+1:], ":")

	if s2 < 0 {
		host = line[s1+1:]
		return host,port,true
	}
	return line[s1+1:s2+s1+1],line[s2+s1+2:],true
}
