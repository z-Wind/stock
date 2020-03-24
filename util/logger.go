package util

import (
	"fmt"
	"log"
	"os"
)

// Log custom
type Log struct {
	f *os.File
}

// Start write log to file
func (l *Log) Start(name string, flag int, perm os.FileMode) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	l.f = f
	log.SetOutput(l.f)
}

// SetFlags set flags
func (l *Log) SetFlags(flag int) {
	log.SetFlags(flag)
}

// Printf show on screen and write to log file
func (l *Log) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
	log.Printf(format, v...)
}

// LPrintf only write to log file
func (l *Log) LPrintf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Fatalf show on screen and write to log file
func (l *Log) Fatalf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
	log.Fatalf(format, v...)
}

// Stop stop to write log
func (l *Log) Stop() {
	l.f.Close()
}
