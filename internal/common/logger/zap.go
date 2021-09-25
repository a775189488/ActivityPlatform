package logger

import (
	"fmt"
	"os"
	"time"

	"entrytask/internal/conf"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	Log *zap.SugaredLogger
}

func (l *Logger) Init() {
	now := time.Now()

	var syncStyle zapcore.WriteSyncer
	if conf.Config.App.LogPath == "" {
		syncStyle = zapcore.AddSync(os.Stdout)
	} else {
		syncStyle = zapcore.AddSync(&lumberjack.Logger{
			Filename:   fmt.Sprintf("%s/%04d%02d%02d%s", conf.Config.App.LogPath, now.Year(), now.Month(), now.Day(), ".log"),
			MaxSize:    2,
			MaxBackups: 100,
			MaxAge:     30, //days
			Compress:   false,
		})
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		syncStyle,
		zap.NewAtomicLevel(),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	l.Log = logger
}

//Info 打印信息
func (l *Logger) Info(args ...interface{}) {
	l.Log.Info(args...)
}

//Infof 打印信息，附带template信息
func (l *Logger) Infof(template string, args ...interface{}) {
	l.Log.Infof(template, args...)
}

//Warn 打印警告
func (l *Logger) Warn(args ...interface{}) {
	l.Log.Warn(args...)
}

//Warnf 打印警告，附带template信息
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.Log.Warnf(template, args...)
}

//Error 打印错误
func (l *Logger) Error(args ...interface{}) {
	l.Log.Error(args...)
}

//Errorf 打印错误，附带template信息
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.Log.Errorf(template, args...)
}

//Panic 打印Panic信息
func (l *Logger) Panic(args ...interface{}) {
	l.Log.Panic(args...)
}

//Panicf 打印Panic信息，附带template信息
func (l *Logger) Panicf(template string, args ...interface{}) {
	l.Log.Panicf(template, args...)
}

//DPanic 打印Panic信息
func (l *Logger) DPanic(args ...interface{}) {
	l.Log.DPanic(args...)
}

//DPanicf 打印DPanic信息，附带template信息
func (l *Logger) DPanicf(template string, args ...interface{}) {
	l.Log.DPanicf(template, args...)
}
