package lzjson

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// selItem contains information of token in a selector
type selItem struct {
	typ selItemType
	val string
}

// selItemType represents token type
type selItemType int

const (
	selItemError     selItemType = iota // lex error
	selItemDot                          // dot symbol
	selItemSpace                        // space
	selItemRightBrac                    // ']'
	selItemLeftBrac                     // '['
	selItemProp                         // property name of json object
	selItemNumber                       // numeric value / array key
	selItemString                       // string
	selItemEnd                          // end of selector string
)

const (
	charCap      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charSmallCap = "abcdefghijklmnoperstuvwxyz"
	charNumeric  = "0123456789"
)

const eof = -1

// lexSel returns a lexer for the selector string
func lexSel(input string) *selLexer {
	return &selLexer{
		input: input,
		items: make(chan selItem),
	}
}

// selLexer helps tokenize a selector string
type selLexer struct {
	input string
	state selStateFn
	pos   int
	start int
	width int
	items chan selItem
}

// next returns the next rune in the input.
func (l *selLexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *selLexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume the next rune in the input.
func (l *selLexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// emit passes an item back to the client.
func (l *selLexer) emit(t selItemType) {
	l.items <- selItem{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *selLexer) ignore() {
	l.start = l.pos
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *selLexer) errorf(format string, args ...interface{}) selStateFn {
	l.items <- selItem{selItemError, fmt.Sprintf(format, args...)}
	return nil
}

// run runs the state machine for the lexer.
func (l *selLexer) run() {
	for l.state = selLexText; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *selLexer) nextItem() selItem {
	item := <-l.items
	return item
}

// accept consumes the next rune if it's from the valid set.
func (l *selLexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *selLexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// selStateFn represents the state of the scanner as a function that returns the next state.
type selStateFn func(*selLexer) selStateFn

// selLexText process selector and determine what selStateFn
// to process the upcoming string
func selLexText(l *selLexer) selStateFn {
	next := l.next()
	switch next {
	case '.':
		l.emit(selItemDot)
		return selLexProc
	case '[':
		l.emit(selItemLeftBrac)
		return selLexInsideBrac
	case ']':
		l.emit(selItemRightBrac)
		return selLexText
	case eof:
		break
	default:
		return selLexProc
	}
	if l.pos > l.start {
		l.emit(selItemProp)
	}
	l.emit(selItemEnd)
	return nil
}

// selLexProc process selector string like object property
// name until reaching non-alphanumerical character
func selLexProc(l *selLexer) selStateFn {
	l.acceptRun(charCap + charSmallCap + charNumeric)
	l.emit(selItemProp)
	return selLexText
}

func selLexInsideBrac(l *selLexer) selStateFn {
	for {
		switch l.next() {
		case ']':
			l.backup()
			if l.pos > l.start {
				l.emit(selItemNumber)
			}
			return selLexText
		case '"':
			l.ignore()
			return selLexInsideDoubleQuoteString
		case '\'':
			l.ignore()
			return selLexInsideQuoteString
		case eof:
			return l.errorf("unclosed bracket")
		}
	}
}

func selLexInsideQuoteString(l *selLexer) selStateFn {
	for {
		switch l.next() {
		case '\\':
			if l.peek() == '\'' {
				l.next()
			}
		case '\'':
			l.backup()
			if l.pos > l.start {
				l.emit(selItemString)
			}
			l.next()
			l.ignore()
			return selLexText
		case eof:
			return l.errorf("unclosed single quoted string")
		}
	}
}

func selLexInsideDoubleQuoteString(l *selLexer) selStateFn {
	for {
		switch l.next() {
		case '\\':
			if l.peek() == '"' {
				l.next()
			}
		case '"':
			l.backup()
			if l.pos > l.start {
				l.emit(selItemString)
			}
			l.next()
			l.ignore()
			return selLexText
		case eof:
			return l.errorf("unclosed single quoted string")
		}
	}
}
