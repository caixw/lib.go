// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"errors"
	"io"
	"strings"

	"github.com/caixw/lib.go/logs/writer"
)

// writer的初始化函数。
type WriterInitializer func(map[string]string) (io.Writer, error)

// 注册的writer，所有注册的writer，都可以通过配置文件配置。
var regInitializer = map[string]WriterInitializer{}
var regNames []string

// 注册一个initizlizer
// 返回值反映是否注册成功。若已经存在相同名称的，则返回false
func Register(name string, init WriterInitializer) bool {
	if IsRegisted(name) {
		return false
	}

	regInitializer[name] = init
	regNames = append(regNames, name)
	return true
}

// 查询指定名称的Writer是否已经被注册
func IsRegisted(name string) bool {
	_, found := regInitializer[name]
	return found
}

// 返回所有已注册的writer名称
func Registed() []string {
	return regNames
}

func stmpInitializer(args map[string]string) (io.Writer, error) {
	username, found := args["username"]
	if !found {
		return nil, errors.New("")
	}

	password, found := args["password"]
	if !found {
		return nil, errors.New("")
	}

	subject, found := args["subject"]
	if !found {
		return nil, errors.New("")
	}

	host, found := args["host"]
	if !found {
		return nil, errors.New("")
	}

	sendToStr, found := args["sendTo"]
	if !found {
		return nil, errors.New("")
	}

	sendTo := strings.Split(sendToStr, ";")

	return writer.NewSmtp(username, password, subject, host, sendTo), nil
}

func init() {
	if !Register("stmp", stmpInitializer) {
		panic("注册stmp时失败")
	}
}
