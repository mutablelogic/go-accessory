package tokenizer

import (
	"fmt"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// Classifies the kind of token
type TokenKind uint

// Token is decomposed from []byte stream to represent a kind of
// token and the vaoue of the token
type Token struct {
	Kind TokenKind
	Val  string
	Pos  Pos
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	Any TokenKind = iota
	String
	Expr
	Space
	Newline
	Ident
	NumberInteger
	NumberFloat
	NumberOctal
	NumberHex
	NumberBinary
	Punkt
	Question
	Colon
	SemiColon
	Comma
	OpenParen
	CloseParen
	OpenSquare
	CloseSquare
	OpenBrace
	CloseBrace
	Ampersand
	Equal
	Less
	Greater
	Plus
	Minus
	Multiply
	Divide
	Not
	True
	False
	Null
	Comment
	EOF
	Lowest = Equal // Lowest precedence
)

var (
	// Special end of file rune
	eof = rune(0)

	// Special characters
	tokenKindMap = map[rune]TokenKind{
		'.': Punkt,
		'?': Question,
		':': Colon,
		';': SemiColon,
		',': Comma,
		'(': OpenParen,
		')': CloseParen,
		'[': OpenSquare,
		']': CloseSquare,
		'{': OpenBrace,
		'}': CloseBrace,
		'&': Ampersand,
		'=': Equal,
		'<': Less,
		'>': Greater,
		'!': Not,
		'+': Plus,
		'-': Minus,
		'*': Multiply,
		'/': Divide,
		eof: EOF,
	}
	// Reserved words
	tokenKeywordMap = map[string]TokenKind{
		"true":  True,
		"false": False,
		"null":  Null,
	}
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewToken(kind TokenKind, val string, pos Pos) *Token {
	return &Token{Kind: kind, Val: val, Pos: pos}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (k TokenKind) String() string {
	switch k {
	case Any:
		return "Any"
	case String:
		return "String"
	case Expr:
		return "Expr"
	case Space:
		return "Space"
	case Newline:
		return "Newline"
	case Ident:
		return "Ident"
	case NumberInteger:
		return "NumberInteger"
	case NumberFloat:
		return "NumberFloat"
	case NumberOctal:
		return "NumberOctal"
	case NumberHex:
		return "NumberHex"
	case NumberBinary:
		return "NumberBinary"
	case Punkt:
		return "Punkt"
	case Question:
		return "Question"
	case Colon:
		return "Colon"
	case SemiColon:
		return "SemiColon"
	case Comma:
		return "Comma"
	case Comment:
		return "Comment"
	case OpenParen:
		return "OpenParen"
	case CloseParen:
		return "CloseParen"
	case OpenSquare:
		return "OpenSquare"
	case CloseSquare:
		return "CloseSquare"
	case OpenBrace:
		return "OpenBrace"
	case CloseBrace:
		return "CloseBrace"
	case Ampersand:
		return "Ampersand"
	case Equal:
		return "Equal"
	case Less:
		return "Less"
	case Greater:
		return "Greater"
	case Plus:
		return "Plus"
	case Minus:
		return "Minus"
	case Multiply:
		return "Star"
	case Divide:
		return "Slash"
	case Not:
		return "Not"
	case True:
		return "True"
	case False:
		return "False"
	case Null:
		return "Null"
	case EOF:
		return "EOF"
	default:
		return "[?? Invalid TokenKind value]"
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%v<%q>", t.Kind, t.Val)
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (t *Token) toString() string {
	return t.Val
}
