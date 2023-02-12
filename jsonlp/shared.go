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
		extra.UnicodeWordBoundary(),
		Zero(trueAst),
	)
}

func falseValue() ParserFn {
	return FlatGroup(
		erase(Seq("false")),
		extra.UnicodeWordBoundary(),
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
		extra.UnicodeWordBoundary(),
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
		extra.UnicodeWordBoundary(),
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
		extra.UnicodeWordBoundary(),
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
		extra.UnicodeWordBoundary(),
		func(ctx ParserContext) (ParserContext, error) {
			if ctx.Tag.(parseOptions).interop {
				return nilParser(ctx)
			} else {
				return nanParser(ctx)
			}
		},
	)
}

func radixNumberParser(prefix string, radix int, radixNumbrrStr ParserFn) ParserFn {
	return Trans(
		FlatGroup(
			First(
				CharClass("+", "-"),
				Zero(Ast{Type: AstType_String, Value: ""}),
			),
			FlatGroup(erase(SeqI(prefix))),
			erase(ZeroOrMoreTimes(Seq("_"))),
			First(
				FlatGroup(
					First(
						FlatGroup(
							Trans(
								FlatGroup(
									radixNumbrrStr,
									Seq("."),
									erase(ZeroOrMoreTimes(Seq("_"))),
									ZeroOrOnce(radixNumbrrStr),
								),
								Concat,
							),
							SeqI("p"),
						),
						FlatGroup(
							Trans(
								FlatGroup(
									Seq("."),
									erase(ZeroOrMoreTimes(Seq("_"))),
									radixNumbrrStr,
								),
								Concat,
							),
							SeqI("p"),
						),
						FlatGroup(radixNumbrrStr, SeqI("p")),
					),
					extra.IntegerNumberStr(),
				),
				FlatGroup(
					radixNumbrrStr,
					First(
						SeqI("s64"),
						SeqI("u64"),
						Zero(Ast{Value: "f"}),
					),
				),
			),
		),
		radixNumberTransformer(prefix, radix),
	)
}

func numberValueInner() ParserFn {
	return First(
		FlatGroup(
			First(
				radixNumberParser("0b", 2, extra.BinaryNumberStr()),
				radixNumberParser("0o", 8, extra.OctalNumberStr()),
				radixNumberParser("0x", 16, extra.HexNumberStr()),
				Trans(
					extra.FloatNumberStr(),
					ParseFloat,
					ChangeClassName("Float"),
				),
				Trans(
					FlatGroup(
						erase(ZeroOrOnce(Seq("+"))),
						extra.IntegerNumberStr(),
						First(
							FlatGroup(SeqI("s64")),
							FlatGroup(SeqI("u64")),
							FlatGroup(Zero(Ast{Value: "f"})),
						),
					),
					decimalNumberTransformer,
				),
			),
			extra.UnicodeWordBoundary(),
		),
		positiveInfinityValue(),
		negativeInfinityValue(),
		nanValue(),
	)
}

func numberValue() ParserFn {
	return numberValueInner()
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
			extra.UnicodeWordBoundary(),
		),
		extra.ParseDate,
	)
}

func dateTimeValue() ParserFn {
	return Trans(
		FlatGroup(
			extra.DateTimeStr(),
			extra.UnicodeWordBoundary(),
		),
		extra.ParseDateTime,
	)
}

func timeValue() ParserFn {
	return Trans(
		FlatGroup(
			extra.TimeStr(),
			extra.UnicodeWordBoundary(),
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
					sp0(),
					sp0NoLb(),
				),
				erase(CharClass(".")),
				If(allowLb,
					sp0(),
					sp0NoLb(),
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
