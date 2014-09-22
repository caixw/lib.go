// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mux

import (
	"net/http"
	"testing"

	"github.com/caixw/lib.go/assert"
)

func TestHost(t *testing.T) {
	defHandler := func(w http.ResponseWriter, r *http.Request) bool {
		return true
	}

	defFunc := MatcherFunc(defHandler)

	fn := func(host string, m Matcher, wont bool) {
		r, err := http.NewRequest("GET", "", nil)
		assert.NotError(t, err)

		r.Host = host
		assert.Equal(t, r.Host, host)
		assert.Equal(t, m.ServeHTTP2(nil, r), wont)
	}

	h := NewHost(defFunc, "www.example.com")
	fn("www.example.com", h, true)
	fn("www.abc.com", h, false)

	h = NewHost(defFunc, "\\w+.example.com")
	fn("www.example.com", h, true)
	fn("api.example.com", h, true)
	fn("www.abc.com", h, false)
}
