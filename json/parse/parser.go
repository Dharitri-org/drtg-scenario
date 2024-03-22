package scenjsonparse

import (
	ei "github.com/Dharitri-org/drtg-scenario/expression/interpreter"
	fr "github.com/Dharitri-org/drtg-scenario/fileresolver"
)

// Parser performs parsing of both json tests (older) and scenarios (new).
type Parser struct {
	ExprInterpreter                  ei.ExprInterpreter
	AllowDctTxLegacySyntax           bool
	AllowDctLegacySetSyntax          bool
	AllowDctLegacyCheckSyntax        bool
	AllowSingleValueInCheckValueList bool
}

// NewParser provides a new Parser instance.
func NewParser(fileResolver fr.FileResolver) Parser {
	return Parser{
		ExprInterpreter: ei.ExprInterpreter{
			FileResolver: fileResolver,
		},
		AllowDctTxLegacySyntax:           true,
		AllowDctLegacySetSyntax:          true,
		AllowDctLegacyCheckSyntax:        true,
		AllowSingleValueInCheckValueList: true,
	}
}
