// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"encoding/xml"
	"fmt"
	"github.com/caixw/lib.go/logs/writer"
	"io"
	"os"
)

// 用于表示config.xml中的配置数据。
type Config struct {
	Parent *Config
	Name   string             // writer的名称，一般为节点名
	Attrs  map[string]string  // 参数列表
	Items  map[string]*Config // 若是容器，则还有子项
}

// 从xml 文件初始化Config实例。
func loadFromXmlFile(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return loadFromXml(f)
}

// 从一个xml reader初始化Config
func loadFromXml(r io.Reader) (*Config, error) {
	var cfg *Config = nil //&config{Parent: nil}
	var t xml.Token
	var err error

	d := xml.NewDecoder(r)
	for t, err = d.Token(); err == nil; t, err = d.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			c := &Config{
				Parent: cfg,
				Name:   token.Name.Local,
				Attrs:  make(map[string]string),
			}
			for _, v := range token.Attr {
				c.Attrs[v.Name.Local] = v.Value
			}

			if cfg != nil {
				if cfg.Items == nil {
					cfg.Items = make(map[string]*Config)
				}
				cfg.Items[token.Name.Local] = c
			}
			cfg = c
		case xml.EndElement:
			if cfg.Parent != nil {
				cfg = cfg.Parent
			}
		default: // 可能还有ProcInst,CharData,Comment等用不到的标签
			continue
		}
	} // end for

	if err != io.EOF {
		return nil, err
	}

	return cfg, nil
}

// 将当前的config转换成io.Writer
func (c *Config) toWriter() (io.Writer, error) {
	initializer, found := regInitializer[c.Name]
	if !found {
		return nil, fmt.Errorf("未注册的初始化函数:[%v]", c.Name)
	}

	w, err := initializer(c.Attrs)
	if err != nil {
		return nil, err
	}

	if len(c.Items) == 0 {
		return w, err
	}

	cont, ok := w.(writer.WriterContainer)
	if !ok {
		return nil, fmt.Errorf("[%v]并未实现writer.WriterContainer接口", c.Name)
	}

	for _, cfg := range c.Items {
		wr, err := cfg.toWriter()
		if err != nil {
			return nil, err
		}
		cont.AddWriter(wr)
	}

	return w, nil
}
