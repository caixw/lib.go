// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package validation

import (
	"strings"
	"sync"
)

type errors struct {
	items []string
}

func (e *errors) add(msg string) {
	e.items = append(e.items, msg)
}

func (e *errors) GetErrors() []string {
	return e.items
}

// fmt.Stringer
func (e *errors) String() string {
	l := len(e.items)
	switch l {
	case 0:
		return "<nil>"
	case 1:
		return e.items[0]
	default:
		return strings.Join(e.items, ";")
	}
}

func (e *errors) free() {
	e.items = e.items[:0]
	errFree.Put(e)
}

var errFree = sync.Pool{
	New: func() interface{} { return new(errors) },
}

func newErrors(msg string) *errors {
	err := errFree.Get().(*errors)
	err.add(msg)

	return err
}
