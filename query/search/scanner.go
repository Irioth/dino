package search

import (
	"bufio"
	"bytes"
	"io"
)

type Token int

const (
	EOF Token = iota
	ILLEGAL
	WS // whitespace
	IDENT
	PIPE
	EQUAL
	INEQUAL
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EUQAL
)

// special marker for eof
const eof = rune(0)

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) Scan() (Token, string) {
	ch := s.read()

	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) || isDigit(ch) {
		s.unread()
		return s.scanIdent()
	}

	switch ch {
	case eof:
		return EOF, ""
	case '|':
		return PIPE, "|"
	case '=':
		return EQUAL, "="
	case '!':
		ch := s.read()
		if ch == '=' {
			return INEQUAL, "!="
		}
		s.unread()
		return ILLEGAL, "!"
	case '>':
		ch := s.read()
		if ch == '=' {
			return GREATER_EUQAL, ">="
		}
		s.unread()
		return GREATER, ">"
	case '<':
		ch := s.read()
		if ch == '=' {
			return LESS_EQUAL, "<="
		}
		s.unread()
		return GREATER, "<"
	}

	return ILLEGAL, string(ch)
}

// isWhitespace returns true if the rune is a space, tab, or newline.
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

func (s *Scanner) scanWhitespace() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

loop:
	for {
		ch := s.read()
		switch {
		case ch == eof:
			break loop
		case !isWhitespace(ch):
			s.unread()
			break loop
		default:
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

func (s *Scanner) scanIdent() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

loop:
	for {
		ch := s.read()
		switch {
		case ch == eof:
			break loop
		case !isLetter(ch) && !isDigit(ch) && ch != '_':
			s.unread()
			break loop
		default:
			buf.WriteRune(ch)
		}
	}

	return IDENT, buf.String()
}

func (s *Scanner) read() rune {
	if ch, _, err := s.r.ReadRune(); err == nil {
		return ch
	}
	return eof
}

func (s *Scanner) unread() {
	if err := s.r.UnreadRune(); err != nil {
		panic(err) // never
	}
}
