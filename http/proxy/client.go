package proxy

import (
	"github.com/lwl1989/spinx/conf"
	"net/http"
	"io"
	"strings"
)

type ClientProxy struct {
	Req *http.Request
	Res *http.Response
}

type Request struct {
	Rwc io.Reader
	Header map[string]string  //必要设置的Header
	KeepConn bool
	Host, Port, Method, RequestURI, Proto string
	Cf *conf.HostMap
}

func New(req *Request) (client *ClientProxy, err error) {

	reqHttp,err := http.NewRequest(req.Method, req.Cf.Proxy+req.RequestURI, req.Rwc)
	reqHttp.Header = parseHeader(req.Header)

	client = &ClientProxy{
		Req:reqHttp,
	}
	return client,nil
}

func parseHeader(oldHeader map[string]string) http.Header {
	header := make(http.Header)
	for k,v := range oldHeader {
		header[k] = make([]string,0)
		if strings.Index(v, ";") != -1 {
			split := strings.Split(v, ";")
			for _,value := range split {
				header[k] = append(header[k], value)
			}
		}else{
			header[k] = append(header[k], v)
		}
	}

	return header
}

func (client *ClientProxy) DoRequest() (err error) {
	c := &http.Client{}
	client.Res,err = c.Do(client.Req)
	if err != nil {
		return err
	}

	return err
}