package core

import (
	"flag"
	"fmt"
	"github.com/natefinch/lumberjack"
	"os"
)

func init() {
	var confPath string
	flag.StringVar(&confPath, "c", "./service.yml", "配置文件路径")
	flag.Parse()
	loadConfig(&Conf, confPath)
	logPath := fmt.Sprintf("%v%v", Conf.Logger.Path, Conf.Logger.Filename)
	var lumberJackLogger = lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
		LocalTime:  true,
	}
	os.MkdirAll(Conf.Logger.Path, os.ModePerm)
	initLogger(lumberJackLogger, Conf.Logger.Level)
}
