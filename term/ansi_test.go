// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package term

import (
	"testing"
)

func TestColor(t *testing.T) {
	t.Logf("%m", Erase(2))

	t.Logf("%mFRed", FRed)
	t.Logf("%mFGreen", FGreen)
	t.Logf("%mFYellow", FYellow)
	t.Logf("%mFBlue", FBlue)
	t.Logf("%mFMagenta", FMagenta)
	t.Logf("%mFCyan", FCyan)
	t.Logf("%mFWhite", FWhite)
	t.Logf("%mFDefault", FDefault)

	t.Logf("%mBRed", BRed)
	t.Logf("%mBGreen", BGreen)
	t.Logf("%mBYellow", BYellow)
	t.Logf("%mBBlue", BBlue)
	t.Logf("%mBMagenta", BMagenta)
	t.Logf("%mBCyan", BCyan)
	t.Logf("%mBWhite", BWhite)
	t.Logf("%mBDefault", BDefault)

	for i := 0; i < 256; i++ {
		t.Logf("%m字体颜色%d", FColor256(i), i)
	}

	for i := 0; i < 256; i++ {
		t.Logf("%m背景颜色%d", BColor256(i), i)
	}

	t.Logf("%m", Reset)
}
