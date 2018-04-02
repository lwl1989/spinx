package core

import (
	"io"
	"bufio"
)

type bufWriter struct {
	closer io.Closer
	*bufio.Writer
}

func (b *bufWriter) Close() error {
	if err := b.Writer.Flush(); err != nil {
		b.closer.Close()
		return err
	}
	return b.closer.Close()
}


func newWriter(c *FCGIClient, recType uint8, reqId uint16) *bufWriter {
	s := &streamWriter{c: c, recType: recType, reqId: reqId}
	return &bufWriter{s, bufio.NewWriterSize(s, maxWrite)}
}

// streamWriter abstracts out the separation of a stream into discrete records.
// It only writes maxWrite bytes at a time.
type streamWriter struct {
	c       *FCGIClient
	recType uint8
	reqId   uint16
}

func (w *streamWriter) Write(p []byte) (int, error) {
	nn := 0
	for len(p) > 0 {
		n := len(p)
		if n > maxWrite {
			n = maxWrite
		}

		if w.recType == 5 {
			//fmt.Println(w.recType, w.reqId, p[:])
			//str := "dbname=wordpress&uname=username&pwd=password&dbhost=localhost&prefix=wp_&language=&submit=Submit"

			if err := w.c.writeRecord(w.recType, w.reqId, p); err != nil {
				return nn, err
			}
		} else {
			if err := w.c.writeRecord(w.recType, w.reqId, p[:n]); err != nil {
				return nn, err
			}
		}

		nn += n
		p = p[n:]
	}
	return nn, nil
}

func (w *streamWriter) Close() error {
	// send empty record to close the stream
	return w.c.writeRecord(w.recType, w.reqId, nil)
}