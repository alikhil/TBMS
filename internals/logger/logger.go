package logger

import (
	"io"
	// "io/ioutil"
	"log"
	"os"
)

var (
	// Trace - Very low level messages
	Trace *log.Logger
	// Info - something that is could be useful for
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func init() {
	Init(os.Stdout /*ioutil.Discard*/, os.Stdout, os.Stdout, os.Stderr)
	Info.Printf("Logger has been initialized")
}
