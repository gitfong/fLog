package fLog

import (
	"os"
	"sync"
	"time"
)

type fLogger struct {
	fileDir    string
	fileName   string
	fileHandle *os.File
	level      int
	nextHour   time.Time
	mu         sync.Mutex
}

const (
	logLevelDebug = iota
	logLevelInfo
	logLevelWarning
	logLevelError
	logLevelFatal
)
