package utils

import (
	"bytes"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	Logger *logrus.Logger
)

type CustomFormatter struct{}

// Format 是CustomFormatter必须实现的方法，用于格式化日志条目
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := entry.Level.String()
	msg := entry.Message
	fileName := entry.Caller.File
	lineNumber := entry.Caller.Line
	fmt.Fprintf(b, "时间:%s----- 级别:%s----- 内容:%s -----位置:%s:%d\n", timestamp, level, msg, fileName, lineNumber)

	return b.Bytes(), nil
}

func init() {
	Logger = logrus.New()

	Logger.SetLevel(logrus.DebugLevel)
	formatter := new(CustomFormatter)
	Logger.SetFormatter(formatter)
	Logger.SetReportCaller(true)
	Logger.Out = os.Stdout
}
