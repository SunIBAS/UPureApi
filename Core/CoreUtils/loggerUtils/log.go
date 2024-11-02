package loggerUtils

import (
	"log"
	"os"
)

type KeepLogger struct {
	infoLogger *log.Logger
	warnLogger *log.Logger
	errLogger  *log.Logger
}

func CreateLogger(filepath string) KeepLogger {
	kLogger := KeepLogger{}
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	kLogger.infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	kLogger.warnLogger = log.New(file, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	kLogger.errLogger = log.New(file, "ERR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return kLogger
}

func (kl KeepLogger) Info(str ...any) {
	kl.infoLogger.Println(str)
}

func (kl KeepLogger) Err(str ...any) {
	kl.errLogger.Println(str)
}

func (kl KeepLogger) Warn(str ...any) {
	kl.warnLogger.Println(str)
}
