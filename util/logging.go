package util

import (
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
)

const logPath = "."

func SetupLogging() {
	log.SetOutput(&lumberjack.Logger{
		Filename: "covidiometro.log",
		MaxSize:  10,
		MaxAge:   60,
		Compress: true,
	})
}

func ErrFatal(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
