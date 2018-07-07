package log

import (
	"fmt"
	"log"
	"os"
)

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
