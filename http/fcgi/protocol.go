package fcgi

import (
	"sync"
	"io"
	"bytes"
	"net"
	"encoding/binary"
	"time"
	"errors"
	"bufio"
	"github.com/lwl1989/spinx/conf"
	"strings"
	"strconv"
)

const (
	typeBeginRequest    uint8 = 1
	typeAbortRequest    uint8 = 2
	typeEndRequest      uint8 = 3
	typeParams          uint8 = 4
	typeStdin           uint8 = 5
	typeStdout          uint8 = 6
	typeStderr          uint8 = 7
	typeData            uint8 = 8
	typeGetValues       uint8 = 9
	typeGetValuesResult uint8 = 10
	typeUnknownType     uint8 = 11
)

// keep the connection between web-server and responder open after request
const flagKeepConn = 1

const (
	roleResponder = iota + 1 // only Responders are implemented.
	roleAuthorizer
	roleFilter
)

const (
	statusRequestComplete = iota
	statusCantMultiplex
	statusOverloaded
	statusUnknownRole
)

type CgiClient struct {
	mutex     sync.Mutex
	rwc       io.ReadWriteCloser
	h         header
	buf       bytes.Buffer
	request   *Request
}

//get new fcgi proxy
func New(rule ,addr string) (cgi *CgiClient, err error) {
	var conn net.Conn

	if rule != "unix" {
		rule = "tcp"
	}
	conn, err = net.DialTimeout(rule, addr, 3*time.Second)
	cgi = &CgiClient{
		rwc:       conn,
	}
	return
}

//write content to proxy
func (cgi *CgiClient) writeRecord(recType uint8, reqId uint16, content []byte) (err error) {
	cgi.mutex.Lock()
	defer cgi.mutex.Unlock()
	cgi.buf.Reset()
	cgi.h.init(recType, reqId, len(content))

	if err := binary.Write(&cgi.buf, binary.BigEndian, cgi.h); err != nil {
		return err
	}
	if _, err := cgi.buf.Write(content); err != nil {
		return err
	}
	if _, err := cgi.buf.Write(pad[:cgi.h.PaddingLength]); err != nil {
		return err
	}
	_, err = cgi.rwc.Write(cgi.buf.Bytes())
	return err
}

//write fcgi abort flag
func (cgi *CgiClient) writeAbortRequest(reqId uint16) error {
	return cgi.writeRecord(typeAbortRequest, reqId, nil)
}

//write fcgi begin flag
func (cgi *CgiClient) writeBeginRequest(reqId uint16, role uint8, flags uint8) error {
	b := [8]byte{byte(role >> 8), byte(role), flags}
	return cgi.writeRecord(typeBeginRequest, reqId, b[:])
}

//write fcgi end
func (cgi *CgiClient) writeEndRequest(reqId uint16, appStatus int, protocolStatus uint8) error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint32(b, uint32(appStatus))
	b[4] = protocolStatus
	return cgi.writeRecord(typeEndRequest, reqId, b)
}

