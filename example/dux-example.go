package main

import (
	"fmt"
	"os"

	"github.com/daichitakahashi/zapmux"
	"github.com/hnakamur/zap-ltsv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	accessLogFile, _ := os.OpenFile("access.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	errorLogFile, _ := os.OpenFile("error.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

	enc := ltsv.NewLTSVEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	core := zapmux.NewDuxCore(
		zapcore.NewCore(enc, accessLogFile, zapcore.InfoLevel),
		zapcore.NewCore(enc, errorLogFile, zapcore.ErrorLevel),
	).WithInterceptor(func(main, sub *zapmux.CoreEntry) {
		if main.Entry.Level >= zapcore.ErrorLevel {
			main.With(zap.String("msg", "error occurred: see error.log"))
			sub.With(zap.String("msg", sub.Entry.Message), zap.Stack("stacktrace"))
			fmt.Println("over error")
		} else {
			main.With(zap.String("msg", main.Entry.Message))
			sub.Core = nil
			fmt.Println("under error")
		}
	})

	logger := zap.New(core).With(zap.String("ip", ""), zap.String("status", ""))

	logger.Info("test without any fields")
	logger.Error("test error without any fields")
	logger.Info("test without any fields2")
}
