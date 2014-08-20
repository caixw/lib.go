lib.go [![Build Status](https://travis-ci.org/caixw/lib.go.svg?branch=master)](https://travis-ci.org/caixw/lib.go) [![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/github.com/caixw/lib.go) [![GoDoc](https://godoc.org/github.com/caixw/lib.go/assert?status.svg)](https://godoc.org/github.com/caixw/lib.go) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/caixw/lib.go/blob/master/LICENSE)
=====

Go语言包的一个集合，包含各个方面，比较杂乱。

每个包都提供了一个版本号(Version)表示包的兼容情况：若主版本发
生变化，则表示不兼容；其它子版本号变化，表示一些小修改，但不会涉及到兼容性问题。

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
