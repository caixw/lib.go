// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package term

import (
	"testing"
)

func TestColor(t *testing.T) {
	t.Logf("%v%vFRed, BDefault%v\n", FRed, BDefault, Reset)
	t.Logf("%v%vFGreen, BWhite%v\n", FGreen, BWhite, Reset)
	t.Logf("%v%vFYellow, BCyan%v\n", FYellow, BCyan, Reset)
	t.Logf("%v%vFBlue, BMagenta%v\n", FBlue, BMagenta, Reset)
	t.Logf("%v%vFMagenta, BBlue%v\n", FMagenta, BBlue, Reset)
	t.Logf("%v%vFCyan, BYellow%v\n", FCyan, BYellow, Reset)
	t.Logf("%v%vFWhite, BGreen%v\n", FWhite, BGreen, Reset)
	t.Logf("%v%vFDefault, BRed%v\n", FDefault, BRed, Reset)

	for i := 0; i < 256; i += 10 {
		t.Logf("%v%v字体颜色%d, 背景颜色%d%v\n", FColor256(i), BColor256(255-i), i, 255-i, Reset)
	}
}
