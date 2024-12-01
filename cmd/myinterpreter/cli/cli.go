package cli

import "os"

type OperationType string

const (
	Parse    = "parse"
	Tokenize = "tokenize"
	Eval     = "evaluate"
	Run      = "run"
)

type ParsedArgs struct {
	Type     OperationType
	FilePath string
}

func ParseArgs() (*ParsedArgs, error) {
	if len(os.Args) < 3 {
		return nil, WrongAmountOfArgsError
	}
	commandType := os.Args[1]
	filePath := os.Args[2]

	return &ParsedArgs{
		Type:     OperationType(commandType),
		FilePath: filePath,
	}, nil
}
