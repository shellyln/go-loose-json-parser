package jsonlp

import (
	. "github.com/shellyln/takenoco/base"
	"github.com/shellyln/takenoco/extra"
	. "github.com/shellyln/takenoco/string"
)

type parseOptions struct {
	interop bool
	isTOML  bool
}

// Remove the resulting AST.
func erase(fn ParserFn) ParserFn {
	return Trans(fn, Erase)
}

// Whitespaces
func sp0() ParserFn {
	return erase(ZeroOrMoreTimes(First(Whitespace(), comment())))
}

// Whitespaces
func sp0NoLb() ParserFn {
	return erase(ZeroOrMoreTimes(First(WhitespaceNoLineBreak(), commentLookAheadLb())))
}

// Whitespaces
func sp1NoLb() ParserFn {
	return erase(OneOrMoreTimes(First(WhitespaceNoLineBreak(), commentLookAheadLb())))
}

func lineComment(lookAheadLb bool) ParserFn {
	return erase(FlatGroup(
		Seq("//"),
		FlatGroup(
			ZeroOrMoreTimes(CharClassN("\r\n", "\n", "\r")),
			First(
				If(lookAheadLb,
					LookAhead(CharClass("\r\n", "\n", "\r")),
					CharClass("\r\n", "\n", "\r"),
				),
				LookAhead(End()),
			),
		),
	))
}

func hashLineComment(lookAheadLb bool) ParserFn {
	return erase(FlatGroup(
		Seq("#"),
		FlatGroup(
			ZeroOrMoreTimes(CharClassN("\r\n", "\n", "\r")),
			First(
				If(lookAheadLb,
					LookAhead(CharClass("\r\n", "\n", "\r")),
					CharClass("\r\n", "\n", "\r"),
				),
				LookAhead(End()),
			),
		),
	))
}

func blockComment() ParserFn {
	return erase(FlatGroup(
		Seq("/*"),
		ZeroOrMoreTimes(CharClassN("*/")),
		First(
			Seq("*/"),
			Error("An unexpected termination has appeared in the block comment."),
		),
	))
}

func comment() ParserFn {
	return First(lineComment(false), hashLineComment(false), blockComment())
}

func commentLookAheadLb() ParserFn {
	return First(lineComment(true), hashLineComment(true), blockComment())
}

func trueValue() ParserFn {
	return FlatGroup(
		erase(Seq("true")),
		WordBoundary(),
		Zero(trueAst),
	)
}

func falseValue() ParserFn {
	return FlatGroup(
		erase(Seq("false")),
		WordBoundary(),
		Zero(falseAst),
	)
}

func boolValue() ParserFn {
	return First(
		trueValue(),
		falseValue(),
	)
}

func nullValue() ParserFn {
	return FlatGroup(
		erase(CharClass("null", "undefined")),
		WordBoundary(),
		Zero(nilAst),
	)
}

func positiveInfinityValue() ParserFn {
	nilParser := Zero(nilAst)
	infParser := Zero(positiveInfinityAst)
	return FlatGroup(
		erase(FlatGroup(
			ZeroOrOnce(Seq("+")),
			CharClass("Infinity", "inf"),
		)),
		WordBoundary(),
		func(ctx ParserContext) (ParserContext, error) {
			if ctx.Tag.(parseOptions).interop {
				return nilParser(ctx)
			} else {
				return infParser(ctx)
			}
		},
	)
}

func negativeInfinityValue() ParserFn {
	nilParser := Zero(nilAst)
	infParser := Zero(negativeInfinityAst)
	return FlatGroup(
		erase(CharClass("-Infinity", "-inf")),
		WordBoundary(),
		func(ctx ParserContext) (ParserContext, error) {
			if ctx.Tag.(parseOptions).interop {
				return nilParser(ctx)
			} else {
				return infParser(ctx)
			}
		},
	)
}

func nanValue() ParserFn {
	nilParser := Zero(nilAst)
	nanParser := Zero(nanAst)
	return FlatGroup(
		erase(CharClass("+nan", "-nan", "NaN", "nan")),
		WordBoundary(),
		func(ctx ParserContext) (ParserContext, error) {
			if ctx.Tag.(parseOptions).interop {
				return nilParser(ctx)
			} else {
				return nanParser(ctx)
			}
		},
	)
}

func numberValue() ParserFn {
	return First(
		FlatGroup(
			First(
				Trans(
					FlatGroup(erase(SeqI("0b")), extra.BinaryNumberStr()),
					ParseIntRadix(2),
					castIntToFloat,
				),
				Trans(
					FlatGroup(erase(SeqI("0o")), extra.OctalNumberStr()),
					ParseIntRadix(8),
					castIntToFloat,
				),
				Trans(
					FlatGroup(erase(SeqI("0x")), extra.HexNumberStr()),
					ParseIntRadix(16),
					castIntToFloat,
				),
				Trans(
					extra.FloatNumberStr(),
					ParseFloat,
					ChangeClassName("Float"),
				),
				Trans(
					extra.IntegerNumberStr(),
					ParseInt,
					castIntToFloat,
				),
			),
			WordBoundary(),
		),
		positiveInfinityValue(),
		negativeInfinityValue(),
		nanValue(),
	)
}

func stringValue() ParserFn {
	tomlStr := tomlStringValue()
	jsonStr := jsonStringValue()
	return func(ctx ParserContext) (ParserContext, error) {
		if ctx.Tag.(parseOptions).isTOML {
			return tomlStr(ctx)
		} else {
			return jsonStr(ctx)
		}
	}
}

func dateValue() ParserFn {
	return Trans(
		FlatGroup(
			extra.DateStr(),
			WordBoundary(),
		),
		extra.ParseDate,
	)
}

func dateTimeValue() ParserFn {
	return Trans(
		FlatGroup(
			extra.DateTimeStr(),
			WordBoundary(),
		),
		extra.ParseDateTime,
	)
}

func timeValue() ParserFn {
	return Trans(
		FlatGroup(
			extra.TimeStr(),
			WordBoundary(),
		),
		extra.ParseTime,
	)
}

func primitiveValue() ParserFn {
	return First(
		stringValue(),
		boolValue(),
		nullValue(),
		timeValue(),
		dateTimeValue(),
		dateValue(),
		numberValue(),
	)
}

func identifier() ParserFn {
	return Trans(
		tomlUnicodeIdentifierStr(),
		ChangeClassName("Idenitifier"),
	)
}

func dottedIdentifier(allowLb bool) ParserFn {
	return Trans(
		FlatGroup(
			First(
				stringValue(),
				identifier(),
			),
			OneOrMoreTimes(
				If(allowLb,
					erase(sp0()),
					erase(sp0NoLb()),
				),
				erase(CharClass(".")),
				If(allowLb,
					erase(sp0()),
					erase(sp0NoLb()),
				),
				First(
					stringValue(),
					identifier(),
				),
			),
		),
		func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
			length := len(asts)
			v := make([]string, length)
			for i := 0; i < length; i++ {
				v[i] = asts[i].Value.(string)
			}
			return AstSlice{{
				ClassName: "DottedIdenitifier",
				Type:      AstType_ListOfAny,
				Value:     v,
			}}, nil
		},
	)
}
