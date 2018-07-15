package log

import (
	"fmt"
	"log"
	"runtime/debug"
)

type LOG_LEVEL int

const (
	INFO LOG_LEVEL = iota
	WARNING
	ERROR
	EXCEPT
)

func loglevel(level LOG_LEVEL) string {
	if level == INFO {
		return "INFO"
	} else if level == WARNING {
		return "WARNING"
	} else if level == ERROR {
		return "ERROR"
	} else {
		return "EXCEPT"
	}
}

var gFileLog *LogFile

func SetLogFile(name string) error {
	file, err := NewLogFile(name, "./", 1024*1024, 5)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.SetOutput(file)
	log.SetFlags(log.Lmicroseconds | log.LstdFlags)
	gFileLog = file
	return nil
}

func Println(level LOG_LEVEL, v ...interface{}) {
	output := fmt.Sprintf("[%s]", loglevel(level))
	output += fmt.Sprint(v...)
	if level >= ERROR {
		output += fmt.Sprint(string(debug.Stack()))
	}
	log.Println(output)
}

func Printf(level LOG_LEVEL, format string, v ...interface{}) {
	output := fmt.Sprintf("[%s]", loglevel(level))
	output += fmt.Sprintf(format, v...)
	if level >= ERROR {
		output += fmt.Sprint(string(debug.Stack()))
	}
	log.Printf(output)
}
