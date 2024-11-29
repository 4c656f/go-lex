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
		{
			name:  "keywoards",
			input: `and class else false for fun if nil or return super this true var while`,
			expectedLines: []string{
				"AND and null",
				"CLASS class null",
				"ELSE else null",
				"FALSE false null",
				"FOR for null",
				"FUN fun null",
				"IF if null",
				"NIL nil null",
				"OR or null",
				"RETURN return null",
				"SUPER super null",
				"THIS this null",
				"TRUE true null",
				"VAR var null",
				"WHILE while null",
				"EOF  null",
			},
		},
		{
			name: "identifiers",
			input: `andy formless fo _ _123 _abc ab123
abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_`,
			expectedLines: []string{
				"IDENTIFIER andy null",
				"IDENTIFIER formless null",
				"IDENTIFIER fo null",
				"IDENTIFIER _ null",
				"IDENTIFIER _123 null",
				"IDENTIFIER _abc null",
				"IDENTIFIER ab123 null",
				"IDENTIFIER abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_ null",
				"EOF  null",
			},
		},
		{
			name: "spaces",
			input: `space    tabs				newlines




end`,
			expectedLines: []string{
				"IDENTIFIER space null",
				"IDENTIFIER tabs null",
				"IDENTIFIER newlines null",
				"IDENTIFIER end null",
				"EOF  null",
			},
		},
		{
			name:  "punctuators",
			input: `(){};,+-*!===<=>=!=<>/.`,
			expectedLines: []string{
				"LEFT_PAREN ( null",
				"RIGHT_PAREN ) null",
				"LEFT_BRACE { null",
				"RIGHT_BRACE } null",
				"SEMICOLON ; null",
				"COMMA , null",
				"PLUS + null",
				"MINUS - null",
				"STAR * null",
				"BANG_EQUAL != null",
				"EQUAL_EQUAL == null",
				"LESS_EQUAL <= null",
				"GREATER_EQUAL >= null",
				"BANG_EQUAL != null",
				"LESS < null",
				"GREATER > null",
				"SLASH / null",
				"DOT . null",
				"EOF  null",
			},
		},
		{
			name:  "unterminated",
			input: `"foo" "unterminated`,
			expectedLines: []string{
				`STRING "foo" foo`,
				`EOF  null`,
			},
		},
		{
			name: "strings",
			input: `""
"string"`,
			expectedLines: []string{
				`STRING "" `,
				`STRING "string" string`,
				`EOF  null`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			l.Lex()
			
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
