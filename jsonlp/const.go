package jsonlp

import (
	"math"

	. "github.com/shellyln/takenoco/base"
)

// Parser constant value
var nilAst = Ast{
	OpCode:    0,
	ClassName: "Null",
	Type:      AstType_Nil,
	Value:     nil,
}

// Parser constant value
var trueAst = Ast{
	OpCode:    0,
	ClassName: "Bool",
	Type:      AstType_Bool,
	Value:     true,
}

// Parser constant value
var falseAst = Ast{
	OpCode:    0,
	ClassName: "Bool",
	Type:      AstType_Bool,
	Value:     false,
}

// Parser constant value
// var zeroStrAst = Ast{
// 	OpCode:    0,
// 	ClassName: "String",
// 	Type:      AstType_String,
// 	Value:     "",
// }

// Parser constant value
var positiveInfinityAst = Ast{
	OpCode:    0,
	ClassName: "Float",
	Type:      AstType_Float,
	Value:     math.Inf(1),
}

// Parser constant value
var negativeInfinityAst = Ast{
	OpCode:    0,
	ClassName: "Float",
	Type:      AstType_Float,
	Value:     math.Inf(-1),
}

// Parser constant value
var nanAst = Ast{
	OpCode:    0,
	ClassName: "Float",
	Type:      AstType_Float,
	Value:     math.NaN(),
}
