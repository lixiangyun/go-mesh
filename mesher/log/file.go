package log

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type LogFile struct {
	sync.Mutex

	name    string /* 日志文件名称 */
	dir     string /* 日志所在的目录 */
	maxsize int64  /* 文件上限 */
	maxnum  int64  /* 文件数量上限 */

	file    *os.File /* 当前正在写入的文件句柄 */
	cursize int64    /* 当前文件大小 */
}

func timeStampGet() string {
	tm := time.Now()
	return fmt.Sprintf("%4d%02d%02d%02d%02d%02d%3d",
		tm.Year(), tm.Month(), tm.Day(),
		tm.Hour(), tm.Minute(), tm.Second(),
		time.Duration(tm.Nanosecond())/time.Millisecond)
}

func NewLogFile(name string, dir string, size int64, num int64) (*LogFile, error) {

	logfile := &LogFile{name: name, dir: dir, maxsize: size, maxnum: num}

	fileinfo, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	if !fileinfo.IsDir() {
		return nil, errors.New(name + "is not dir.")
	}

	return logfile, nil
}

func (lf *LogFile) packfile() error {
	fileold := fmt.Sprintf("%s/%s.log", lf.dir, lf.name)
	filenew := fmt.Sprintf("%s/%s_%s.log", lf.dir, lf.name, timeStampGet())
	return os.Rename(fileold, filenew)
}

func (lf *LogFile) openfile() error {
	path := fmt.Sprintf("%s/%s.log", lf.dir, lf.name)

	file, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		file, err = os.Create(path)
		if err != nil {
			return err
		}
		lf.cursize = 0
	} else {
		fileinfo, err := file.Stat()
		if err == nil {
			file.Seek(fileinfo.Size(), 0)
			lf.cursize = fileinfo.Size()
		} else {
			lf.cursize = 0
		}
	}
	lf.file = file
	return nil
}

func (lf *LogFile) Write(p []byte) (n int, err error) {

	lf.Lock()
	defer lf.Unlock()

	for {
		if lf.file == nil {
			err = lf.openfile()
			if err != nil {
				return 0, err
			}
		} else {
			if lf.cursize > lf.maxsize {
				lf.file.Close()
				lf.file = nil
				err = lf.packfile()
				if err != nil {
					return 0, err
				}
			} else {
				os.Stderr.Write(p)
				return lf.file.Write(p)
			}
		}
	}

	return 0, nil
}
