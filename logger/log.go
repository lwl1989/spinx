package logger

import (
	"sync"
	"log"
)

type Level uint16

const (
	CRITICAL Level = iota
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

const (
	TypeFile  string = "file"
	TypeMgo          = "mgo"
	TypeMysql        = "mysql"
)

var Logger *log.Logger
var logOnce sync.Once

type Log struct {
	Log DoLog
}

type FileLogConfig interface {
	GetFilePath() interface{}
}

func GetFileLogger(config FileLogConfig) (*log.Logger) {
	logOnce.Do(func() {
		//fmt.Println("logPath:", config.Logs.File.FilePath)
		//todo from config get Logger
		nowLogger := &Log{
			Log: &FileLog{
				FilePath: config.GetFilePath().(string),
			},
		}
		Logger = log.New(nowLogger, "", log.LstdFlags|log.Llongfile)
	})
	return Logger
}

type DoLog interface {
	DoWrite(p []byte) (n int, err error)
}

func (Log *Log) Write(p []byte) (n int, err error) {
	n, err = Log.Log.DoWrite(p)
	return n, err
	//return Log.Log.DoWrite(p)
}
