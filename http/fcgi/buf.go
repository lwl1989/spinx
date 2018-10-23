package fcgi

import (
	"io"
	"log"
	"bufio"
)

const enter byte = 13
const line  byte = 10



func Read(reader io.ReadCloser) []byte {
	buffer := make([]byte, 0)
	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf)
		if n < 1024 {
			buf = buf[0:n]
		}
		if err != nil {
			if err == io.EOF {
				log.Println("Eof", len(buffer), n)
				return buffer
			}
		}
		if n < 1 {
			log.Println("read over")
			return buffer
		}
		buffer = append(buffer, buf[:]...)
	}

	return buffer
}



type bufWriter struct {
	closer io.Closer
	*bufio.Writer
}

// close BuffWrite impl io.Closer
func (b *bufWriter) Close() error {
	if err := b.Writer.Flush(); err != nil {
		b.closer.Close()
		return err
	}
	return b.closer.Close()
}

//get a new buf
func newWriter(c *CgiClient, recType uint8, reqId uint16) *bufWriter {
	s := &streamWriter{c: c, recType: recType, reqId: reqId}
	return &bufWriter{s, bufio.NewWriterSize(s, maxWrite)}
}

// streamWriter abstracts out the separation of a stream into discrete records.
// It only writes maxWrite bytes at a time.
type streamWriter struct {
	c       *CgiClient
	recType uint8
	reqId   uint16
}

//write stream impl io.Writer
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
//write stream impl io.Closer
func (w *streamWriter) Close() error {
	// send empty record to close the stream
	return w.c.writeRecord(w.recType, w.reqId, nil)
}