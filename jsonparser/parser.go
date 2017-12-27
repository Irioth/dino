package jsonparser

import (
	"fmt"
	"io"
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

func (iter *Iterator) reportError(format string, params ...interface{}) {
	if iter.Error == nil {
		iter.Error = fmt.Errorf(format, params...)
	}
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
	last := iter.head
	for {
		for i := last; i < iter.tail; i++ {
			if iter.buf[i] == '"' {
				s := iter.buf[iter.head:i]
				iter.head = i + 1
				return s
			}
		}

		nbuf := iter.buf
		if iter.tail-iter.head > len(iter.buf)/2 {
			nbuf = make([]byte, len(iter.buf)*2)
		}

		iter.tail = copy(nbuf, iter.buf[iter.head:iter.tail])
		last = iter.tail
		iter.head = 0
		iter.buf = nbuf

		if err := iter.readMore(); err != nil {
			iter.Error = err
			return nil
		}
	}
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
			iter.Error = err
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
