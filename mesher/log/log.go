package log

import (
	"log"
	"os"
)

var logfile *os.File

func SetLogFile(name string) error {

	var err error

	logfile, err = os.Open(name)
	if err == nil {

		fstat, err := logfile.Stat()
		if err != nil {
			log.Println(err.Error())
			return err
		}

		logfile.Seek(fstat.Size(), 0)

		log.SetOutput(logfile)
		return nil
	}

	if os.IsNotExist(err) {
		logfile, err := os.Create(name)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		log.SetOutput(logfile)
		return nil
	}
	log.Println(err.Error())

	return err
}
