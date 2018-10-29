package http

import (
	"bufio"
	"net/textproto"
	"sync"
)

var textprotoReaderPool sync.Pool

func newTextprotoReader(br *bufio.Reader) *textproto.Reader {
	if v := textprotoReaderPool.Get(); v != nil {
		tr := v.(*textproto.Reader)
		tr.R = br
		return tr
	}
	return textproto.NewReader(br)
}

func readRequest(b *bufio.Reader, deleteHostHeader bool) (err error) {
	return nil
}

