// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package term

const (
	SGRReset           = "0"
	SGRBold            = "1"
	SGRUnderline       = "4"
	SGRBlink           = "5" // 闪烁
	SGRReverseVideo    = "7" // 反显
	SGRConceal         = "8"
	SGRBoldOff         = "22"
	SGRUnderlineOff    = "24"
	SGRBlinkOff        = "25"
	SGRReverseVideoOff = "27"
	SGRConcealOff      = "28"

	SGRFBlack   = "30"
	SGRFRed     = "31"
	SGRFGreen   = "32"
	SGRFYellow  = "33"
	SGRFBlue    = "34"
	SGRFMagenta = "35"
	SGRFCyan    = "36"
	SGRFWhite   = "37"
	SGRFDefault = "39" // 默认前景色

	SGRBBlack   = "40"
	SGRBRed     = "41"
	SGRBGreen   = "42"
	SGRBYellow  = "43"
	SGRBBlue    = "44"
	SGRBMagenta = "45"
	SGRBCyan    = "46"
	SGRBWhite   = "47"
	SGRBDefault = "49" // 默认背景色
)

// 将几个SGR控制符合成一个ansi对象
func SGR(args ...string) ansi {
	if len(args) == 0 {
		return ansi(SGRReset + "m")
	}

	ret := ""
	for _, v := range args {
		ret += v + ";"
	}

	return ansi(ret[0:len(ret)-1] + "m")
}
