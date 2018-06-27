package http

import (
	"net"
	"time"
	"sync"
)

type Server struct {
	Config 	map[string]interface{}
	ServerConfig map[string]ServerConfig
	listeners  map[net.Listener]struct{}

	mu sync.Mutex
	rwc net.Conn
	proxy net.Conn
}

type ServerConfig struct {
	Name string
	Net string
	Addr string
	TryFiles string
	DocumentRoot string
	Index string
}
// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func (srv *Server) ListenAndServe(addr string) error {

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}

func (srv *Server) Serve(l net.Listener) error {

	defer l.Close()

}

func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()


}