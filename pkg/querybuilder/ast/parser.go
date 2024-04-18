package ast

import (
	"fmt"
	"io"
	"strconv"

	// Packages
	"github.com/mutablelogic/go-accessory/pkg/querybuilder/tokenizer"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Parser struct {
	*tokenizer.Scanner
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader, pos tokenizer.Pos) *Parser {
	return &Parser{
		tokenizer.NewScanner(r, pos),
	}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Parse is the main function call given a lexer instance it will parse
// and construct an abstract syntax tree (AST) for the given input.
func (p *Parser) Parse() (Node, error) {
	result := &exprListNode{}
	for {
		tok := p.Next()
		if tok == nil {
			return nil, fmt.Errorf("Syntax error")
		} else if tok.Kind == tokenizer.EOF {
			break
		} else if tok.Kind == tokenizer.Space || tok.Kind == tokenizer.Comment {
			// Ignore whitespace and comments
			continue
		} else if node := p.parse(tok); node != nil {
			result.v = append(result.v, node)
		}
	}

	// Return the parsed expression
	return result, nil
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (p *Parser) parse(t *tokenizer.Token) Node {
	switch t.Kind {
	case tokenizer.NumberFloat:
		if f, err := strconv.ParseFloat(t.Val.(string), 64); err != nil {
			p.err = p.s.NewError(ErrBadParameter.Withf("Unexpected number: %q", t.Val))
			return nil
		} else {
			return &FloatNumberNode{v: f}
		}
	case tokenizer.NumberHex, tokenizer.NumberOctal, tokenizer.NumberInteger, tokenizer.NumberBinary:
		if f, err := strconv.ParseInt(t.Val.(string), 0, 64); err != nil {
			p.err = p.s.NewError(ErrBadParameter.Withf("Unexpected number: %q", t.Val))
			return nil
		} else {
			return &IntNumberNode{v: f}
		}
	case tokenizer.True:
		return &BooleanNode{v: true}
	case tokenizer.False:
		return &BooleanNode{v: false}
	case tokenizer.Null:
		return &NullNode{}
	case tokenizer.String:
		return &StringNode{v: t.Val.(string)}
	}

	// Unhandled parse case
	//p.err = p.s.NewError(ErrBadParameter.Withf("Unexpected value: %q", t.Val))
	return nil
}
