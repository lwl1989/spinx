// Copyright 2012 Junqing Tan <ivan@mysqlab.net> and The Go Authors
// Use of cgi source code is governed by a BSD-style
// Part of source code is from Go fcgi package

// Fix bug: Can't recive more than 1 record untill FCGI_END_REQUEST 2012-09-15
// By: wofeiwo

package core

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"sync"
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

func (cgi *FCGIClient) writeAbortRequest(reqId uint16) error {
	return cgi.writeRecord(typeAbortRequest, reqId, nil)
}

func (cgi *FCGIClient) writeBeginRequest(reqId uint16, role uint8, flags uint8) error {
	b := [8]byte{byte(role >> 8), byte(role), flags}
	return cgi.writeRecord(typeBeginRequest, reqId, b[:])
}

func (cgi *FCGIClient) writeEndRequest(reqId uint16, appStatus int, protocolStatus uint8) error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint32(b, uint32(appStatus))
	b[4] = protocolStatus
	return cgi.writeRecord(typeEndRequest, reqId, b)
}

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

func (cgi *FCGIClient) DoRequest(request *Request) (retout []byte, reterr []byte, err error) {
	pool := GetIdPool(65535)
	reqId  := pool.Alloc()
	defer cgi.writeEndRequest(reqId,200,0)
	defer pool.Release(reqId)
	defer cgi.rwc.Close()

	if request.KeepConn {
		//err = cgi.writeBeginRequest(reqId, roleResponder, 1)
		//todo: is keepAlived
		err = cgi.writeBeginRequest(reqId, roleResponder, 0)
	} else {
		err = cgi.writeBeginRequest(reqId, roleResponder, 0)
	}

	if err != nil {
		return
	}

	err = cgi.writePairs(typeParams, reqId, request.Params)
	if err != nil {
		return
	}
	p := make([]byte,1024)
	n,_ :=request.Stdin.Read(p)
	err = cgi.writeRecord(typeStdin, reqId, p[:n])
	err = cgi.writeRecord(typeStdin, reqId, nil)
	//err = cgi.writeBody(FCGI_STDIN, reqId, request)
	//if err != nil {
	//	return
	//}

	rec := &record{}
	var err1 error

	// recive untill EOF or FCGI_END_REQUEST
	for {
		err1 = rec.read(cgi.rwc)
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
			fallthrough
		default:
			break
		}
	}

	return
}
