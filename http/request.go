package http

import (
	"bufio"
	"github.com/lwl1989/spinx/conf"
	"strings"
)

const enter byte = 13
const line  byte = 10
//用于获取头的位置
const SPLIT_STR  =  "\n\n"
//用户获取特定头
const ENTER_SPACE  = "\n"
//请求头部字符串
const REQUEST_URL  = "Request URL:"


type Request struct {
	Id  uint16
	Rwc *bufio.Reader
	KeepConn bool
	Host, Port, Method, RequestURI, Proto string
}



//执行此方法 获取 Request URL:
//获取解析到HOST PORT 获得CGI转发的配置
//并且获取到双换行的位置
//然后通过配置获取到的参数 修改头 协议需要的 增加头
//document_file  index等等
func (req *Request) Parse() (cf *conf.HostMap, e error)  {
	l, _, err := req.Rwc.ReadLine()
	if err != nil {
		return nil,GetError(500, err.Error())
	}

	var ok bool
	req.Method, req.RequestURI, req.Proto, ok = ParseRequestLine(string(l[:]))
	req.Method = strings.ToUpper(req.Method)
	if !ok {
		return nil,GetError(500, err.Error())
	}
	//fmt.Println(Method, prouestURI, Proto, ok)

	//Host: localhost:8888
	l, _, err = req.Rwc.ReadLine()
	if err != nil {
		return nil,GetError(404, err.Error())
	}

	req.Host, req.Port, ok = ParseHostLine(string(l[:]))
	if !ok {
		return nil,GetError(404, err.Error())
	}


	cf,err = conf.HostMaps.GetHostMap(req.Port, req.Host)
	if err != nil {
		return nil,GetError(404, err.Error())
	}

	//根据配置是走proxy还是走fcgi还是走cache
	return cf,nil
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
