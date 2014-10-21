// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"fmt"
	"io"
	"log"

	"github.com/caixw/lib.go/logs/writer"
)

const (
	LevelInfo = iota
	LevelDebug
	LevelTrace
	LevelWarn
	LevelError
	LevelCritical
)

type LevelLogger struct {
	logs map[int]*log.Logger
}

var _ writer.WriterContainer = &LevelLogger{}

// 仅为实现接口，不作任何输出
func (l *LevelLogger) Write(bs []byte) (int, error) {
	return 0, nil
}

func (l *LevelLogger) AddWriter(w io.Writer) error {
	log, ok := w.(*logWriter)
	if !ok {
		return fmt.Errorf("必须为logWriter接口")
	}

	l.logs[log.level] = log.toLogger()
	return nil
}

func (l *LevelLogger) ToStdLogger(level int) (log *log.Logger, ok bool) {
	log, ok = l.logs[level]
	return
}

func (l *LevelLogger) Println(level int, v ...interface{}) {
	if log, found := l.logs[level]; found {
		log.Println(v...)
	}
}

func (l *LevelLogger) Printf(level int, format string, v ...interface{}) {
	if log, found := l.logs[level]; found {
		log.Printf(format, v...)
	}
}

// Info相当于LevelLogger.Println(v...)的简写方式
func (l *LevelLogger) Info(v ...interface{}) {
	l.Println(LevelInfo, v...)
}

// Infof根目录于LevelLogger.Printf(format, v...)的简写方式
func (l *LevelLogger) Infof(format string, v ...interface{}) {
	l.Printf(LevelInfo, format, v...)
}

// Debug相当于LevelLogger.Println(v...)的简写方式
func (l *LevelLogger) Debug(v ...interface{}) {
	l.Println(LevelDebug, v...)
}

// Debugf根目录于LevelLogger.Printf(format, v...)的简写方式
func (l *LevelLogger) Debugf(format string, v ...interface{}) {
	l.Printf(LevelDebug, format, v...)
}

// Trace相当于LevelLogger.Println(v...)的简写方式
func (l *LevelLogger) Trace(v ...interface{}) {
	l.Println(LevelTrace, v...)
}

// Tracef根目录于LevelLogger.Printf(format, v...)的简写方式
func (l *LevelLogger) Tracef(format string, v ...interface{}) {
	l.Printf(LevelTrace, format, v...)
}

// Warn相当于LevelLogger.Println(v...)的简写方式
func (l *LevelLogger) Warn(v ...interface{}) {
	l.Println(LevelWarn, v...)
}

// Warnf根目录于LevelLogger.Printf(format, v...)的简写方式
func (l *LevelLogger) Warnf(format string, v ...interface{}) {
	l.Printf(LevelWarn, format, v...)
}

// Error相当于LevelLogger.Println(v...)的简写方式
func (l *LevelLogger) Error(v ...interface{}) {
	l.Println(LevelError, v...)
}

// Errorf根目录于LevelLogger.Printf(format, v...)的简写方式
func (l *LevelLogger) Errorf(format string, v ...interface{}) {
	l.Printf(LevelError, format, v...)
}

// Critical相当于LevelLogger.Println(v...)的简写方式
func (l *LevelLogger) Critical(v ...interface{}) {
	l.Println(LevelCritical, v...)
}

// Criticalf根目录于LevelLogger.Printf(format, v...)的简写方式
func (l *LevelLogger) Criticalf(format string, v ...interface{}) {
	l.Printf(LevelCritical, format, v...)
}
