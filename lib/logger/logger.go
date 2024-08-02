package logger

import (
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

type LoggerConfig struct {
	LogPath string
}

func (l *LoggerConfig) createLogFile(logPath string) error {
	date := time.Now().Format("20060102")
	fileName = fmt.Sprintf("%scircle-%s.log", logPath, date)

	var err error
	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (l *LoggerConfig) InitLogger() {
	err := l.createLogFile(l.LogPath)
	if err != nil {
		log.Fatalln("Error create circle log file: " + err.Error())
	}

	flag := log.LstdFlags | log.Llongfile

	Info = log.New(file, "INFO:  ", flag)
	Warning = log.New(file, "WARNING: ", flag)
	Error = log.New(file, "ERROR: ", flag)
	Debug = log.New(file, "DEBUG: ", flag)

	Info.Printf("Circle new start!")
}

func (l *LoggerConfig) InitializeLoggerScheduler(sch *gocron.Scheduler) {
	err := sch.Every(1).Day().At("00:00").Do(l.InitLogger)
	if err != nil {
		panic(fmt.Errorf("can't create scheduler of circle log file, err: %v", err))
	}
}
