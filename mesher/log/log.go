package log

import (
	"fmt"
	"log"
	"os"
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

var logfile *os.File

func openlogfile(name string) error {

	file, err := os.OpenFile(name, os.O_WRONLY, 0)
	if err != nil {
		file, err = os.Create(name)
		if err != nil {
			fmt.Println("create file error!", err.Error())
			return err
		}
		fmt.Println("create log file!", name)
	} else {
		fileinfo, err := file.Stat()
		if err == nil {
			file.Seek(fileinfo.Size(), 0)
			fmt.Println("append log to file ", file.Name())
		}
	}
	logfile = file

	return nil
}

func SetLogFile(name string) error {

	err := openlogfile(name)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.SetOutput(logfile)

	return nil
}

func Println(level LOG_LEVEL, v ...interface{}) {
	output := fmt.Sprintf("[%s]", loglevel(level))
	output += fmt.Sprintln(v...)
	if level >= ERROR {
		output += fmt.Sprintln(string(debug.Stack()))
	}
	log.Println(output)
}

func Printf(level LOG_LEVEL, format string, v ...interface{}) {
	output := fmt.Sprintf("[%s]", loglevel(level))
	output += fmt.Sprintf(format, v...)
	if level >= ERROR {
		output += fmt.Sprintln(string(debug.Stack()))
	}
	log.Println(output)
}
