package http

import (
	"errors"
	"github.com/lwl1989/spinx/conf"
	"github.com/lwl1989/spinx/http/fcgi"
	"github.com/lwl1989/spinx/http/proxy"
)

//context
//has conf
//build req
//received chan res or chan error
//do Response
type Context struct {
	Cf	*conf.HostMap
	req *Request
	res chan interface{}
	err chan error
}

// do ctxcotol
// check config and get how to do
// ctxxy ? cache ? fastcgi ?
//
func (ctx *Context) Do() {
	if ctx.Cf.Proxy != "" {
		ctx.DoProxy()
	}

	if ctx.Cf.CacheRule != "" {
		ctx.DoCache()
	}

	if ctx.Cf.CgiProxy != "" {
		ctx.DoCgi()
	}else {
		ctx.err <- errors.New("can't do this ctxcotol")
	}
}

// cache
func (ctx *Context) DoCache()  {

}
// read config and  ctxxy
func (ctx *Context) DoProxy() {
	req := &proxy.Request{
		Cf:ctx.Cf,
		Rwc: ctx.req.Rwc,
		Method: ctx.req.Method,
		Host: ctx.req.Host,
		Port: ctx.req.Port,
		Header: make(map[string]string),
		KeepConn:ctx.req.KeepConn,
		RequestURI:ctx.req.RequestURI,
		Proto:ctx.req.Proto,
	}
	cgi,err := proxy.New(req)
	if err != nil {
		ctx.err <- err
		return
	}
	content, err := cgi.DoRequest()
	if err != nil {
		ctx.err <- err
		return
	}
	ctx.res <- &Response{
		content: content,
	}

}

// read config and build cgi ctxtocol
func (ctx *Context) DoCgi() {
	req := &fcgi.Request{
		Cf:ctx.Cf,
		Rwc: ctx.req.Rwc,
		Method: ctx.req.Method,
		Host: ctx.req.Host,
		Port: ctx.req.Port,
		Header: make(map[string]string),
		KeepConn:ctx.req.KeepConn,
		RequestURI:ctx.req.RequestURI,
		Proto:ctx.req.Proto,
	}
	cgi,err := fcgi.New(req)
	if err != nil {
		ctx.err <- err
		return
	}

	content, err := cgi.DoRequest()
	if err != nil {
		ctx.err <- err
		return
	}
	ctx.res <- &Response{
		content: content,
	}
}
