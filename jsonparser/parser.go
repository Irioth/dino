package jsonparser

import (
	"fmt"
	"io"
	"math"
	"strconv"
)

type Config struct {
}

type Iterator struct {
	cfg        Config
	r          io.Reader
	buf        []byte
	head, tail int
	Error      error
}

type ValueType byte

const (
	ILLEGAL ValueType = iota
	EOF
	STRING
	NUMBER
	OBJECT
	ARRAY
	TRUE
	FALSE
	NULL
)

func NewIterator(cfg *Config, r io.Reader) *Iterator {
	if cfg == nil {
		cfg = &Config{}
	}

	return &Iterator{
		cfg: *cfg,
		r:   r,
		buf: make([]byte, 10),
	}
}

func (iter *Iterator) Skip() bool {
	switch iter.WhatIsNext() {
	case NULL:
		return iter.readTokens('n', 'u', 'l', 'l')
	case TRUE:
		return iter.readTokens('t', 'r', 'u', 'e')
	case FALSE:
		return iter.readTokens('f', 'a', 'l', 's', 'e')
	case STRING:
		iter.ReadString()
		return iter.Error == nil
	case OBJECT:
		return iter.ReadObject(skipFieldCallback)
	case NUMBER:
		iter.ReadNumber()
		return iter.Error == nil
	case ARRAY:
		return iter.ReadArray((*Iterator).Skip)
	}
	return false
}

func skipFieldCallback(iter *Iterator, field string) bool {
	return iter.Skip()
}

func (iter *Iterator) WhatIsNext() ValueType {
	c := iter.nextToken()
	switch c {
	case '{':
		return OBJECT
	case '[':
		return ARRAY
	case '"':
		return STRING
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return NUMBER
	case 't':
		return TRUE
	case 'f':
		return FALSE
	case 'n':
		return NULL
	case eof:
		return EOF
	}
	return ILLEGAL
}

func (iter *Iterator) reportError(format string, params ...interface{}) {
	if iter.Error == nil {
		iter.Error = fmt.Errorf(format, params...)
	}
}

func (iter *Iterator) ReadArray(callback func(iter *Iterator) bool) bool {
	if !iter.readToken('[') {
		iter.reportError("'[' expected, but '%c' found", iter.nextToken())
		return false
	}
	if iter.readToken(']') {
		return true
	}
	for {
		if !callback(iter) {
			return false
		}
		if !iter.readToken(',') {
			break
		}
	}
	if !iter.readToken(']') {
		iter.reportError("']' expected, but '%c' found", iter.nextToken())
		return false
	}
	return true
}

func (iter *Iterator) ReadObject(callback func(iter *Iterator, field string) bool) bool {
	if !iter.readToken('{') {
		iter.reportError("'{' expected, but '%c' found", iter.nextToken())
		return false
	}
	if iter.readToken('}') {
		return true
	}
	for {
		field := iter.readKey()
		if !iter.readToken(':') {
			iter.reportError("':' expected, but '%c' found", iter.nextToken())
			return false
		}
		if !callback(iter, field) {
			return false
		}
		if !iter.readToken(',') {
			break
		}
	}
	if !iter.readToken('}') {
		iter.reportError("'}' expected, but '%c' found", iter.nextToken())
		return false
	}

	return true
}

func (iter *Iterator) readKey() string {
	return string(iter.ReadString())
}

func (iter *Iterator) ReadString() []byte {
	if !iter.readToken('"') {
		iter.reportError("'\"' expected, but '%c' found", iter.nextToken())
		return nil
	}
	offset := 0
	for {
		for i := iter.head + offset; i < iter.tail; i++ {
			if iter.buf[i] == '"' {
				if i == iter.head || iter.buf[i-1] != '\\' {
					s := iter.buf[iter.head:i]
					iter.head = i + 1
					return s
				}
			}
			// control char must be escaped
			if iter.buf[i] <= 0x1f {
				iter.reportError("Unescaped control character")
			}
		}
		offset = iter.tail - iter.head

		iter.growBuffer()

		if err := iter.readMore(); err != nil {
			iter.Error = err
			return nil
		}
	}
}

func (iter *Iterator) growBuffer() {
	nbuf := iter.buf
	if iter.tail-iter.head > len(iter.buf)/2 {
		nbuf = make([]byte, len(iter.buf)*2)
	}

	iter.tail = copy(nbuf, iter.buf[iter.head:iter.tail])
	iter.head = 0
	iter.buf = nbuf
}

