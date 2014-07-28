// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 扩展fmt.Printf，实现对ansi控制码的输出。
//  // 普通的fmt格式化字符串，使用%m占位符
//  fmt.Printf("%m这是红色的字", term.FRed)
//  fmt.Printf("%m这是红色字，绿色背景", term.SGR(term.SGRFRed,term.SGRBGreen))
//  fmt.Printf("%m%m这是红色字，绿色背景", term.FRed,term.BGreen)
//
//  // 申请一个Writer，可以使用链式写法。
//  w := term.NewWriter(writer)
//  w.Left(5).SGR(term.SGRFRed).Printf("%s", "string").Move(1,1)
//
// ansi的相关文档，可参考以下内容：
//  http://en.wikipedia.org/wiki/ANSI_escape_code
//  http://www.mudpedia.org/mediawiki/index.php/ANSI_colors
package term

// 当前库的版本
const Version = "0.1.0.140727"
