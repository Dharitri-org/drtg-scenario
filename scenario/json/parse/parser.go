package scenjsonparse

import (
	fr "github.com/Dharitri-org/drtg-scenario/scenario/expression/fileresolver"
	ei "github.com/Dharitri-org/drtg-scenario/scenario/expression/interpreter"
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
func NewParser(fileResolver fr.FileResolver, vmType []byte) Parser {
	return Parser{
		ExprInterpreter: ei.ExprInterpreter{
			FileResolver: fileResolver,
			VMType:       vmType,
		},
		AllowDctTxLegacySyntax:           true,
		AllowDctLegacySetSyntax:          true,
		AllowDctLegacyCheckSyntax:        true,
		AllowSingleValueInCheckValueList: true,
	}
}
