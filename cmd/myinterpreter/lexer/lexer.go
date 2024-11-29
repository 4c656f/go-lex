package lexer

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type Lexer struct {
	source string
	start  int
	end    int
	line   int
	tokens []*token.Token
	errors []error
}

func New(src string) *Lexer {
	return &Lexer{
		line:   1,
		source: src,
	}
}

func (l *Lexer) Lex() {
	for l.hasNext() {
		char := l.advance()
		switch char {
		case '(':
			l.addToken(token.NewToken(token.LEFT_PAREN, l.line, "(", token.NewNullValue()))
		case ')':
			l.addToken(token.NewToken(token.RIGHT_PAREN, l.line, ")", token.NewNullValue()))
		case '{':
			l.addToken(token.NewToken(token.LEFT_BRACE, l.line, "{", token.NewNullValue()))
		case '}':
			l.addToken(token.NewToken(token.RIGHT_BRACE, l.line, "}", token.NewNullValue()))
		case ';':
			l.addToken(token.NewToken(token.SEMICOLON, l.line, ";", token.NewNullValue()))
		case ',':
			l.addToken(token.NewToken(token.COMMA, l.line, ",", token.NewNullValue()))
		case '+':
			l.addToken(token.NewToken(token.PLUS, l.line, "+", token.NewNullValue()))
		case '-':
			l.addToken(token.NewToken(token.MINUS, l.line, "-", token.NewNullValue()))
		case '*':
			l.addToken(token.NewToken(token.STAR, l.line, "*", token.NewNullValue()))
		case '!':
			if l.matchCur('=') {
				l.addToken(token.NewToken(token.BANG_EQUAL, l.line, "!=", token.NewNullValue()))
			} else {
				l.addToken(token.NewToken(token.BANG, l.line, "!", token.NewNullValue()))
			}
		case '=':
			if l.matchCur('=') {
				l.addToken(token.NewToken(token.EQUAL_EQUAL, l.line, "==", token.NewNullValue()))
			} else {
				l.addToken(token.NewToken(token.EQUAL, l.line, "=", token.NewNullValue()))
			}
		case '<':
			if l.matchCur('=') {
				l.addToken(token.NewToken(token.LESS_EQUAL, l.line, "<=", token.NewNullValue()))
			} else {
				l.addToken(token.NewToken(token.LESS, l.line, "<", token.NewNullValue()))
			}
		case '>':
			if l.matchCur('=') {
				l.addToken(token.NewToken(token.GREATER_EQUAL, l.line, ">=", token.NewNullValue()))
			} else {
				l.addToken(token.NewToken(token.GREATER, l.line, ">", token.NewNullValue()))
			}
		case '/':
			//match comment
			if l.matchCur('/') {
				for l.hasNext() && l.peek() != '\n' {
					l.advance()
				}
				continue
			}
			l.addToken(token.NewToken(token.SLASH, l.line, "/", token.NewNullValue()))
		case '.':
			l.addToken(token.NewToken(token.DOT, l.line, ".", token.NewNullValue()))
		case '"':
			token, err := l.lexString()
			if err != nil {
				l.onError(err)
				continue
			}
			l.addToken(token)
		case ' ', '\r', '\t':
			// Ignore whitespace
			continue
		case '\n':
			l.line++
		default:
			if isAlpha(char) {
				token, err := l.lexIdent()
				if err != nil {
					l.onError(err)
					continue
				}
				l.addToken(token)
				continue
			}
			if isNumeric(char) {
				token, err := l.lexNumber()
				if err != nil {
					l.onError(err)
					continue
				}
				l.addToken(token)
				continue
			}
			l.onError(NewLexError(l.line, "Unexpected character", string(char)))
		}
	}
	l.addToken(token.NewToken(token.EOF, l.line, "", token.NewNullValue()))

}

func (l *Lexer) addToken(t *token.Token) {
	l.tokens = append(l.tokens, t)
}

func (l *Lexer) onError(err error) {
	l.errors = append(l.errors, err)
}

func (l *Lexer) HasErrors() ([]error, bool) {
	if l.errors == nil {
		return nil, false
	}
	return l.errors, true
}

func (l *Lexer) lexString() (*token.Token, error) {
	l.start = l.end - 1
	startLine := l.line
	for l.hasNext() {
		cur := l.peek()
		if cur == '\n' {
			l.line++
			continue
		}
		if cur == '"' {
			break
		}
		l.advance()
	}
	if !l.hasNext() {
		return nil, nil
	}
	l.advance()
	return token.NewToken(token.STRING, startLine, l.source[l.start:l.end], token.NewStringValue(l.source[l.start+1:l.end-1])), nil
}

func (l *Lexer) lexNumber() (*token.Token, error) {
	//start of substr should start from prev advanced token
	l.start = l.end - 1
	for l.hasNext() && isNumeric(l.peek()) {
		l.advance()
	}

	if l.peek() == '.' && isNumeric(l.peekNext()) {
		//consume '.'
		l.advance()
		for l.hasNext() && isNumeric(l.peek()) {
			l.advance()
		}
		parsedFloat, err := strconv.ParseFloat(l.source[l.start:l.end], 64)
		return token.NewToken(token.NUMBER, l.line, l.source[l.start:l.end], token.NewFloatValue(parsedFloat)), err
	}

	parsedInt, err := strconv.Atoi(l.source[l.start:l.end])
	return token.NewToken(token.NUMBER, l.line, l.source[l.start:l.end], token.NewIntValue(parsedInt)), err
}

func (l *Lexer) lexIdent() (*token.Token, error) {
	l.start = l.end - 1
	for l.hasNext() && isAlpaNumeric(l.peek()) {
		l.advance()
	}
	literal := l.source[l.start:l.end]
	tokType, ok := token.MatchStringToKeywoard(literal)
	if !ok {
		return token.NewToken(token.IDENTIFIER, l.line, literal, token.NewNullValue()), nil
	}
	return token.NewToken(tokType, l.line, literal, token.NewNullValue()), nil
}

func (l Lexer) hasNext() bool {
	return l.end < len(l.source)
}

func (l *Lexer) advance() byte {
	c := l.source[l.end]
	l.end++
	return c
}

func (l *Lexer) matchCur(char byte) bool {
	if !l.hasNext() {
		return false
	}
	matched := char == l.peek()
	if !matched {
		return false
	}
	l.advance()
	return true
}

func (l *Lexer) peek() byte {
	if !l.hasNext() {
		return 0
	}
	return l.source[l.end]
}

func (l *Lexer) peekNext() byte {
	nxtIdx := l.end + 1
	if nxtIdx >= len(l.source) {
		return 0
	}
	return l.source[nxtIdx]
}

func (l Lexer) String() string {
	strSlice := make([]string, len(l.tokens))
	for i, t := range l.tokens {
		strSlice[i] = t.String()
	}
	return strings.Join(strSlice, "\n")
}

func isAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '_'
}

func isNumeric(char byte) bool {
	return char >= '0' && char <= '9'
}

func isAlpaNumeric(char byte) bool {
	return isAlpha(char) || isNumeric(char)
}

func Report(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
}
