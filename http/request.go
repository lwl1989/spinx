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
	Host, Port string
	Header map[string]string  //必要设置的Header
	KeepConn bool
	content []byte
	Cf	*conf.HostMap
	Method, RequestURI, Proto string
}


func (req *Request) setHeader(key, value string)  {

}


func (req *Request) getHeader(key string) string {
	return ""
}


func (req *Request) getHeaderPosition() (start, end int64)  {
	return start,end
}

//执行此方法 获取 Request URL:
//获取解析到HOST PORT 获得CGI转发的配置
//并且获取到双换行的位置
//然后通过配置获取到的参数 修改头 协议需要的 增加头
//document_file  index等等
func (req *Request) Parse() (e error)  {
	l, _, err := req.Rwc.ReadLine()
	if err != nil {
		return GetError(500, err.Error())
	}

	var ok bool
	req.Method, req.RequestURI, req.Proto, ok = ParseRequestLine(string(l[:]))
	req.Method = strings.ToUpper(req.Method)
	if !ok {
		return GetError(500, err.Error())
	}
	//fmt.Println(Method, RequestURI, Proto, ok)

	//Host: localhost:8888
	l, _, err = req.Rwc.ReadLine()
	if err != nil {
		return GetError(404, err.Error())
	}

	req.Host, req.Port, ok = ParseHostLine(string(l[:]))
	if !ok {
		return GetError(404, err.Error())
	}

	cf,err := conf.HostMaps.GetHostMap(req.Port, req.Host)
	if err != nil {
		return GetError(404, err.Error())
	}

	req.Cf = cf
	//根据配置是走proxy还是走fcgi还是走cache
	return nil
}
func (req *Request) Do() (e error) {
	return nil
}