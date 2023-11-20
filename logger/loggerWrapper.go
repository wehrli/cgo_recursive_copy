package logger

/*
#include "logger.c"
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

type loggerWrapper struct {
	destFile   string
	logChannel chan LogMessage
	waitGroup  sync.WaitGroup
}

type LogMessage struct {
	Text  string
	Level C.LogLevel
}

func NewLoggerWrapper(destFile string) *loggerWrapper {
	lw := &loggerWrapper{
		destFile:   destFile,
		logChannel: make(chan LogMessage),
		waitGroup:  sync.WaitGroup{},
	}

	go lw.startAsyncLogger()

	return lw
}

func (lw *loggerWrapper) startAsyncLogger() {
	defer lw.waitGroup.Done()
	lw.waitGroup.Add(1)

	for message := range lw.logChannel {
		cMessage := C.CString(message.Text)
		cDestFile := C.CString(lw.destFile)
		defer C.free(unsafe.Pointer(cMessage))
		defer C.free(unsafe.Pointer(cDestFile))

		resp := C.WriteLogWithLevel(cMessage, message.Level, cDestFile)
		if resp != 0 {
			fmt.Printf("failed to write log:  %s\n", message.Text)
		}
	}
}

func (lw *loggerWrapper) Debug(message string) {
	lw.logChannel <- LogMessage{message, C.DEBUG}
}

func (lw *loggerWrapper) Info(message string) {
	lw.logChannel <- LogMessage{message, C.INFO}
}

func (lw *loggerWrapper) Warn(message string) {
	lw.logChannel <- LogMessage{message, C.WARNING}
}

func (lw *loggerWrapper) Err(message string) {
	lw.logChannel <- LogMessage{message, C.ERROR}
}

func (lw *loggerWrapper) Close() {
	close(lw.logChannel)
	lw.waitGroup.Wait()
}
