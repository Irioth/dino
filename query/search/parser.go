package search

import (
	"fmt"
	"io"
	"strings"
)

type Query struct {
	table string
	ops   []Operation
}

type Operation interface{}

type Where struct {
	column string
	cmpOp  Token
	value  string
}

type Comparison int

func ParseQuery(query string) (*Query, error) {
	parser := NewParser(strings.NewReader(query))
	return parser.Parse()
}

type Parser struct {
	s *Scanner
}

// NewParser returns a new instance of Parser
func NewParser(r io.Reader) *Parser {
	return &Parser{NewScanner(r)}
}

func (p *Parser) Parse() (*Query, error) {
	q := &Query{}
	tok, lit := p.scanIgnoreWhiteSpace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	q.table = lit

	for {
		tok, lit := p.scanIgnoreWhiteSpace()
		if tok == EOF {
			break
		}
		if tok != PIPE {
			return nil, fmt.Errorf("found %q, expected | (pipe)", lit)
		}
		op, err := p.parseOperation()
		if err != nil {
			return nil, err
		}
		q.ops = append(q.ops, op)
	}
	return q, nil
}

func (p *Parser) parseOperation() (Operation, error) {
	tok, lit := p.scanIgnoreWhiteSpace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected operation name", lit)
	}
	switch lit {
	case "where":
		return p.parseWhere()
	default:
		return nil, fmt.Errorf("unknown operation name %q", lit)
	}
}

func (p *Parser) parseWhere() (*Where, error) {
	w := &Where{}
	tok, lit := p.scanIgnoreWhiteSpace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected column name", lit)
	}
	w.column = lit

	tok, lit = p.scanIgnoreWhiteSpace()
	if tok != EQUAL {
		return nil, fmt.Errorf("found %q, expected equal operator", lit)
	}
	w.cmpOp = tok

	tok, lit = p.scanIgnoreWhiteSpace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected value", lit)
	}
	w.value = lit

	return w, nil
}

func (p *Parser) scanIgnoreWhiteSpace() (Token, string) {
	tok, lit := p.scan()
	if tok == WS {
		return p.scan()
	}
	return tok, lit
}

func (p *Parser) scan() (Token, string) {
	return p.s.Scan()
}
