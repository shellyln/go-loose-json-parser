package jsonlp

import (
	. "github.com/shellyln/takenoco/base"
)

func castIntToFloat(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	return AstSlice{{
		ClassName: "Float",
		Type:      AstType_Float,
		Value:     float64(asts[0].Value.(int64)),
	}}, nil
}

func setStr(s string) TransformerFn {
	return func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
		return AstSlice{{
			ClassName: "setStr",
			Type:      AstType_String,
			Value:     s,
		}}, nil
	}
}

func replaceStr(fn ParserFn, s string) ParserFn {
	return Trans(fn, setStr(s))
}
