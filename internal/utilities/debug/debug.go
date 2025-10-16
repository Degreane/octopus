package debug

import (
	"log"
	"os"
)

const (
	Blue   = "\033[34m"
	Orange = "\033[33m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Reset  = "\033[0m"
)

var DEBUG bool

type Level int

const (
	Info Level = iota
	Important
	Warning
	Error
)

func Debug(level Level, v ...any) {
	if level == Warning {
		warning(v...)
	} else if level == Error {
		err(v...)
	} else if level == Important {
		important(v...)
	} else if level == Info && (os.Getenv("DEBUG") == "true" || DEBUG) {
		info(v...)
	}
}

func info(v ...interface{}) {
	log.Printf(Blue+"[INFO] "+Reset+" %+v", v...)
}

func warning(v ...interface{}) {
	log.Printf(Orange+"[WARNING] "+Reset+" %+v", v...)
}

func err(v ...interface{}) {
	log.Printf(Red+"[ERROR] "+Reset+" %+v", v...)
}
func important(v ...interface{}) {
	log.Printf(Green+"[IMPORTANT] "+Reset+" %+v", v...)
}