//write fcgi header
func (cgi *CgiClient) writeHeader(recType uint8, reqId uint16, req *Request) (err error) {
	writer := newWriter(cgi, recType, reqId)
	defer writer.Close()

	headers := make(map[string]string)
	for ; ; {
		by := make([]byte, 0)
		by,_,err := req.rwc.ReadLine()
		if len(by) == 0 {
			break
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		str := string(by[:])
		s1 := strings.Index(str, ":")
		if s1>0 {
			headers[str[:s1]] = str[s1+1:]
		}
		//by = append(by, 13, 10)
	}

	err,pairs := buildEnv(req)
	if err != nil {
		return err
	}

	for k,v := range headers {
		pairs["HTTP_"+strings.Replace(strings.ToUpper(k), "-", "_", -1)] = v
	}

	//_, err := writer.Write(req.content[:])
	b := make([]byte, 8)
	for k, v := range pairs {
		//fmt.Println(k,v)
		n := encodeSize(b, uint32(len(k)))
		n += encodeSize(b[n:], uint32(len(v)))
		//fmt.Println(k,v,b[:n])
		if _, err := writer.Write(b[:n]); err != nil {
			return err
		}
		if _, err := writer.WriteString(k); err != nil {
			return err
		}
		if _, err := writer.WriteString(v); err != nil {
			return err
		}
	}
	return err
}

//write content with http content
func (cgi *CgiClient) writeBody(recType uint8, reqId uint16, req *Request) (err error) {
	// write the stdin stream
	writer := newWriter(cgi, recType, reqId)
	defer writer.Close()
	l,err := strconv.Atoi(req.Header["CONTENT_LENGTH"])
	if err != nil {
		l = 0
	}
	if l > 0 {
		p := make([]byte, 1024)
		var count int
		for {
			count, err = req.rwc.Read(p)
			if err == io.EOF {
				err = nil
			} else if err != nil {
				return err
			}
			if count == 0 {
				break
			}
			_, err = writer.Write(p[:count])
			if err != nil {
				return err
			}
		}
	}

	return err
}

// bufWriter encapsulates bufio.Writer but also closes the underlying stream when
// Closed.
func Handler(conn net.Conn) {
	pool := GetIdPool(65535)
	reqId := pool.Alloc()
	//close connection and release id
	defer func() {
		conn.Close()
		pool.Release(reqId)
	}()


	req := &Request{
		Id: reqId,
		KeepConn:false,
		rwc: bufio.NewReader(conn),
	}

	l, _, err := req.rwc.ReadLine()
	if err != nil {
		Response(conn, "500", "")
		return
	}
	var ok bool
	req.Method, req.RequestURI, req.Proto, ok = ParseRequestLine(string(l[:]))
	req.Method = strings.ToUpper(req.Method)
	if !ok {
		Response(conn, "500", "")
		return
	}
	//fmt.Println(Method, RequestURI, Proto, ok)

	//Host: localhost:8888
	l, _, err = req.rwc.ReadLine()
	if err != nil {
		Response(conn, "404", "")
		return
	}

	req.Host, req.Port, ok = ParseHostLine(string(l[:]))
	if !ok {
		Response(conn, "404", "")
		return
	}

	cf,err := conf.HostMaps.GetHostMap(req.Port, req.Host)
	if err != nil {
		Response(conn, "404", "")
		return
	}

	req.cf = cf
	//i := 0
	//for ; ;  {
	//	l, _, err = req.rwc.ReadLine()
	//	//正文这里结束了，因为readLine会跳过换行
	//	if len(l) == 0 {
	//		i++
	//		break
	//	}
	////	fmt.Println(string(l[:]))
	//}
	//HEADER END \r\n
	//\r\n  continue
	//BODY
	//因此跳过的行数是1
	//fmt.Println(i)
	//for ; ; {
	//	b, e := req.rwc.ReadByte()
	//	fmt.Println(b,e)
	//}

	//response.Body = bytes.NewReader()
	//处理完成 从这里获取 host:port 得到配置  然后处理request uri
	//最后进行转发

	cgi,err := New(cf.Net, cf.Addr)
	cgi.request = req
	if err != nil {
		Response(conn, "502", err.Error())
		return
	}

	content, err := cgi.DoRequest(req)

	if err != nil {
		Response(conn, "500", err.Error())
		return
	}

	content = append(bytes.NewBufferString("HTTP/1.1 200").Bytes(), content...)
	conn.Write(content)
}

func Response(conn net.Conn, code, content string) {
	conn.Write(bytes.NewBufferString("HTTP/1.1  "+code+" \r\n\r\n<h1>"+code+"</h1><p>"+content+"</p>").Bytes())
}

//if is proxy request
//do request and get response
func (cgi *CgiClient) DoRequest(request *Request) (retout []byte, err error) {
	reqId := request.Id
	defer cgi.writeEndRequest(reqId, 200, 0)
	defer pool.Release(reqId)
	defer cgi.rwc.Close()

	if request.KeepConn {
		//if it's keep-alive
		//set flags 1
		err = cgi.writeBeginRequest(reqId, roleResponder, 1)
	} else {
		err = cgi.writeBeginRequest(reqId, roleResponder, 0)
	}

	if err != nil {
		return nil, err
	}

	err = cgi.writeHeader(typeParams, reqId, request)
	if err != nil {
		return nil, err
	}
	//todo: 这个时间应该从配置中读取
	timer := time.NewTimer(1*time.Second)

	if request.Method != "GET" && request.Method != "HEAD" {
		err = cgi.writeBody(typeStdin, reqId, request)
		if err != nil {
			return nil, err
		}
	}

	res := &ResponseContent{
		received:make(chan bool),
		err:make(chan error),
		buf:make([]byte,0),
	}
	go readResponse(cgi, res)
	// recive untill EOF or FCGI_END_REQUEST
	// todo :if time out  add  Connection: close
	for {
		select {
		case <- timer.C:
			//超时发送终止请求
			cgi.writeEndRequest(reqId, 502, 1)
		    err = errors.New("502 timeout")
			return retout,err
		case <-res.received:
			retout = res.content()
			//fmt.Println(string(retout[:])+" has received")
			//fmt.Println(string(res.buf[:]))
			return retout,err
		case e:= <-res.err:
			err = e
			return retout,err
		}
	}

	return retout,err
}

func readResponse(cgi *CgiClient,res *ResponseContent) {
	//bb := Read(cgi.rwc)
	//fmt.Println(string(bb[:]))
	for {
		rec := &record{}
		err1 := rec.read(cgi.rwc)
		//if !keep-alive the end has EOF
		if err1 != nil {
			if err1 != io.EOF {
				res.err <- err1
			} else {
				res.received <- true
			}
			break
		}
		//fmt.Println(rec.h.Type)
		switch {
		case rec.h.Type == typeStdout:
			res.buf = append(res.buf, rec.content()...)
			//fmt.Println(string(rec.buf[:]))
		case rec.h.Type == typeStderr:
			//fmt.Println(string(rec.buf[:]))
			res.buf = append(res.buf, rec.content()...)
		case rec.h.Type == typeEndRequest:
			//if keep-alive
			//It's had return
			//But connection Not close
			//fmt.Println("end")
			//fmt.Println(string(rec.buf[:]))
			res.buf = append(res.content(), rec.content()...)
			res.received <- true
			return
		default:
			//fallthrough
		}
	}
}
