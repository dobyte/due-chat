package mysql

import (
	"context"
	"fmt"
	"github.com/dobyte/due/v2/errors"
	"github.com/dobyte/due/v2/log"
	loggers "gorm.io/gorm/logger"
	"time"
)

const traceStr = "%s\n[%.3fms] [rows:%v] %s"

type logger struct {
	logLevel                  loggers.LogLevel
	ignoreRecordNotFoundError bool
	slowThreshold             time.Duration
}

var _ loggers.Interface = &logger{}

func (l *logger) LogMode(level loggers.LogLevel) loggers.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

func (l *logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= loggers.Info {
		log.Printf(log.InfoLevel, msg, data)
	}
}

func (l *logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= loggers.Warn {
		log.Printf(log.WarnLevel, msg, data)
	}
}

func (l *logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= loggers.Error {
		log.Printf(log.ErrorLevel, msg, data)
	}
}

// Trace print sql message
func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= loggers.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.logLevel >= loggers.Error && (!errors.Is(err, loggers.ErrRecordNotFound) || !l.ignoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			log.Printf(log.ErrorLevel, traceStr, err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Printf(log.ErrorLevel, traceStr, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= loggers.Error:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.slowThreshold)
		if rows == -1 {
			log.Printf(log.WarnLevel, traceStr, slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Printf(log.WarnLevel, traceStr, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.logLevel == loggers.Info:
		sql, rows := fc()
		if rows == -1 {
			log.Printf(log.InfoLevel, traceStr, "", float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Printf(log.InfoLevel, traceStr, "", float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
