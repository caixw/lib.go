// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package assert

import (
	"testing"
)

func TestAssert(t *testing.T) {
	Assert(1, t, true, "Assert(true) falid")
	Assert(1, t, !false, "Assert(!false) falid")
	Assert(1, t, 5 == 5, "Assert(5==5) falid")
}

func TestTrue(t *testing.T) {
	True(t, true, "True falid")
	True(t, 1 == 1, "True(1==1) falid")
}

func TestFalse(t *testing.T) {
	False(t, false, "False falid")
	False(t, 1 == 2, "False(1==2) falid")
}

func TestNil(t *testing.T) {
	Nil(t, nil, "Nil falid")

	var v interface{}
	Nil(t, v, "Nil(v) falid")
}

func TestNotNil(t *testing.T) {
	NotNil(t, 5, "NotNil falid")

	var v interface{} = 5
	NotNil(t, v, "NotNil falid")
}

func TestEqual(t *testing.T) {
	Equal(t, 5, 5, "Equal(5,5) falid")

	var v1, v2 interface{}
	v1 = 5
	v2 = 5

	Equal(t, 5, v1, "Equal(5,v1) falid")
	Equal(t, v1, v2, "Equal(v1,v2) falid")
}

func TestNotEqual(t *testing.T) {
	NotEqual(t, 5, 6, "NotEqual(5,6) falid")

	var v1, v2 interface{}
	v1 = 5
	v2 = 6

	NotEqual(t, 5, v2, "NotEqual(5,v2) falid")
	NotEqual(t, v1, v2, "NotEqual(v1,v2) falid")
}
