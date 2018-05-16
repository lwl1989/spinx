package core

import (
	"net/http"
	"io"
	"sync"
)

//User request
type Request struct {
	Raw      *http.Request
	Role     uint8
	Params   map[string]string
	Stdin    io.ReadCloser
	Data     io.ReadCloser
	KeepConn bool
}

// Request hold information of a standard
// FastCGI request
type IdPool struct {
	IDs chan uint16
}

// AllocID implements Client.AllocID
func (p *IdPool) Alloc() uint16 {
	return <-p.IDs
}

// ReleaseID implements Client.ReleaseID
func (p *IdPool) Release(id uint16) {
	go func() {
		// release the ID back to channel for reuse
		// use goroutine to prev0, ent blocking ReleaseID
		p.IDs <- id
	}()
}


var pool *IdPool
var one sync.Once
func GetIdPool(limit uint16) *IdPool {
	one.Do(func(){
		var ids = make(chan uint16)

		go func(maxID uint16) {
			for i := uint16(0); i < maxID; i++ {
				ids <- i
			}
			ids <- uint16(maxID)
		}(uint16(limit - 1))
		pool =  &IdPool{IDs:ids}
	})
	//sync.Once.Do(func() {
	//
	//})
	return pool
}
