package dbconn

import (
	"context"
	"errors"
	"time"

	"github.com/dwnGnL/testWork/lib/goerrors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	gormlogger "gorm.io/gorm/logger"
)

type logger struct {
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
}

func NewLogger() *logger {
	return &logger{
		SkipErrRecordNotFound: true,
	}
}

func (l *logger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *logger) Info(ctx context.Context, s string, args ...interface{}) {
	goerrors.Log().WithContext(ctx).Infof(s, args...)
}

func (l *logger) Warn(ctx context.Context, s string, args ...interface{}) {
	goerrors.Log().WithContext(ctx).Warnf(s, args...)
}

func (l *logger) Error(ctx context.Context, s string, args ...interface{}) {
	goerrors.Log().WithContext(ctx).Errorf(s, args...)
}

func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logrus.Fields{}
	// if l.SourceField != "" {
	// 	fields[l.SourceField] = utils.FileWithLineNum()
	// }
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		goerrors.Log().WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		goerrors.Log().WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}

	goerrors.Log().WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
}
