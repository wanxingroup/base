package log

import (
	"fmt"
	"path"
	"regexp"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

func RegisterFilePath(logger *logrus.Logger) {

	// Set log formatter to output source code file name, line number and function name.
	var re = regexp.MustCompile(`^dev-gitlab.wanxingrowth.com/`)
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339Nano,
		FullTimestamp:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			fileName := path.Base(f.File)
			return fmt.Sprintf("%s()", re.ReplaceAllString(f.Function, "")), fmt.Sprintf("%s:%d", fileName, f.Line)
		},
	})
}
