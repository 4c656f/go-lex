package lexer

import (
	"testing"
)

func TestLexerNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "integer number",
			input: `123
123.456
.456
123.`,
			expected: 
`NUMBER 123 123.0
NUMBER 123.456 123.456
DOT . null
NUMBER 456 456.0
NUMBER 123 123.0
DOT . null
EOF  null`,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Run number test:", i, tt)
			l := New(tt.input)
			err := l.Lex()
			if err != nil {
				t.Errorf("Lex() error = %v", err)
				return
			}
			lexed := l.String()
			if lexed != tt.expected{
				t.Errorf(`
					Lex() got =
					%v
					want = 
					%v`, lexed, tt.expected)
			}
			
		})
	}
}
