package log

import "github.com/alexcesaro/log/stdlog"
import "github.com/alexcesaro/log"

type GormLogger struct {
	logger log.Logger
}

func (l *GormLogger) Print(v ...interface{}) {
	l.logger.Info(v...)
}

func NewGormLogger(l log.Logger) *GormLogger {
	return &GormLogger{logger: l}
}

var Log log.Logger

func InitLog() {
	Log = stdlog.GetFromFlags()
}
