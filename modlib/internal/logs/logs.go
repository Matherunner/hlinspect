package logs

import (
	"errors"
	"os"
	"unsafe"

	"github.com/op/go-logging"
)

/*
#include <windows.h>
*/
import "C"

// DLLLog main logger
var DLLLog = mustInitLogs()

type debugOutputWriter struct{}

func (w debugOutputWriter) Write(p []byte) (n int, err error) {
	s := C.CString(string(p))
	defer C.free(unsafe.Pointer(s))
	C.OutputDebugStringA(s)
	n = len(p)
	return
}

type fileOutputWriter struct {
	file *os.File
}

// func createFileOutputWriter() fileOutputWriter {
// 	file, err := os.OpenFile("hlinspect-debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return fileOutputWriter{}
// 	}
// 	return fileOutputWriter{
// 		file: file,
// 	}
// }

func (w fileOutputWriter) Write(p []byte) (n int, err error) {
	if w.file != nil {
		return w.file.Write(p)
	}
	err = errors.New("log file is not opened")
	return
}

func mustInitLogs() *logging.Logger {
	logger := logging.MustGetLogger("hlinspect")
	writer := debugOutputWriter{}
	// writer := createFileOutputWriter()
	format := logging.MustStringFormatter(`%{time:15:04:05.000000} %{shortfunc} %{level:.4s} %{id:03x}: %{message}`)
	backend := logging.NewLogBackend(writer, "hlinspect ", 0)
	formatter := logging.NewBackendFormatter(backend, format)
	levelled := logging.AddModuleLevel(formatter)
	levelled.SetLevel(logging.DEBUG, "hlinspect")
	logging.SetBackend(levelled)
	return logger
}
