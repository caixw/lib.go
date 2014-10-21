lib.go [![Build Status](https://travis-ci.org/caixw/lib.go.svg?branch=master)](https://travis-ci.org/caixw/lib.go) [![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/github.com/caixw/lib.go) [![GoDoc](https://godoc.org/github.com/caixw/lib.go/assert?status.svg)](https://godoc.org/github.com/caixw/lib.go) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/caixw/lib.go/blob/master/LICENSE)
=====

Go语言包的一个集合，包含各个方面，比较杂乱。

#### assert
常用的断言操作。对testing.T的一个简单封装。

#### term
提供了对ansi控制码的一些基本操作。不支持windows。

#### conv
各类型数据之间的相互转换。

#### errors
对系统errors的简单扩展。

#### session
web服务器的session管理包。

#### encoding/version
版本号的比较和解析。

#### encoding/ini
ini文件的解析。

#### encoding/tag
对固定格式的struct tag的分析。

#### mux
对net/http.ServeMux的简单扩展，可以实现大部分路由功能。

#### validation
验证工具。

#### validation/validator
一些常用的验证函数。

#### logs
一个默认的日志系统

#### logs/writer
供日志系统使用的几个io.Writer
