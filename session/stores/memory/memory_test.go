// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package memory

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caixw/lib.go/assert"
	sess "github.com/caixw/lib.go/session"
)

func hello(a *assert.Assertion, w http.ResponseWriter) {
	n, err := w.Write([]byte("hello world"))
	a.NotError(err)
	a.NotEmpty(n)
}

func freeSession(a *assert.Assertion, inst *sess.Instance, w http.ResponseWriter, r *http.Request) {
	err := inst.FreeSession(w, r)
	a.NotError(err)
}

// 测试freeSession之后，值是否被清空
func freeSession2(a *assert.Assertion, inst *sess.Instance, w http.ResponseWriter, r *http.Request) {
	sess, err := inst.StartSession(w, r)
	a.NotError(err)
	a.NotNil(sess)

	uid, found := sess.Get("uid")
	a.False(found)
	a.Equal(uid, nil)
}

func testSession(a *assert.Assertion, inst *sess.Instance, w http.ResponseWriter, r *http.Request) {
	// 获取session接口实例
	sess, err := inst.StartSession(w, r)
	a.NotError(err)
	a.NotNil(sess)

	// 设置值
	sess.Set("uid", 1)
	uid, found := sess.Get("uid")
	a.True(found)
	a.Equal(uid, 1)

	// 不存在的值
	username := sess.MustGet("username", "abc")
	a.Equal(username, "abc")

	// 添加值
	sess.Set("username", "caixw")
	username = sess.MustGet("username", "abc")
	a.Equal(username, "caixw")

	// 修改值
	sess.Set("uid", "2")
	uid, found = sess.Get("uid")
	a.True(found)
	a.Equal(uid, "2")
}

// 测试cookie能否正常使用
func testSession2(a *assert.Assertion, inst *sess.Instance, w http.ResponseWriter, r *http.Request) {
	sess, err := inst.StartSession(w, r)
	a.NotError(err)
	a.NotNil(sess)

	uid, found := sess.Get("uid")
	a.True(found)
	a.Equal(uid, "2")
}

func TestMemory(t *testing.T) {
	a := assert.New(t)
	inst := sess.New(New(), "gosid", 3600, false)
	a.NotNil(inst)

	// 构建服务端
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		a.NotError(err)
		if action, found := r.Form["action"]; found {
			switch action[0] {
			case "testSession":
				testSession(a, inst, w, r)
			case "freeSession2":
				freeSession2(a, inst, w, r)
			case "freeSession":
				freeSession(a, inst, w, r)
			case "testSession2":
				testSession2(a, inst, w, r)
			}
			return
		}

		// 没有action参数的情况下，输出普通内容。
		hello(a, w)
	}))
	defer ts.Close()

	// 普通客户端请求
	response, err := http.Get(ts.URL)
	a.NotError(err)
	a.NotNil(response)
	txt, err := ioutil.ReadAll(response.Body)
	a.NotError(err)
	a.NotEmpty(txt)
	a.T().Log(string(txt))

	// testSession
	response, err = http.Get(ts.URL + "?action=testSession")
	a.NotError(err)
	a.NotNil(response)

	// 提交cookie,测试是否可以还原数据
	req, err := http.NewRequest("GET", ts.URL+"?action=testSession2", nil)
	a.NotError(err)
	a.NotNil(req)
	req.Header.Add("Cookie", response.Header.Get("Set-Cookie")) // 模拟浏览器提交cookie
	client := &http.Client{}
	response, err = client.Do(req)
	a.NotError(err)
	a.NotNil(response)

	// freeSession
	req, err = http.NewRequest("GET", ts.URL+"?action=freeSession", nil)
	a.NotError(err)
	a.NotNil(req)
	req.Header.Add("Cookie", response.Header.Get("Set-Cookie")) // 模拟浏览器提交cookie
	client = &http.Client{}
	response, err = client.Do(req)
	a.NotError(err)
	a.NotNil(response)

	// freeSession2 测试FreeSession之后，是否正常
	response, err = http.Get(ts.URL + "?action=freeSession2")
	a.NotError(err)
	a.NotNil(response)
}
