package log

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func Log(level logrus.Level, logMessage interface{}) error {
	// validate data
	if level == 0 {
		return errors.New("uncertain log level")
	}

	if logMessage == nil {
		return errors.New("uncertain log message")
	}

	// log writter
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
		panic(err)
	}
	defer f.Close()

	// set logger config
	log := &logrus.Logger{
		Out:   io.MultiWriter(f, os.Stdout),
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%\n",
		},
	}

	// log proccess
	log.Log(level, logMessage)

	return nil
}
