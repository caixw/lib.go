// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package validation

import (
	"testing"

	"github.com/caixw/lib.go/assert"
)

func TestValidation(t *testing.T) {
	v := New()
	a := assert.New(t)

	v.Apply(true, "true", "").
		Apply(false, "msg1", "key1").
		Apply(false, "msg2", "key2").
		Apply(false, "msg3", "key2").
		Apply(false, "msg4", "")

	a.True(v.HasErrors())
	a.Equal(4, len(v.GetErrors()))
	a.Equal(2, len(v.GetErrorsMap()))
	a.Equal(v.GetErrors(), []string{"msg1", "msg2", "msg3", "msg4"})
	a.Equal(v.GetErrorsMap(), map[string]string{"key1": "msg1", "key2": "msg3"})

	v.Clear()
	a.False(v.HasErrors())
}

func TestResult(t *testing.T) {
	a := assert.New(t)
	v := New()

	ok := func() *Result {
		return &Result{Ok: true}
	}

	fail := func() *Result {
		return &Result{Ok: false, Message: "Message", Key: "Key"}
	}

	v.ApplyResult(ok()).ApplyResult(fail())

	a.True(v.HasErrors())
	a.Equal(v.GetErrors(), []string{"Message"})
	a.Equal(v.GetErrorsMap(), map[string]string{"Key": "Message"})
	v.Clear()
	a.False(v.HasErrors())
}
