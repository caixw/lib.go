// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package ini

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// 表示ini的语法错误信息
type SyntaxError struct {
	Line int
	Msg  string
}

func (s *SyntaxError) Error() string {
	return fmt.Sprintf("unicode/ini，在第%d行发生语法错误：%v", s.Line, s.Msg)
}

// 几种元素类型
const (
	Undefined = iota // 未定义，初始状态
	Element
	Section
	Comment
	EOF // 已经读取完毕。
)

type Token struct {
	Type  int
	Value string
	Key   string
}

func (t *Token) reset() {
	t.Type = Undefined
	t.Value = t.Value[:0]
	t.Key = t.Key[:0]
}

// 复制一个新的Token
func (t *Token) Copy() *Token {
	return &Token{
		Type:  t.Type,
		Value: t.Value,
		Key:   t.Key,
	}
}

// ini数据的读取操作类。
// 注释只支持以`#`,`;`开头的行，不支持行尾注释；
//
// 读取的内容，Element元素的首尾空格将被去除
// 但注释不会作此处理，包括换行符都将原样输出。
type Reader struct {
	reader *bufio.Reader
	eof    bool // 已经读取完毕
	line   int  // 当前正在处理的行数。
	token  *Token
}

// 从一个io.Reader初始化Reader
func NewReader(r io.Reader) *Reader {
	return &Reader{reader: bufio.NewReader(r), token: &Token{}}
}

// 从一个[]byte初始化Reader
func NewReaderBytes(data []byte) *Reader {
	return NewReader(bytes.NewReader(data))
}

// 从一个字符串初始化Reader
func NewReaderString(str string) *Reader {
	return NewReader(strings.NewReader(str))
}

// 返回下一个Token，当内容读取完毕之后，将返回Type值为EOF的Token
//
// 返回的Token变量，在下将调用Token()时，数据会被重置，若需要保存
// Token的数据，可使用Token.Copy()函数复制一份。
func (r *Reader) Token() (*Token, error) {
	r.token.reset()

START:
	if r.eof {
		r.token.Type = EOF
		return r.token, nil
	}

	// 读取新行
	buffer, err := r.reader.ReadString('\n')
	r.line++
	if err != nil {
		// 真的发生错误了
		if err != io.EOF {
			return nil, err
		}

		// 读取完毕
		r.eof = true
		if len(buffer) == 0 {
			r.token.Type = EOF
			return r.token, nil
		}
	}

	buffer = strings.TrimLeftFunc(buffer, unicode.IsSpace)
	if len(buffer) == 0 {
		goto START
	}

	return r.parseLine(buffer)
}

// 分析一行字符串，并返回结果。
func (r *Reader) parseLine(line string) (*Token, error) {
	switch line[0] {
	case '[': // section
		line = strings.TrimRightFunc(line, unicode.IsSpace) // 防止]之后有空格
		if line[len(line)-1] != ']' {
			return nil, r.newSyntaxError("section名称没有以]作为结尾")
		}

		r.token.Type = Section
		r.token.Value = line[1 : len(line)-1]
		return r.token, nil
	case '#', ';': // comment
		r.token.Type = Comment
		r.token.Value = line[1:]
		return r.token, nil
	default: // element
		pos := strings.IndexRune(line, '=')
		if pos < 0 {
			return nil, r.newSyntaxError("表达式中未找到`=`符号")
		}

		r.token.Type = Element
		r.token.Key = strings.TrimRightFunc(line[:pos], unicode.IsSpace)
		r.token.Value = strings.TrimSpace(line[pos+1:])
		return r.token, nil
	}

	r.token.Type = EOF
	return r.token, nil
}

// 返回一个语法错误的error
func (r *Reader) newSyntaxError(msg string) error {
	return &SyntaxError{
		Msg:  msg,
		Line: r.line,
	}
}

// 将ini转换成map[string]interface{}返回。
//
// 没有与之相对就的MarshalMap，因为map是无序的，若一个map带了section，则转换结
// 果未必是正确的。
//
// 若section参数不为空，则表示只返回section的内容，若没有对应内容，则返回空值
func UnmarshalMap(data []byte, section string) (map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, &SyntaxError{Msg: "没有内容", Line: 0}
	}

	m := make(map[string]interface{})
	currSection := m

	sectionFlag := false // 是否指定了section值
	if section != "" {
		sectionFlag = true
		currSection = nil
	}

	r := NewReaderBytes(data)
	for {
		token, err := r.Token()
		if err != nil {
			return nil, err
		}

		switch token.Type {
		case Comment:
			continue
		case EOF:
			return m, nil
		case Element:
			if currSection != nil {
				currSection[token.Key] = token.Value
			}
		case Section:
			if sectionFlag {
				if section == token.Value {
					currSection = m
				} else {
					currSection = nil
				}
			} else {
				currSection = make(map[string]interface{})
				m[token.Value] = currSection
			}
		default:
			return nil, errors.New("未知的元素类型")
		}
	}
	return m, nil
}

// 将data中的数据写入v中。注释将被忽略。
func Unmarshal(data []byte, v interface{}) error {
	tree, err := scan(v)
	if err != nil {
		return err
	}

	return tree.unmarshal(NewReaderBytes(data))
}
