package lexer

var keywords = map[string]TokenType{
	"let":     LET,
	"let!":    LET_BANG,
	"let?":    LET_CHECK,
	"let!?":   LET_BANG_CHECK,
	"fn":      FN,
	"fn!":     FN_BANG,
	"asm":     ASM,
	"for":     FOR,
	"loop":    LOOP,
	"while":   WHILE,
	"until":   UNTIL,
	"do":      DO,
	"in":      IN,
	"pulse":   PULSE,
	"before":  BEFORE,
	"inside":  INSIDE,
	"after":   AFTER,
	"error":   ERROR,
	"nil":     NIL,
	"case":    CASE,
	"of":      OF,
	"if":      IF,
	"else":    ELSE,
	"elif":    ELIF,
	"use":     USE,
	"as":      AS,
	"merge":   MERGE,
	"obj":     OBJ,
	"pub":     PUB,
	"when":    WHEN,
	"trigger": TRIGGER,
	"trait":   TRAIT,
	"this":    THIS,

	"put":     PUT,
	"print":   PRINT,
	"printf":  PRINTF,
	"println": PRINTLN,

	"true":  TRUE,
	"false": FALSE,

	"int":    INT,
	"i8":     I8,
	"i16":    I16,
	"i32":    I32,
	"i64":    I64,
	"uint":   UINT,
	"u8":     U8,
	"u16":    U16,
	"u32":    U32,
	"u64":    U64,
	"float":  FLOAT,
	"f32":    F32,
	"f64":    F64,
	"bool":   BOOL,
	"char":   CHAR,
	"string": STRING,
	"byte":   BYTE,
	"any":    ANY,
}

func (t TokenType) IsType() bool {
	return t == INT || t == I8 || t == I16 || t == I32 || t == I64 || t == UINT || t == U8 || t == U16 || t == U32 || t == U64 || t == FLOAT || t == F32 || t == F64 || t == BOOL || t == CHAR || t == STRING || t == BYTE || t == ANY
}

func (t TokenType) IsValidType() bool {
	return t == INT || t == UINT || t == FLOAT || t == BOOL || t == CHAR || t == STRING || t == UNDEFINED
}
