// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package assert

import (
	"errors"
	"testing"
)

func TestTrue(t *testing.T) {
	True(t, true)
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

	Equal(t, 5, v1)
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

func TestEmpty(t *testing.T) {
	Empty(t, 0, "Empty(0) falid")
	Empty(t, "", "Empty(``) falid")
	Empty(t, false, "Empty(false) falid")
	Empty(t, []string{}, "Empty(slice{}) falid")
	Empty(t, []int{}, "Empty(slice{}) falid")
}

func TestNotEmpty(t *testing.T) {
	NotEmpty(t, 1, "NotEmpty(1) falid")
	NotEmpty(t, true, "NotEmpty(true) falid")
	NotEmpty(t, []string{"ab"}, "NotEmpty(slice(abc)) falid")
}

func TestError(t *testing.T) {
	err := errors.New("test")
	Error(t, err, "Error(err) falid")

}

func TestNotError(t *testing.T) {
	NotError(t, "123", "NotError(123) falid")
}

func TestFileExists(t *testing.T) {
	FileExists(t, "./assert.go", "FileExists() falid")
}

func TestFileNotExists(t *testing.T) {
	FileNotExists(t, "c:/win", "FileNotExists() falid")
}
