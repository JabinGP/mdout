package config

import (
	"os"

	"github.com/JabinGP/mdout/static"
	"github.com/JabinGP/mdout/tool"
	log "github.com/sirupsen/logrus"
)

type CombLogger struct {
	// Stdout 标准输出实例
	Stdout *log.Logger
	// File 日志文件输出实例
	File *log.Logger
}

func initLogger() {
	stdout := newStdout()
	file := newFile()

	PublicLogger = CombLogger{
		Stdout: stdout,
		File:   file,
	}
}

func initLogFile() {
	if !tool.IsExists(static.LogFolderFullName) {
		err := os.Mkdir(static.LogFolderFullName, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	if !tool.IsExists(static.LogFileFullName) {
		_, err := os.Create(static.LogFileFullName)
		if err != nil {
			panic(err)
		}
	}
}

func newStdout() *log.Logger {
	level, err := tool.TransformToLogrusLevel(Obj.Runtime.StdoutLogLevel)
	if err != nil {
		panic(err)
	}
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
	logger.SetLevel(log.Level(level))
	return logger
}

func newFile() *log.Logger {
	initLogFile()
	level, err := tool.TransformToLogrusLevel(Obj.Runtime.FileLogLevel)
	if err != nil {
		panic(err)
	}
	logger := log.New()
	file, err := os.OpenFile(static.LogFileFullName, os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(file)
	logger.SetLevel(log.Level(level))
	return logger
}

// SetStdoutLevel ...
func (c *CombLogger) SetStdoutLevel(level uint32) {
	c.Stdout.SetLevel(log.Level(level))
}

// Debugf ...
func (c *CombLogger) Debugf(format string, args ...interface{}) {
	c.Stdout.Debugf(format, args...)
	c.File.Debugf(format, args...)
}

// Infof ...
func (c *CombLogger) Infof(format string, args ...interface{}) {
	c.Stdout.Infof(format, args...)
	c.File.Infof(format, args...)
}

// Errorf ..
func (c *CombLogger) Errorf(format string, args ...interface{}) {
	c.Stdout.Errorf(format, args...)
	c.File.Errorf(format, args...)
}

// Debugln ..
func (c *CombLogger) Debugln(args ...interface{}) {
	c.Stdout.Debugln(args...)
	c.File.Debugln(args...)
}

// Infoln ...
func (c *CombLogger) Infoln(args ...interface{}) {
	c.Stdout.Infoln(args...)
	c.File.Infoln(args...)
}

// Errorln ...
func (c *CombLogger) Errorln(args ...interface{}) {
	c.Stdout.Errorln(args...)
	c.File.Errorln(args...)
}
