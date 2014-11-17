// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"reflect"
	"testing"

	"github.com/caixw/lib.go/assert"
)

type FetchObj struct {
	Id       int    `orm:"name:id;ai:1,2;"`
	Email    string `orm:"unique:unique_index;nullable;pk:pk_name"`
	Username string `orm:"index:index"`
	Group    int    `orm:"name:group;fk:fk_group,group,id"`

	Regdate int `orm:"-"`
}

func TestParseObj(t *testing.T) {
	a := assert.New(t)
	obj := &FetchObj{Id: 5}
	mapped := map[string]reflect.Value{}

	v := reflect.ValueOf(obj).Elem()
	a.True(v.IsValid())

	parseObj(v, &mapped)
	a.Equal(4, len(mapped))

	// 忽略的字段
	_, found := mapped["Regdate"]
	a.False(found)

	// 判断字段是否存在
	vi, found := mapped["id"]
	a.True(found).True(vi.IsValid())

	// 设置字段的值
	mapped["id"].Set(reflect.ValueOf(36))
	a.Equal(36, obj.Id)
	mapped["Email"].SetString("email")
	a.Equal("email", obj.Email)
	mapped["Username"].SetString("username")
	a.Equal("username", obj.Username)
	mapped["group"].SetInt(1)
	a.Equal(1, obj.Group)
}
