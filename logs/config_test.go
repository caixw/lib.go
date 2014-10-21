// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"bytes"
	"testing"

	"github.com/caixw/lib.go/assert"
)

var xmlCfg = `
<?xml version="1.0" encoding="utf-8" ?>
<logs>
    <info>
        <buffer size="5">
            <file dir="/var/logs/info" />
        </buffer>
        <console color="\033[12m" />
    </info>

    <debug>
        <file dir="/var/logs/debug" />
        <console color="\033[12" />
    </debug>
</logs>
`

func TestLoadFromXml(t *testing.T) {
	a := assert.New(t)

	r := bytes.NewReader([]byte(xmlCfg))
	a.NotNil(r)

	cfg, err := loadFromXml(r)
	a.NotError(err).NotNil(cfg)
	a.Equal(2, len(cfg.Items)) // info debug

	info, found := cfg.Items["info"]
	a.True(found).NotNil(info).Equal(info.Name, "info")
	a.Equal(2, len(info.Items)) // buffer,console

	buf, found := info.Items["buffer"]
	a.True(found).NotNil(buf)
	a.Equal(buf.Attrs["size"], "5")
}
