package jsonlp

import (
	"unicode"

	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/string"
)

func tomlUnicodeIdentifierStr() ParserFn {
	return Trans(
		FlatGroup(
			OneOrMoreTimes(First(
				// ID_Continue + '$' + U+200C + U+200D + ('-' for TOML)
				CharClass("$", "-"),
				CharClassFn(func(c rune) bool {
					// Alnum(), '_', and ...
					return (unicode.Is(unicode.L, c) ||
						unicode.Is(unicode.Nl, c) ||
						unicode.Is(unicode.Other_ID_Start, c) ||
						unicode.Is(unicode.Mn, c) ||
						unicode.Is(unicode.Mc, c) ||
						unicode.Is(unicode.Nd, c) ||
						unicode.Is(unicode.Pc, c) ||
						unicode.Is(unicode.Other_ID_Continue, c) ||
						c == 0x0200c || c == 0x0200d) &&
						!unicode.Is(unicode.Pattern_Syntax, c) &&
						!unicode.Is(unicode.Pattern_White_Space, c)
				}),
			)),
		),
		Concat,
		ChangeClassName("IdentifierStr"),
	)
}
