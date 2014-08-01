// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package term

import (
	"os"
	"testing"
)

func TestWriter(t *testing.T) {

	w := NewWriter(os.Stdout)

	w.Erase(2)

	for i := 0; i < 256; i++ {
		w.Color256(i, 255-i)
		w.Printf("FColor(%d),BColor(%d)", i, 255-i)
		w.WriteAnsi(Reset)
		w.Println()
	}

	w.WriteAnsi(Reset) //.Move(50, 100)
}