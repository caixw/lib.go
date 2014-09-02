// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mux

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caixw/lib.go/assert"
)

func TestServe(t *testing.T) {
	a := assert.New(t)

	api := func(w http.ResponseWriter, r *http.Request) {
		n, err := w.Write([]byte("api"))
		a.NotError(err)
		a.NotEqual(n, 0)
	}

	Get("/api/", http.HandlerFunc(api))
	ts := httptest.NewServer(defaultServeMux)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/api/")
	a.NotError(err)
	a.NotNil(res)
	txt, err := ioutil.ReadAll(res.Body)
	a.NotError(err)
	a.NotEmpty(txt)
	t.Log(string(txt))
}
