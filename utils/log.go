package utils

import (
	"log"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

var (
	Logger   *logrus.Logger
	level    = "debug"
	rootPath string
)

func getRootPath() (string, int) {
	_, filename, line, _ := runtime.Caller(0)
	return path.Dir(path.Dir(filename) + "/../"), line
}

func init() {
	Logger = logrus.New()

	switch level {
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	case "info":
		Logger.SetLevel(logrus.InfoLevel)
	case "panic":
		Logger.SetLevel(logrus.PanicLevel)
	default:
		log.Fatalf("设置日志级别有误:%s", level)
	}

	Logger.Formatter = &logrus.TextFormatter{
		ForceColors:     true, // 控制台日志显示颜色
		DisableColors:   false,
		FullTimestamp:   true,                  // 显示完整的时间戳
		TimestampFormat: "2006-01-02 15:04:05", // 自定义时间戳格式
	}

	Logger.Out = os.Stdout
}
