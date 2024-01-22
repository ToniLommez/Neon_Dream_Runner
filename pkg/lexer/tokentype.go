package lexer

type TokenType int

const (
	SPACE TokenType = iota

	// Context
	LEFT_PAREN    // (
	RIGHT_PAREN   // )
	LEFT_BRACE    // {
	RIGHT_BRACE   // }
	LEFT_DIAMOND  // <
	RIGHT_DIAMOND // >
	PIPE          // |
	DOT           // .
	COMMA         // ,
	AT            // @
	TAG           // #

	// Operators
	MINNUS         // -
	PLUS           // +
	SLASH          // /
	STAR           // *
	BANG           // !
	CHECK          // ?
	COLON          // :
	SEMICOLON      // ;
	RANGE_DOT      // ..
	INCREMENT      // ++
	DECREMENT      // --
	ELVIS          // ?:
	CHECK_NAV      // ?.
	BANG_NAV       // !.
	QUOTE          // '
	PIPELINE_RIGHT // |>
	PIPELINE_LEFT  // <|
	GO_IN          // <!
	GO_OUT         // !>
	GO_BI          // <!>

	// Assign
	ASSIGN                  // =
	ADD_ASSIGN              // +=
	SUB_ASSIGN              // -=
	MUL_ASSIGN              // *=
	DIV_ASSIGN              // /=
	MOD_ASSIGN              // &=
	BITSHIFT_LEFT_ASSIGN    // <<=
	BITSHIFT_RIGHT_ASSIGN   // >>=
	ROUNDSHIFT_LEFT_ASSIGN  // <<<=
	ROUNDSHIFT_RIGHT_ASSIGN // >>>=
	NOT_ASSIGN              // !=
	AND_ASSIGN              // &=
	NAND_ASSIGN             // ~&=
	OR_ASSIGN               // |=
	NOR_ASSIGN              // ~|=
	XOR_ASSIGN              // ^=
	XNOR_ASSIGN             // ~^=

	// Literals
	IDENTIFIER // function_name
	STRING     // "abc"
	CHAR       // 'a'
	NUMBER     // 12
	COMMENT    // /* */ //

	// Keywords
	LET            // let
	LET_BANG       // let!
	LET_CHECK      // let?
	LET_BANG_CHECK // let!?
	FN             // fn
	ASM            // asm
	FOR            // for
	LOOP           // loop
	WHILE          // while
	UNTIL          // until
	DO             // do
	IN             // in
	PULSE          // pulse
	BEFORE         // before
	INSIDE         // inside
	AFTER          // after
	RETURN         // =>]
	ERROR          // error
	NIL            // nil
	CASE           // case
	OF             // of
	IF             // if
	ELSE           // else
	ELIF           // elif
	USE            // use
	AS             // as
	MERGE          // merge
	OBJ            // obj
	PUB            // pub
	WHEN           // when
	TRIGGER        // trigger
	TRAIT          // trait
	THIS           // this

	// Reserved
	PRINT
	TRUE
	FALSE

	// Types
	BOOL
	INT
	FLOAT
	ANY

	// Special
	EOF
)

