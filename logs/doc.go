// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 日志处理包。
//
// 为什么log.Logger不能更换writer为什么log.Logger不能指定flag为none
//
// logs.Info()
//
// logs.SetWriter(Info, os.Stderr)
// logs.SetWriter(Info, writer.NewStmp())
package logs

const Version = "0.1.2.141011"
