package log

import (
	"os"
	"sync"

	"github.com/JabinGP/mdout/static"
	"github.com/JabinGP/mdout/tool"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

var once sync.Once

// Stdout 控制台输出日志对象
var Stdout *logrus.Logger

// File 文件输出日志对象
var File *logrus.Logger

// 默认的日志输出等级
var stdoutDefaultLevel = logrus.DebugLevel
var fileDefaultLevel = logrus.DebugLevel

func init() {
	once.Do(func() {
		initLogFile()
		Stdout = newStdout()
		File = newFile()
	})
}

func initLogFile() {
	if err := tool.InitFolder(static.ConfigFolderFullName); err != nil {
		panic(err)
	}
	if err := tool.InitFolder(static.LogFolderFullName); err != nil {
		panic(err)
	}

	if !tool.IsExists(static.LogFileFullName) {
		_, err := os.Create(static.LogFileFullName)
		if err != nil {
			panic(err)
		}
	}
}

func newStdout() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	logrus.SetOutput(colorable.NewColorableStdout())
	logger.SetLevel(stdoutDefaultLevel)

	return logger
}

func newFile() *logrus.Logger {
	logger := logrus.New()

	file, err := os.OpenFile(static.LogFileFullName, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(file)
	logger.SetLevel(fileDefaultLevel)
	return logger
}

// SetStdoutLevel ...
func SetStdoutLevel(level uint32) {
	Stdout.SetLevel(logrus.Level(level))
}

// SetFileLevel ...
func SetFileLevel(level uint32) {
	File.SetLevel(logrus.Level(level))
}

// Debugf ...
func Debugf(format string, args ...interface{}) {
	Stdout.Debugf(format, args...)
	File.Debugf(format, args...)
}

// Infof ...
func Infof(format string, args ...interface{}) {
	Stdout.Infof(format, args...)
	File.Infof(format, args...)
}

// Errorf ..
func Errorf(format string, args ...interface{}) {
	Stdout.Errorf(format, args...)
	File.Errorf(format, args...)
}

// Debugln ..
func Debugln(args ...interface{}) {
	Stdout.Debugln(args...)
	File.Debugln(args...)
}

// Infoln ...
func Infoln(args ...interface{}) {
	Stdout.Infoln(args...)
	File.Infoln(args...)
}

// Errorln ...
func Errorln(args ...interface{}) {
	Stdout.Errorln(args...)
	File.Errorln(args...)
}
