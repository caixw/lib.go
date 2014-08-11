// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package errors_test

import (
	"fmt"
	"testing"

	"github.com/caixw/lib.go/assert"
	"github.com/caixw/lib.go/errors"
)

func TestEqual(t *testing.T) {
	err1 := errors.New(5, nil, "abc")
	err2 := errors.New(5, nil, "abc")

	err3 := errors.New(6, err1, "abc")
	err4 := errors.New(6, err1, "abc")

	assert.False(t, err1 == err2)
	assert.False(t, err3 == err4)
}

func TestNew(t *testing.T) {
	err := errors.New(5, nil, "abc")

	a := assert.New(t)
	a.Equal(err.GetCode(), 5, "err.GetCode的值不等于5")
	a.Nil(err.GetPrevious())
	a.Equal(err.Error(), "abc")

	err2 := errors.New(5, err, "abc")
	a.Equal(err2.GetCode(), 5)
	a.Equal(err2.GetPrevious(), err)
	a.Equal(err2.Error(), "abc")
}

func ExampleNew() {
	err := errors.New(5, nil, "abc")
	if err != nil {
		fmt.Print(err.GetCode())
	}

	// Output: 5
}

func ExampleNewf() {
	err := errors.Newf(5, nil, "code=[%v]", 5)
	err2 := errors.New(6, err, "abc")
	if err2 != nil && err2.GetPrevious() != nil {
		fmt.Print(err2.GetPrevious())
	}

	// Output: code=[5]
}
