package GormULog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/universe-30/ULog"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Config struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	LogLevel                  gormlogger.LogLevel
}

type gormLocalLogger struct {
	LocalLogger ULog.Logger
	Config

	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func New_gormLocalLogger(ulog ULog.Logger, config Config) *gormLocalLogger {
	l := &gormLocalLogger{
		LocalLogger:  ulog,
		infoStr:      "%s\n[info] ",
		warnStr:      "%s\n[warn] ",
		errStr:       "%s\n[error] ",
		traceStr:     "%s\n[%.3fms] [rows:%v] %s",
		traceWarnStr: "%s %s\n[%.3fms] [rows:%v] %s",
		traceErrStr:  "%s %s\n[%.3fms] [rows:%v] %s",
	}
	if config.SlowThreshold == 0 {
		config.SlowThreshold = 500 * time.Millisecond
	}
	if config.LogLevel == 0 {
		config.LogLevel = gormlogger.Warn
	}

	l.SlowThreshold = config.SlowThreshold
	l.IgnoreRecordNotFoundError = config.IgnoreRecordNotFoundError
	l.LogLevel = config.LogLevel

	return l
}

func (l *gormLocalLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l *gormLocalLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(err, gormlogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		var s string
		if rows == -1 {
			s = fmt.Sprintf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			//l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			s = fmt.Sprintf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			//l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
		l.LocalLogger.Errorln(s)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		var s string
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			s = fmt.Sprintf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			//l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			s = fmt.Sprintf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			//l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
		l.LocalLogger.Warnln(s)
	case l.LogLevel == gormlogger.Info:
		sql, rows := fc()
		var s string
		if rows == -1 {
			s = fmt.Sprintf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
			//l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			s = fmt.Sprintf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			//l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
		l.LocalLogger.Infoln(s)
	}
}

func (l *gormLocalLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		//l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)

		s := fmt.Sprintf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
		l.LocalLogger.Infoln(s)
	}
}

func (l *gormLocalLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		//l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)

		s := fmt.Sprintf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
		l.LocalLogger.Warnln(s)
	}
}

func (l *gormLocalLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		//l.Printf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)

		s := fmt.Sprintf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
		l.LocalLogger.Errorln(s)
	}
}
