package fcgi

import (
	"bufio"
	"github.com/lwl1989/spinx/conf"
)

type Request struct {
	Id  uint16
	Rwc *bufio.Reader
	Host, Port string
	Header map[string]string  //必要设置的Header
	KeepConn bool
	content []byte
	Method, RequestURI, Proto string
	Cf *conf.HostMap
}
