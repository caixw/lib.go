// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package ini

import (
	"bytes"
	"testing"

	"github.com/caixw/lib.go/assert"
)

func TestWriter(t *testing.T) {
	a := assert.New(t)
	buf := bytes.NewBufferString("")

	w := NewWriter(buf)
	a.NotNil(w)

	w.AddElement("key", "val")
	w.AddComment("comment")
	w.AddSection("section")
	w.AddElement("key", "val")
	w.AddElement("key1", "val1")
	w.Flush()
	data := buf.Bytes()

	str := `key=val
#comment
[section]
key=val
key1=val1
`

	a.Equal([]byte(str), data)
}
