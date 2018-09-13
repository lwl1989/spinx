package fcgi

import (
	"io"
	"log"
)

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
