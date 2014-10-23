// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/caixw/lib.go/logs/writer"
)

// log.Logger并未提供更换output的功能，为了达到config实例中，所有配置
// 节点都是一个io.Writer的功能，logWriter的作用就是将一个log.Logger伪
// 装io.Writer，但未真正实现Write()方法，因为最终在LevelLogger中使用的
// 还是log的实例。
type logWriter struct {
	level  int
	prefix string
	flag   int
	c      *writer.Container
	log    *log.Logger
}

var _ writer.FlushAdder = &logWriter{}
var _ io.Writer = &logWriter{}

func (l *logWriter) Write(bs []byte) (int, error) {
	panic("该函数并未真正实现，仅为支持接口而设")
	return 0, nil
}

// writer.Adder.Add()
func (l *logWriter) Add(w io.Writer) error {
	l.c.Add(w)
	return nil
}

func (l *logWriter) Flush() (int, error) {
	return l.c.Flush()
}

// toLogger将当前类转换成log.Logger实例。
func (l *logWriter) toLogger() *log.Logger {
	if l.log != nil {
		panic("log.Logger已经生成")
	}

	l.log = log.New(l.c, l.prefix, l.flag)
	return l.log
}

var flagMap = map[string]int{
	"log.ldate":         log.Ldate,
	"log.ltime":         log.Ltime,
	"log.lmicroseconds": log.Lmicroseconds,
	"log.llongfile":     log.Llongfile,
	"log.lshortfile":    log.Lshortfile,
	"log.lstdflags":     log.LstdFlags,
}

func logWriterInitializer(level int, args map[string]string) (io.Writer, error) {
	flagStr, found := args["flag"]
	if !found {
		flagStr = "log.lstdflags"
	}

	flag, found := flagMap[strings.ToLower(flagStr)]
	if !found {
		return nil, fmt.Errorf("未知的Flag参数:[%v]", flagStr)
	}

	prefix, found := args["prefix"]

	return &logWriter{
		level:  level,
		flag:   flag,
		prefix: prefix,
		c:      writer.NewContainer(),
	}, nil
}

func init() {
	reg := func(levelName string, level int) {
		fn := func(args map[string]string) (io.Writer, error) {
			return logWriterInitializer(level, args)
		}
		if !Register(levelName, fn) {
			panic(fmt.Sprintf("注册[%v]未成功", levelName))
		}
	}

	reg("info", LevelInfo)
	reg("debug", LevelDebug)
	reg("trace", LevelTrace)
	reg("warn", LevelWarn)
	reg("error", LevelError)
	reg("critical", LevelCritical)
}
