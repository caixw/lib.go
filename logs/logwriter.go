// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"io"
	"log"

	"github.com/caixw/lib.go/logs/writer"
)

type logWriter struct {
	level  int
	prefix string
	flag   int
	c      *writer.Container
}

var _ writer.WriterContainer = &logWriter{}

func (l *logWriter) Write(bs []byte) (int, error) {
	return 0, nil
}

func (l *logWriter) AddWriter(w io.Writer) error {
	l.c.AddWriter(w)
	return nil
}

func (l *logWriter) toLogger() *log.Logger {
	return log.New(l.c, l.prefix, l.flag)
}

func logWriterInitializer(level, args map[string]string) (io.Writer, error) {
	//
	return nil, nil
}
