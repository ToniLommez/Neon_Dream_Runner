package lexer

type TokenType string

const (
	NEW_LINE TokenType = "NEW_LINE"

	// Context
	LEFT_PAREN  TokenType = "LEFT_PAREN"  // (
	RIGHT_PAREN TokenType = "RIGHT_PAREN" // )
	LEFT_BRACE  TokenType = "LEFT_BRACE"  // {
	RIGHT_BRACE TokenType = "RIGHT_BRACE" // }
	PIPE        TokenType = "PIPE"        // |
	DOT         TokenType = "DOT"         // .
	COMMA       TokenType = "COMMA"       // ,
	AT          TokenType = "AT"          // @
	TAG         TokenType = "TAG"         // #

	// Operators
	MINUS            TokenType = "MINUS"            // -
	PLUS             TokenType = "PLUS"             // +
	SLASH            TokenType = "SLASH"            // /
	STAR             TokenType = "STAR"             // *
	POW              TokenType = "POW"              // *
	MOD              TokenType = "MOD"              // %
	AND_BITWISE      TokenType = "AND_BITWISE"      // &
	OR_BITWISE       TokenType = "OR_BITWISE"       // |
	XOR_BITWISE      TokenType = "XOR_BITWISE"      // ^
	NAND_BITWISE     TokenType = "NAND_BITWISE"     // ~&
	NOR_BITWISE      TokenType = "NOR_BITWISE"      // ~|
	XNOR_BITWISE     TokenType = "XNOR_BITWISE"     // ~^
	EQUAL            TokenType = "EQUAL"            // ==
	NOT_EQUAL        TokenType = "NOT_EQUAL"        // !=
	GREATER_EQUAL    TokenType = "GREATER_EQUAL"    // >=
	LESS_EQUAL       TokenType = "LESS_EQUAL"       // <=
	AND_LOGIC        TokenType = "AND_LOGIC"        // &&
	OR_LOGIC         TokenType = "OR_LOGIC"         // ||
	NOT_BITWISE      TokenType = "NOT_BITWISE"      // ~
	LESS             TokenType = "LESS"             // <
	GREATER          TokenType = "GREATER"          // >
	SHIFT_LEFT       TokenType = "SHIFT_LEFT"       // <<
	SHIFT_RIGHT      TokenType = "SHIFT_RIGHT"      // >>
	ROUNDSHIFT_LEFT  TokenType = "ROUNDSHIFT_LEFT"  // <<<
	ROUNDSHIFT_RIGHT TokenType = "ROUNDSHIFT_RIGHT" // >>>
	BANG             TokenType = "BANG"             // !
	CHECK            TokenType = "CHECK"            // ?
	COLON            TokenType = "COLON"            // :
	SEMICOLON        TokenType = "SEMICOLON"        // ;
	RANGE_DOT        TokenType = "RANGE_DOT"        // ..
	INCREMENT        TokenType = "INCREMENT"        // ++
	DECREMENT        TokenType = "DECREMENT"        // --
	ELVIS            TokenType = "ELVIS"            // ?:
	CHECK_NAV        TokenType = "CHECK_NAV"        // ?.
	BANG_NAV         TokenType = "BANG_NAV"         // !.
	QUOTE            TokenType = "QUOTE"            // '
	PIPELINE_RIGHT   TokenType = "PIPELINE_RIGHT"   // |>
	PIPELINE_LEFT    TokenType = "PIPELINE_LEFT"    // <|
	GO_IN            TokenType = "GO_IN"            // <!
	GO_OUT           TokenType = "GO_OUT"           // !>
	GO_BI            TokenType = "GO_BI"            // <!>
	RETURN           TokenType = "RETURN"           // =>

	// Assign
	ASSIGN                  TokenType = "ASSIGN"                  // =
	ADD_ASSIGN              TokenType = "ADD_ASSIGN"              // +=
	SUB_ASSIGN              TokenType = "SUB_ASSIGN"              // -=
	MUL_ASSIGN              TokenType = "MUL_ASSIGN"              // *=
	POW_ASSIGN              TokenType = "POW_ASSIGN"              // **=
	DIV_ASSIGN              TokenType = "DIV_ASSIGN"              // /=
	MOD_ASSIGN              TokenType = "MOD_ASSIGN"              // %=
	BITSHIFT_LEFT_ASSIGN    TokenType = "BITSHIFT_LEFT_ASSIGN"    // <<=
	BITSHIFT_RIGHT_ASSIGN   TokenType = "BITSHIFT_RIGHT_ASSIGN"   // >>=
	ROUNDSHIFT_LEFT_ASSIGN  TokenType = "ROUNDSHIFT_LEFT_ASSIGN"  // <<<=
	ROUNDSHIFT_RIGHT_ASSIGN TokenType = "ROUNDSHIFT_RIGHT_ASSIGN" // >>>=
	AND_ASSIGN              TokenType = "AND_ASSIGN"              // &=
	NAND_ASSIGN             TokenType = "NAND_ASSIGN"             // ~&=
	OR_ASSIGN               TokenType = "OR_ASSIGN"               // |=
	NOR_ASSIGN              TokenType = "NOR_ASSIGN"              // ~|=
	XOR_ASSIGN              TokenType = "XOR_ASSIGN"              // ^=
	XNOR_ASSIGN             TokenType = "XNOR_ASSIGN"             // ~^=
	NOT_ASSIGN              TokenType = "NOT_ASSIGN"              // ~=

	// Literals
	IDENTIFIER     TokenType = "IDENTIFIER"     // function_name
	STRING_LITERAL TokenType = "STRING_LITERAL" // "abc"
	NUMBER_LITERAL TokenType = "NUMBER_LITERAL" // 12
	FLOAT_LITERAL  TokenType = "FLOAT_LITERAL"  // 12.3

	LET            TokenType = "LET"     // let
	LET_BANG       TokenType = "LET!"    // let
	LET_CHECK      TokenType = "LET?"    // let
	LET_BANG_CHECK TokenType = "LET!?"   // let
	FN             TokenType = "FN"      // fn
	FN_BANG        TokenType = "FN!"     // fn!
	ASM            TokenType = "ASM"     // asm
	FOR            TokenType = "FOR"     // for
	LOOP           TokenType = "LOOP"    // loop
	WHILE          TokenType = "WHILE"   // while
	UNTIL          TokenType = "UNTIL"   // until
	DO             TokenType = "DO"      // do
	IN             TokenType = "IN"      // in
	PULSE          TokenType = "PULSE"   // pulse
	BEFORE         TokenType = "BEFORE"  // before
	INSIDE         TokenType = "INSIDE"  // inside
	AFTER          TokenType = "AFTER"   // after
	ERROR          TokenType = "ERROR"   // error
	NIL            TokenType = "NIL"     // nil
	CASE           TokenType = "CASE"    // case
	OF             TokenType = "OF"      // of
	IF             TokenType = "IF"      // if
	ELSE           TokenType = "ELSE"    // else
	ELIF           TokenType = "ELIF"    // elif
	USE            TokenType = "USE"     // use
	AS             TokenType = "AS"      // as
	MERGE          TokenType = "MERGE"   // merge
	OBJ            TokenType = "OBJ"     // obj
	PUB            TokenType = "PUB"     // pub
	WHEN           TokenType = "WHEN"    // when
	TRIGGER        TokenType = "TRIGGER" // trigger
	TRAIT          TokenType = "TRAIT"   // trait
	THIS           TokenType = "THIS"    // this
	PRINT          TokenType = "PRINT"   // print
	PRINTF         TokenType = "PRINTF"  // printf
	PRINTLN        TokenType = "PRINTLN" // println
	TRUE           TokenType = "TRUE"    // true
	FALSE          TokenType = "FALSE"   // false
	INT            TokenType = "INT"     // int
	INT8           TokenType = "INT8"    // int8
	INT16          TokenType = "INT16"   // int16
	INT32          TokenType = "INT32"   // int32
	INT64          TokenType = "INT64"   // int64
	UINT           TokenType = "UINT"    // uint
	UINT8          TokenType = "UINT8"   // uint8
	UINT16         TokenType = "UINT16"  // uint16
	UINT32         TokenType = "UINT32"  // uint32
	UINT64         TokenType = "UINT64"  // uint64
	FLOAT          TokenType = "FLOAT"   // float
	FLOAT32        TokenType = "FLOAT32" // float32
	FLOAT64        TokenType = "FLOAT64" // float64
	BOOL           TokenType = "BOOL"    // bool
	CHAR           TokenType = "CHAR"    // char
	STRING         TokenType = "STRING"  // string
	BYTE           TokenType = "BYTE"    // byte
	ANY            TokenType = "ANY"     // any

	// Special
	EOF TokenType = "eof" // EOF
)
