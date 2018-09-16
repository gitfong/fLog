package fLog

import (
	"fmt"
	"log"
	"os"
	"path"
	//"sync"
	"encoding/json"
	"io/ioutil"
	"time"
)

type FLogger struct {
	pathName   string
	fileName   string
	fileHandle *os.File
	baseLogger *log.Logger
	level      int
	nextHour   time.Time
	flushImmediately bool
	//mu         sync.Mutex
}

type jsonCfg struct {
	PathName string `json:"pathName"`
	Level    int    `json:"level"`
	FlushImmediately bool `json:"flushImmediately"`
}

const (
	levelDebug = iota
	levelInfo
	levelWarn
	levelError
	levelFatal
)

const (
	flagLevelDebug = "[debug]"
	flagLevelInfo  = "[info ]"
	flagLevelWarn  = "[warn ]"
	flagLevelError = "[error]"
	flagLevelFatal = "[fatal]"
)

func New() *FLogger {
	cfg := getCfg()

	if cfg.Level > levelFatal || cfg.Level < levelDebug {
		return nil
	}

	flog := new(FLogger)
	flog.pathName = cfg.PathName
	flog.level = cfg.Level
	flog.flushImmediately = cfg.FlushImmediately
	err := flog.makeOutFile()
	if err != nil {
		return nil
	}

	return flog
}

func getCfg() jsonCfg {
	data, err := ioutil.ReadFile("logCfg.json")
	if err != nil {
		fmt.Println("read logCfg.json fail. err:", err)
		os.Exit(0)
	}

	cfg := jsonCfg{}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Println("json unmarshal err:", err)
		os.Exit(0)
	}

	return cfg
}

func (l *FLogger) makeOutFile() error {
	fmt.Println("on makeoutfile")
	if l.fileHandle != nil{
		l.fileHandle.Sync()
		l.fileHandle.Close()
	}

	now := time.Now()
	fileName := fmt.Sprintf("%d%02d%02d_%02d.log", now.Year(), now.Month(), now.Day(), now.Hour())
	fullName := path.Join(l.pathName, fileName)
	//file,err := os.Create(fullName)
	file, err := os.OpenFile(fullName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}

	nextHour := time.Unix(now.Unix()+3600, 0)
	nextHour = time.Date(nextHour.Year(), nextHour.Month(), nextHour.Day(), nextHour.Hour(), 0, 0, 0, nextHour.Location())
	l.nextHour = nextHour
	l.baseLogger = log.New(file, "", log.LstdFlags|log.Lmicroseconds |log.Llongfile)
	l.fileName = fileName
	l.fileHandle = file

	return nil
}

func (logger *FLogger) log(flagLevel string, format string, a ...interface{}) {
	now := time.Now()
	if logger.fileHandle == nil || now.Unix() > logger.nextHour.Unix() {
		logger.makeOutFile()
	}

	logger.baseLogger.SetPrefix(flagLevel)

	str := fmt.Sprintf(format, a...)
	logger.baseLogger.Output(3, str)

	if logger.flushImmediately {
		logger.fileHandle.Sync()
	}

	if flagLevel == flagLevelFatal {
		os.Exit(1)
	}
}

func (logger *FLogger) Debug(format string, a ...interface{}) {
	if levelDebug < logger.level {
		return
	}
	logger.log(flagLevelDebug, format, a...)
}

func (logger *FLogger) Info(format string, a ...interface{}) {
	if levelInfo < logger.level {
		return
	}
	logger.log(flagLevelInfo, format, a...)
}

func (logger *FLogger) Warn(format string, a ...interface{}) {
	if levelWarn < logger.level {
		return
	}
	logger.log(flagLevelWarn, format, a...)
}

func (logger *FLogger) Error(format string, a ...interface{}) {
	if levelError < logger.level {
		return
	}
	logger.log(flagLevelError, format, a...)
}

func (logger *FLogger) Fatal(format string, a ...interface{}) {
	if levelFatal < logger.level {
		return
	}
	logger.log(flagLevelFatal, format, a...)
}
