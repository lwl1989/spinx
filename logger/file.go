package logger

import (
	"os"
)

type FileLog struct {
	FilePath  string
	SaveLevel Level
}

func (fileLog *FileLog) DoWrite(buf []byte) (n int, err error) {
	fd, _ := os.OpenFile(fileLog.FilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer fd.Close()

	return fd.Write(buf)
}
