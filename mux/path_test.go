// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mux

import (
	"net/http"
	"testing"

	"github.com/caixw/lib.go/assert"
)

func TestPath(t *testing.T) {
	defHandler := func(w http.ResponseWriter, r *http.Request) bool {
		return true
	}

	defFunc := MatcherFunc(defHandler)

	fn := func(pattern string, m Matcher, wont bool) {
		r, err := http.NewRequest("GET", pattern, nil)
		assert.NotError(t, err)
		assert.Equal(t, m.ServeHTTP2(nil, r), wont)
	}

	p := NewPath(defFunc, "/api")
	fn("/api", p, true)
	fn("/api/v1", p, true)

	p = NewPath(defFunc, "/api/v(\\d+)")
	fn("/api", p, false)
	fn("/api/v1", p, true)
	fn("/api/v1/post/1", p, true)
}
