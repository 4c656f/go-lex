package lexer

import (
	"testing"
)

func TestLexerNumbers(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedLines []string
	}{
		{
			name:  "integer number",
			input: `123 123.456 .456 123.`,
			expectedLines: []string{
				"NUMBER 123 123.0",
				"NUMBER 123.456 123.456",
				"DOT . null",
				"NUMBER 456 456.0",
				"NUMBER 123 123.0",
				"DOT . null",
				"EOF  null",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			err := l.Lex()
			if err != nil {
				t.Errorf("Lex() error = %v", err)
				return
			}
			if len(l.tokens) != len(tt.expectedLines) {
				t.Errorf("TEST %s Wrong amount of tokens:  %s, %s", tt.name, l, tt.expectedLines)
			}
			for i, tok := range l.tokens {
				if tok.String() != tt.expectedLines[i] {
					t.Errorf(
`TEST %s Lex() got at index %v =
%v
want =
%v
`, tt.name, i, tok, tt.expectedLines[i])
				}
			}

		})
	}
}
