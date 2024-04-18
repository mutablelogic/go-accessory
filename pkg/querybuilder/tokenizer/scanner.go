package tokenizer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// Scanner represents a lexical scanner.
type Scanner struct {
	r     *bufio.Reader
	pos   Pos
	buf   []*Token
	flags Flags
}

// Flags are the features enabled for the scanner
type Flags uint

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	tokDefaultCapacity = 20
)

const (
	None          Flags = 0           // No flags
	HashComments  Flags = (1 << iota) // Enable comments which start with a hash and end in newline
	LineComments                      // Enable comments which start with a double slash and end in newline
	BlockComments                     // Enable comments which start with a slash and asterisk, and end With an asterisk and slash
	SQLComments                       // Enable comments which start with a double dash and end in newline
	NewlineToken                      // Enable newline tokens (without this flag, they are recognized as whitespace)
)

var (
	hexDigits = &unicode.RangeTable{
		R16: []unicode.Range16{
			{Lo: 0x0030, Hi: 0x0039, Stride: 1}, // 0-9
			{Lo: 0x0041, Hi: 0x0046, Stride: 1}, // A-F
			{Lo: 0x0061, Hi: 0x0066, Stride: 1}, // a-f
		},
	}
	octalDigits = &unicode.RangeTable{
		R16: []unicode.Range16{
			{Lo: 0x0030, Hi: 0x0037, Stride: 1}, // 0-7
		},
	}
	binaryDigits = &unicode.RangeTable{
		R16: []unicode.Range16{
			{Lo: 0x0030, Hi: 0x0031, Stride: 1}, // 0-1
		},
	}
	numberPrefix = &unicode.RangeTable{
		R16: []unicode.Range16{
			{Lo: 0x002B, Hi: 0x002B, Stride: 1}, // +
			{Lo: 0x002D, Hi: 0x002E, Stride: 1}, // - .
		},
	}
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader, pos Pos, flags Flags) *Scanner {
	return &Scanner{
		r:   bufio.NewReader(r),
		pos: pos,
		buf: make([]*Token, 0, tokDefaultCapacity),
	}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Tokens returns tokens from the scanner until EOF or illegal input is encountered
func (s *Scanner) Tokens() ([]*Token, error) {
	result := make([]*Token, 0, tokDefaultCapacity)
	for {
		tok := s.next()
		if tok == nil {
			return nil, NewPosError(ErrBadParameter.With("Illegal input"), s.pos)
		} else if tok.Kind == EOF {
			break
		}
		result = append(result, tok)
	}
	// Return tokens
	return result, nil
}

// Next returns the next token. If the scanner is at EOF, continue to return EOF
func (s *Scanner) Next() *Token {
	// Obtain the next token
	token := s.Peak()

	// If it's not an EOF token then remove it from the buffer
	if token.Kind != EOF {
		s.buf = s.buf[1:]
	}

	// Return the token
	return token
}

// Peak returns the next token without advancing the scanner. If the scanner is
// at EOF, Peak continues to return EOF
func (s *Scanner) Peak() *Token {
	// If there is a single EOF token in the buffer, return it
	if len(s.buf) == 1 && s.buf[0].Kind == EOF {
		return s.buf[0]
	}
	// Consume a token, add it to the buffer and return it
	token := s.next()
	s.buf = append(s.buf, token)
	return token
}

// NewError returns an error with positional information
func (s *Scanner) NewError(err error) error {
	return NewPosError(err, s.pos)
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Returns true if the flag is set
func (f Flags) is(flag Flags) bool {
	return f&flag > 0
}

// Returns the next token and literal value.
func (s *Scanner) next() *Token {
	// Read the next rune, and advance the position
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if ch == '\n' && s.flags.is(NewlineToken) {
		return NewToken(Space, "\n", s.pos)
	} else if unicode.IsSpace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if unicode.IsLetter(ch) {
		s.unread()
		return s.scanIdent()
	} else if unicode.IsDigit(ch) || unicode.Is(numberPrefix, ch) {
		s.unread()
		return s.scanNumber()
	} else if ch == '"' {
		s.unread()
		return s.scanString()
	} else if ch == '\'' {
		s.unread()
		return s.scanString()
	} else if ch == '#' && s.flags.is(HashComments) {
		s.unread()
		return s.scanHashComment()
	} else if ch == '-' && s.flags.is(SQLComments) {
		s.unread()
		return s.scanLineComment(ch)
	} else if ch == '/' && s.flags.is(LineComments) {
		s.unread()
		return s.scanLineComment(ch)
	}

	// Otherwise read the individual character
	if kind, exists := tokenKindMap[ch]; exists {
		return NewToken(kind, string(ch), s.pos)
	} else {
		return nil
	}
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() *Token {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if ch == '\n' && s.flags.is(NewlineToken) {
			s.unread()
			break
		} else if !unicode.IsSpace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return NewToken(Space, buf.String(), s.pos)
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() *Token {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// Reserved words
	keyword := strings.ToUpper(buf.String())
	if kind, exists := tokenKeywordMap[keyword]; exists {
		return NewToken(kind, buf.String(), s.pos)
	}

	// Otherwise return as a regular identifier.
	return NewToken(Ident, buf.String(), s.pos)
}

// scanString consumes a contiguous string of non-quote characters.
// Quote characters can be consumed if they're first escaped with a backslash.
func (s *Scanner) scanString() *Token {
	// Read the delimiter
	ending := s.read()

	// Create a buffer and read the current character into it.
	var buf bytes.Buffer

	// Read every subsequent character into the buffer.
	for {
		if ch := s.read(); ch == eof {
			// Return nil if the string is not terminated
			return nil
		} else if ch == ending {
			break
		} else if ch == '\\' {
			// If the next character is an escape then write the escaped char
			next := s.read()
			if next == eof {
				// Unterminated escape
				return nil
			} else if next == 'n' {
				buf.WriteRune('\n')
			} else if next == 'r' {
				buf.WriteRune('\r')
			} else if next == 't' {
				buf.WriteRune('\t')
			} else if next == '\\' {
				buf.WriteRune('\\')
			} else if next == '"' {
				buf.WriteRune('"')
			} else if next == '\'' {
				buf.WriteRune('\'')
			} else {
				// Invalid escape
				return nil
			}
		} else {
			buf.WriteRune(ch)
		}
	}

	// Return the string
	return NewToken(String, buf.String(), s.pos)
}

// scanNumber consumes all kinds of numbers
func (s *Scanner) scanNumber() *Token {
	// Create a buffer
	var buf bytes.Buffer
	// Exponent is true if we've seen an E
	var Exponent bool
	// Set the default kind
	kind := NumberInteger

	// Read every digit into the buffer
FOR_LOOP:
	for {
		ch := s.read()
		switch {
		case ch == eof:
			// EOF will cause the loop to exit
			break FOR_LOOP
		case kind == NumberInteger && ch == '.':
			// Switch to float
			kind = NumberFloat
		case kind == NumberInteger && unicode.Is(numberPrefix, ch):
			if buf.Len() > 0 {
				// Plus, Minus or Punkt is not at the beginning of the number
				s.unread()
				break FOR_LOOP
			}
		case kind == NumberInteger && unicode.IsDigit(ch):
			// Switch to octal if first digit is zero
			if buf.String() == "0" || buf.String() == "-0" || buf.String() == "+0" {
				kind = NumberOctal
			}
		case kind == NumberHex:
			if !unicode.Is(hexDigits, ch) {
				return nil
			}
		case kind == NumberOctal:
			// Continuation of octal
			if !unicode.Is(octalDigits, ch) {
				return nil
			}
		case kind == NumberBinary:
			// Continuation of binary
			if !unicode.Is(binaryDigits, ch) {
				return nil
			}
		case kind == NumberFloat && unicode.IsDigit(ch):
			// Continuation of float
		case (kind == NumberFloat || kind == NumberInteger) && buf.Len() > 0 && (ch == 'e' || ch == 'E'):
			// Switch to exponent mode
			kind = NumberFloat
			if Exponent {
				return nil
			} else {
				Exponent = true
			}
		case kind == NumberInteger && buf.Len() > 0 && (ch == 'x' || ch == 'X'):
			// Switch to octal if first digit is zero
			if buf.String() == "0" || buf.String() == "-0" || buf.String() == "+0" {
				kind = NumberHex
			} else {
				return nil
			}
		case kind == NumberInteger && buf.Len() > 0 && (ch == 'b' || ch == 'B'):
			// Switch to binary if first digit is zero
			if buf.String() == "0" || buf.String() == "-0" || buf.String() == "+0" {
				kind = NumberBinary
			} else {
				return nil
			}
		default:
			s.unread()
			break FOR_LOOP
		}

		// Write the rune into the buffer
		_, _ = buf.WriteRune(ch)
	}

	// Error - no digits
	if buf.Len() == 0 {
		return nil
	}

	// Error when Float ends on an E
	if kind == NumberFloat && buf.String()[buf.Len()-1] == 'e' || buf.String()[buf.Len()-1] == 'E' {
		return nil
	}

	fmt.Printf("   kind=%v buf=%q remaining=%d\n", kind, buf.String(), s.r.Buffered())

	// If number prefix on it's own, then return as a regular ident
	if buf.Len() == 1 {
		ch := rune(buf.String()[0])
		if kind, exists := tokenKindMap[ch]; exists {
			fmt.Println("   returning ", kind, "buf=", buf.String())
			return NewToken(kind, buf.String(), s.pos)
		}
	}

	// Create a new token and return it
	fmt.Println("   returning ", kind, "buf=", buf.String())
	return NewToken(kind, buf.String(), s.pos)
}

// scanHashComment consumes the current rune and all contiguous runes until a newline
func (s *Scanner) scanHashComment() *Token {
	// Create a buffer for the comment
	var buf bytes.Buffer

	// Read every subsequent character into the buffer.
	// Newlines and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if ch == '\n' {
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return NewToken(Comment, buf.String(), s.pos)
}

// scanLineComment consumes the current rune, another identical rune and then all
// contiguous runes until a newline
func (s *Scanner) scanLineComment(delim rune) *Token {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer

	// Count the number of delimiter runes
	count := 0

	// Read every subsequent character into the buffer.
	// Newlines and EOF will cause the loop to exit.
FOR_LOOP:
	for {
		ch := s.read()
		switch {
		case ch == delim && count <= 2:
			// Delimiter for comment
			count++
		case count >= 2 && ch != '\n':
			// Consume comment
			buf.WriteRune(ch)
		case count >= 2 && ch == '\n':
			// End of line, emit comment token
			break FOR_LOOP
		default:
			// Invalid state, return error
			s.unread()
			return nil
		}
	}

	return NewToken(Comment, buf.String(), s.pos)
}

// read the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	// Read a rune
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}

	// Mark previous position
	s.pos.x, s.pos.y = s.pos.Line, s.pos.Col

	// Advance position
	if ch == '\n' {
		s.pos.Line++
		s.pos.Col = 0
	} else if ch != eof {
		s.pos.Col++
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
	// Restore previous position
	s.pos.Line, s.pos.Col, s.pos.x, s.pos.y = s.pos.x, s.pos.y, 0, 0
}
