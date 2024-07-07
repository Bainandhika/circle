package logger

import (
	"circle-fiber/app/config"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jasonlvhit/gocron"
)

var (
	Info     *log.Logger
	Debug    *log.Logger
	Warning  *log.Logger
	Error    *log.Logger
	file     *os.File
	fileName string
)

func createLogFile() error {
	date := time.Now().Format("20060102")
	fileName = fmt.Sprintf("%scircle-fiber-%s.log", config.App.LogPath, date)

	var err error
	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	return nil
}

func newLogger() {
	err := createLogFile()
	if err != nil {
		panic(fmt.Errorf("can't create circle log file, err: %v", err))
	}

	flag := log.LstdFlags | log.Llongfile

	Info = log.New(file, "INFO:  ", flag)
	Warning = log.New(file, "WARNING: ", flag)
	Error = log.New(file, "ERROR: ", flag)
	Debug = log.New(file, "DEBUG: ", flag)

	Info.Printf("Circle new start!")
}

func InitializeLoggerScheduler(sch *gocron.Scheduler) {
	err := sch.Every(1).Day().At("00:00").Do(newLogger)
	if err != nil {
		panic(fmt.Errorf("can't create scheduler of circle log file, err: %v", err))
	}
}

func init() {
	newLogger()
}