func (t TokenType) String() string {
	switch t {
	case SPACE:
		return "SPACE"
	case LEFT_PAREN:
		return "LEFT_PAREN"
	case RIGHT_PAREN:
		return "RIGHT_PAREN"
	case LEFT_BRACE:
		return "LEFT_BRACE"
	case RIGHT_BRACE:
		return "RIGHT_BRACE"
	case LEFT_DIAMOND:
		return "LEFT_DIAMOND"
	case RIGHT_DIAMOND:
		return "RIGHT_DIAMOND"
	case PIPE:
		return "PIPE"
	case DOT:
		return "DOT"
	case COMMA:
		return "COMMA"
	case AT:
		return "AT"
	case TAG:
		return "TAG"
	case MINNUS:
		return "MINNUS"
	case PLUS:
		return "PLUS"
	case SLASH:
		return "SLASH"
	case STAR:
		return "STAR"
	case BANG:
		return "BANG"
	case CHECK:
		return "CHECK"
	case COLON:
		return "COLON"
	case SEMICOLON:
		return "SEMICOLON"
	case RANGE_DOT:
		return "RANGE_DOT"
	case INCREMENT:
		return "INCREMENT"
	case DECREMENT:
		return "DECREMENT"
	case ELVIS:
		return "ELVIS"
	case CHECK_NAV:
		return "CHECK_NAV"
	case BANG_NAV:
		return "BANG_NAV"
	case QUOTE:
		return "QUOTE"
	case PIPELINE_RIGHT:
		return "PIPELINE_RIGHT"
	case PIPELINE_LEFT:
		return "PIPELINE_LEFT"
	case GO_IN:
		return "GO_IN"
	case GO_OUT:
		return "GO_OUT"
	case GO_BI:
		return "GO_BI"
	case ASSIGN:
		return "ASSIGN"
	case ADD_ASSIGN:
		return "ADD_ASSIGN"
	case SUB_ASSIGN:
		return "SUB_ASSIGN"
	case MUL_ASSIGN:
		return "MUL_ASSIGN"
	case DIV_ASSIGN:
		return "DIV_ASSIGN"
	case MOD_ASSIGN:
		return "MOD_ASSIGN"
	case BITSHIFT_LEFT_ASSIGN:
		return "BITSHIFT_LEFT_ASSIGN"
	case BITSHIFT_RIGHT_ASSIGN:
		return "BITSHIFT_RIGHT_ASSIGN"
	case ROUNDSHIFT_LEFT_ASSIGN:
		return "ROUNDSHIFT_LEFT_ASSIGN"
	case ROUNDSHIFT_RIGHT_ASSIGN:
		return "ROUNDSHIFT_RIGHT_ASSIGN"
	case NOT_ASSIGN:
		return "NOT_ASSIGN"
	case AND_ASSIGN:
		return "AND_ASSIGN"
	case NAND_ASSIGN:
		return "NAND_ASSIGN"
	case OR_ASSIGN:
		return "OR_ASSIGN"
	case NOR_ASSIGN:
		return "NOR_ASSIGN"
	case XOR_ASSIGN:
		return "XOR_ASSIGN"
	case XNOR_ASSIGN:
		return "XNOR_ASSIGN"
	case IDENTIFIER:
		return "IDENTIFIER"
	case STRING:
		return "STRING"
	case CHAR:
		return "CHAR"
	case NUMBER:
		return "NUMBER"
	case COMMENT:
		return "COMMENT"
	case LET:
		return "LET"
	case LET_BANG:
		return "LET_BANG"
	case LET_CHECK:
		return "LET_CHECK"
	case LET_BANG_CHECK:
		return "LET_BANG_CHECK"
	case FN:
		return "FN"
	case ASM:
		return "ASM"
	case FOR:
		return "FOR"
	case LOOP:
		return "LOOP"
	case WHILE:
		return "WHILE"
	case UNTIL:
		return "UNTIL"
	case DO:
		return "DO"
	case IN:
		return "IN"
	case PULSE:
		return "PULSE"
	case BEFORE:
		return "BEFORE"
	case INSIDE:
		return "INSIDE"
	case AFTER:
		return "AFTER"
	case RETURN:
		return "RETURN"
	case ERROR:
		return "ERROR"
	case NIL:
		return "NIL"
	case CASE:
		return "CASE"
	case OF:
		return "OF"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case ELIF:
		return "ELIF"
	case USE:
		return "USE"
	case AS:
		return "AS"
	case MERGE:
		return "MERGE"
	case OBJ:
		return "OBJ"
	case PUB:
		return "PUB"
	case WHEN:
		return "WHEN"
	case TRIGGER:
		return "TRIGGER"
	case TRAIT:
		return "TRAIT"
	case THIS:
		return "THIS"
	case PRINT:
		return "PRINT"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case BOOL:
		return "BOOL"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case ANY:
		return "ANY"
	case EOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}