func (iter *Iterator) ReadNumber() float64 {
	iter.nextToken()
	negative := false
	if iter.nextChar() == '-' {
		iter.head++
		negative = true
	}
	v := 0.0

	if iter.nextChar() >= '1' && iter.nextChar() <= '9' {
		for iter.nextChar() >= '0' && iter.nextChar() <= '9' {
			v = v*10 + float64(iter.nextChar()-'0')
			iter.head++
		}
	} else if iter.nextChar() != '0' {
		iter.reportError("number expected, but '%c' found", iter.nextChar())
		return 0
	} else {
		iter.head++
	}

	exp := 0

	if iter.nextChar() == '.' {
		iter.head++
		if iter.nextChar() < '0' || iter.nextChar() > '9' {
			iter.reportError("Fraction expected, but '%c' found", iter.nextChar())
			return 0
		}
		for iter.nextChar() >= '0' && iter.nextChar() <= '9' {
			v = v*10 + float64(iter.nextChar()-'0')
			exp--
			iter.head++
		}
	}

	if iter.nextChar() == 'e' || iter.nextChar() == 'E' {
		iter.head++
		negative := false
		if iter.nextChar() == '+' {
			iter.head++
		} else if iter.nextChar() == '-' {
			iter.head++
			negative = true
		}

		exp1 := 0
		if iter.nextChar() < '0' || iter.nextChar() > '9' {
			iter.reportError("Exponent expected, but '%c' found", iter.nextChar())
			return 0
		}
		for iter.nextChar() >= '0' && iter.nextChar() <= '9' {
			exp1 = exp1*10 + int(iter.nextChar()-'0')
			iter.head++
		}

		if negative {
			exp -= exp1
		} else {
			exp += exp1
		}
	}

	if negative {
		v = -v
	}

	return v * math.Pow10(exp)

	// offset := 0
	// for {
	// 	for i := iter.head + offset; i < iter.tail; i++ {
	// 		ch := iter.buf[i]

	// 		if (ch < '0' || ch > '9') && ch != '.' && ch != '+' && ch != '-' && ch != 'e' && ch != 'E' {
	// 			s := iter.buf[iter.head:i]
	// 			iter.head = i
	// 			f, err := parseNumber(s)
	// 			if err != nil {
	// 				iter.Error = err
	// 				return 0
	// 			}
	// 			return f
	// 		}
	// 	}
	// 	offset = iter.tail - iter.head

	// 	iter.growBuffer()

	// 	if err := iter.readMore(); err != nil {
	// 		if err == io.EOF {
	// 			s := iter.buf[iter.head:iter.tail]
	// 			iter.head = iter.tail
	// 			f, err := parseNumber(s)
	// 			if err != nil {
	// 				iter.Error = err
	// 				return 0
	// 			}
	// 			return f
	// 		}
	// 		iter.Error = err
	// 		return 0
	// 	}
	// }
}

func parseNumber(d []byte) (float64, error) {
	return strconv.ParseFloat(string(d), 64)
}

func (iter *Iterator) nextChar() byte {
	for iter.head == iter.tail {
		iter.tail = 0
		iter.head = 0
		if err := iter.readMore(); err != nil {
			if err != io.EOF {
				iter.Error = err
			}
			return eof
		}
	}
	return iter.buf[iter.head]
}

func (iter *Iterator) nextToken() byte {
	for {
		for i := iter.head; i < iter.tail; i++ {
			switch iter.buf[i] {
			case ' ', '\t', '\r', '\n':
				continue
			default:
				iter.head = i
				return iter.buf[i]
			}
		}

		iter.tail = 0
		iter.head = 0
		if err := iter.readMore(); err != nil {
			if err != io.EOF {
				iter.Error = err
			}
			return eof
		}
	}
}

func (iter *Iterator) readMore() error {
	n, err := iter.r.Read(iter.buf[iter.tail:])
	if err != nil {
		return err
	}
	iter.tail += n
	return nil
}

const (
	eof = 0
)

func (iter *Iterator) readToken(ch byte) bool {
	if iter.nextToken() == ch {
		iter.head++
		return true
	}
	return false
}
func (iter *Iterator) readTokens(tokens ...byte) bool {
	for _, ch := range tokens {
		if iter.nextToken() != ch {
			return false
		}
		iter.head++
	}
	return true
}
