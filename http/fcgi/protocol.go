package fcgi

import (
	"sync"
	"io"
	"bytes"
	"net"
	"encoding/binary"
	"net/http"
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

type FCGIClient struct {
	mutex     sync.Mutex
	rwc       io.ReadWriteCloser

	h         header
	buf       bytes.Buffer
}

//get new fcgi proxy
func New(rule ,addr string) (fcgi *FCGIClient, err error) {
	var conn net.Conn
	//fastcgi_pass  127.0.0.1:9000;
	//fastcgi_pass   unix:/dev/shm/php-cgi.sock;
	if rule != "unix" {
		rule = "tcp"
	}
	conn, err = net.Dial(rule,addr)
	fcgi = &FCGIClient{
		rwc:       conn,
	}
	return
}

//write content to proxy
func (cgi *FCGIClient) writeRecord(recType uint8, reqId uint16, content []byte) (err error) {
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
func (cgi *FCGIClient) writeAbortRequest(reqId uint16) error {
	return cgi.writeRecord(typeAbortRequest, reqId, nil)
}

//write fcgi begin flag
func (cgi *FCGIClient) writeBeginRequest(reqId uint16, role uint8, flags uint8) error {
	b := [8]byte{byte(role >> 8), byte(role), flags}
	return cgi.writeRecord(typeBeginRequest, reqId, b[:])
}

//write fcgi end
func (cgi *FCGIClient) writeEndRequest(reqId uint16, appStatus int, protocolStatus uint8) error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint32(b, uint32(appStatus))
	b[4] = protocolStatus
	return cgi.writeRecord(typeEndRequest, reqId, b)
}

//write fcgi header
func (cgi *FCGIClient) writePairs(recType uint8, reqId uint16, pairs map[string]string) error {
	w := newWriter(cgi, recType, reqId)


	b := make([]byte, 8)
	for k, v := range pairs {
		n := encodeSize(b, uint32(len(k)))
		n += encodeSize(b[n:], uint32(len(v)))

		if _, err := w.Write(b[:n]); err != nil {
			return err
		}
		if _, err := w.WriteString(k); err != nil {
			return err
		}
		if _, err := w.WriteString(v); err != nil {
			return err
		}
	}
	return w.Close()
}

//write content with http content
func (cgi *FCGIClient) writeBody(recType uint8, reqId uint16, req *Request) (err error) {
	// write the stdin stream
	stdinWriter := newWriter(cgi, recType, reqId)
	if req.Stdin != nil {
		defer req.Stdin.Close()
		p := make([]byte, 1024)
		var count int
		for {
			count, err = req.Stdin.Read(p)
			if err == io.EOF {
				err = nil
			} else if err != nil {
				stdinWriter.Close()
				return err
			}
			if count == 0 {
				break
			}

			_, err = stdinWriter.Write(p[:count])
			if err != nil {
				stdinWriter.Close()
				return err
			}
		}

	}
	if err = stdinWriter.Close(); err != nil {
		return err
	}
	return nil
}

// bufWriter encapsulates bufio.Writer but also closes the underlying stream when
// Closed.
func (cgi *FCGIClient) GetRequest(r *http.Request, env map[string]string) (req *Request) {
	req = &Request{
		Raw:    r,
		Role:   roleResponder,
		Params: env,
		Stdin:	r.Body,
		KeepConn: r.Header.Get("Connection") == "keep-alive",
	}
	return req
}


//if is proxy request
//do request and get response
func (cgi *FCGIClient) DoRequest(request *Request) (retout []byte, reterr []byte, err error) {
	pool := GetIdPool(65535)
	reqId := pool.Alloc()
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
		return
	}

	//err = cgi.writePairs(typeParams, reqId, request.Params)
	err = cgi.writeRecord(typeParams, reqId, request.buf[request.pos.HStart:request.pos.HEnd])
	if err != nil {
		return
	}

	//p := make([]byte, 1024)
	//n, _ := request.Stdin.Read(p)
	//err = cgi.writeRecord(typeStdin, reqId, p[:n])
	err = cgi.writeRecord(typeStdin, reqId, nil)
	if err != nil {
		return
	}
	rec := &record{}
	var err1 error

	//思路错了
	// 应该是   用户请求->golang http ->my proxy->fastcgi
	// 但是     golang http每个用户请求已经是产生了一个协成然后我获取的已经是一个底层的链接了
	// 只能    自己处理http协议 然后 通过  如果是一个 keep-alive
	// 则   连接保持为长连接 然后 没处理一次请求 max -1 同时产生一个定时器 当定时器到了 也直接return
	// 伪代码：
	// newTimer(timeoutSecond).add()  //对没处理完的链接发出502
	// while(max > 1) {
	//	when user -> request
	//  max --;
	//  connectId add to []connection
	//  go send(connectId) //
	//  go read(connectId) //读取返回
	// }


	// recive untill EOF or FCGI_END_REQUEST
	// todo :if time out  add  Connection: close
	for {
		err1 = rec.read(cgi.rwc)
		//if !keep-alive the end has EOF
		if err1 != nil {
			if err1 != io.EOF {
				err = err1
			}
			break
		}

		switch {
		case rec.h.Type == typeStdout:
			retout = append(retout, rec.content()...)
		case rec.h.Type == typeStderr:
			reterr = append(reterr, rec.content()...)
		case rec.h.Type == typeEndRequest:
			//if keep-alive
			//It's had return
			//But connection Not close
			retout = append(retout, rec.content()...)
			return
		default:
			//fallthrough
		}
	}

	return
}
