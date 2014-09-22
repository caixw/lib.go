// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package validation

import (
	"github.com/caixw/lib.go/conv"
)

func MaxInt(max int64, val interface{}) bool {
	if v, err := conv.Int64(val); err != nil {
		return false
	} else {
		return v <= max
	}
}

func MaxFloat(max float64, val interface{}) bool {
	if v, err := conv.Float64(val); err != nil {
		return false
	} else {
		return v <= max
	}
}
