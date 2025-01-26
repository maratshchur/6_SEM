package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

type Token int

const (
	IF Token = iota
	ELSE
	PRINT
	INTEGER
	REAL
	CHARACTER
	END
	THEN
	PROGRAM
	USE
	MODULE
	TYPE
	CALL
	IMPLICIT
	NONE
	FUNCTION
	SUBROUTINE
	DO
	WHILE
	SELECT
	CASE
	IDENT
	ICONST
	RCONST
	SCONST
	PLUS
	MINUS
	MULT
	DIV
	ASSOP    // =
	EQ       // ==
	POW      // **
	GTHAN    // >
	LTHAN    // <
	GEQ      // >=
	LEQ      // <=
	NEQ      // /=
	CAT      // //
	COMMA    // ,
	LPAREN   // (
	RPAREN   // )
	LBRACKET // [
	RBRACKET // ]
	DOT      // .
	DCOLON   // ::
	PERCENT  // %
	ERR
	DONE
)

var tokenMap = map[Token]string{
	IF:        "IF",
	ELSE:      "ELSE",
	PRINT:     "PRINT",
	INTEGER:   "INTEGER",
	REAL:      "REAL",
	CHARACTER: "CHARACTER",
	END:       "END",
	THEN:      "THEN",
	PROGRAM:   "PROGRAM",
	USE:       "USE",
	MODULE:    "MODULE",
	TYPE:      "TYPE",
	CALL:      "CALL",
	IMPLICIT:  "IMPLICIT",
	NONE:      "NONE",
	FUNCTION:  "FUNCTION",
	SUBROUTINE: "SUBROUTINE",
	DO:        "DO",
	WHILE:     "WHILE",
	SELECT:    "SELECT",
	CASE:      "CASE",
	IDENT:     "IDENT",
	ICONST:    "ICONST",
	RCONST:    "RCONST",
	SCONST:    "SCONST",
	PLUS:      "PLUS",
	MINUS:     "MINUS",
	MULT:      "MULT",
	DIV:       "DIV",
	ASSOP:     "ASSOP",
	EQ:        "EQ",
	POW:       "POW",
	GTHAN:     "GTHAN",
	LTHAN:     "LTHAN",
	GEQ:       "GEQ",
	LEQ:       "LEQ",
	NEQ:       "NEQ",
	CAT:       "CAT",
	COMMA:     "COMMA",
	LPAREN:    "LPAREN",
	RPAREN:    "RPAREN",
	LBRACKET:  "LBRACKET",
	RBRACKET:  "RBRACKET",
	DOT:       "DOT",
	DCOLON:    "DCOLON",
	PERCENT:   "PERCENT",
	ERR:       "ERR",
	DONE:      "DONE",
}

type LexItem struct {
	token  Token
	lexeme string
	line   int
	column int
	id     int
}

type Lexer struct {
	input      string
	position   int
	line       int
	column     int
	lexeme     string
	stack      []rune
	constID    int
	idMap      map[string]int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:   input,
		line:    1,
		column:  0,
		stack:   []rune{},
		idMap:   make(map[string]int),
	}
}

func (lx *Lexer) NextToken() LexItem {
	state := "START"
	for lx.position < len(lx.input) {
		ch := lx.input[lx.position]
		lx.position++
		lx.column++

		switch state {
		case "START":
			if unicode.IsSpace(rune(ch)) {
				if ch == '\n' {
					lx.line++
					lx.column = 0
				}
				continue
			}
			if unicode.IsLetter(rune(ch)) || ch == '_' {
				lx.lexeme = string(ch)
				state = "INID"
			} else if unicode.IsDigit(rune(ch)) {
				lx.lexeme = string(ch)
				state = "ININT"
			} else if ch == '"' {
				lx.lexeme = ""
				lx.stack = append(lx.stack, '"')
				state = "INSTRING"
			} else if ch == '(' {
				lx.stack = append(lx.stack, '(')
				return LexItem{token: LPAREN, lexeme: "(", line: lx.line, column: lx.column, id: lx.nextID()}
			} else if ch == ')' {
				if len(lx.stack) > 0 && lx.stack[len(lx.stack)-1] == '(' {
					lx.stack = lx.stack[:len(lx.stack)-1]
				} else {
					return LexItem{
						token:  ERR,
						lexeme: fmt.Sprintf("Unmatched closing parenthesis: ')'", ch),
						line:   lx.line,
						column: lx.column,
						id:     lx.nextID(),
					}
				}
				return LexItem{token: RPAREN, lexeme: ")", line: lx.line, column: lx.column, id: lx.nextID()}
			} else if ch == '[' {
				lx.stack = append(lx.stack, '[')
				return LexItem{token: LBRACKET, lexeme: "[", line: lx.line, column: lx.column, id: lx.nextID()}
			} else if ch == ']' {
				if len(lx.stack) > 0 && lx.stack[len(lx.stack)-1] == '[' {
					lx.stack = lx.stack[:len(lx.stack)-1]
				} else {
					return LexItem{
						token:  ERR,
						lexeme: fmt.Sprintf("Unmatched closing bracket: ']'"),
						line:   lx.line,
						column: lx.column,
						id:     lx.nextID(),
					}
				}
				return LexItem{token: RBRACKET, lexeme: "]", line: lx.line, column: lx.column, id: lx.nextID()}
			} else if ch == '!' {
				state = "COMMENT"
			} else if strings.ContainsRune("+-*/=<>.,:()%$", rune(ch)) {
				lx.lexeme = string(ch)
				state = "SIGN"
			} else {
				return LexItem{
					token:  ERR,
					lexeme: fmt.Sprintf("Unexpected character: '%c'", ch),
					line:   lx.line,
					column: lx.column,
					id:     lx.nextID(),
				}
			}

		case "INID":
			if unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_' {
				lx.lexeme += string(ch)
			} else {
				lx.position--
				lx.column--
				return lx.identOrKeyword()
			}

		case "ININT":
			if unicode.IsDigit(rune(ch)) {
				lx.lexeme += string(ch)
			} else if ch == '.' {
				lx.lexeme += string(ch)
				state = "INREAL"
			} else {
				lx.position--
				lx.column--
				return lx.constantOrError(ICONST)
			}

		case "INREAL":
			if unicode.IsDigit(rune(ch)) || ch == 'E' || ch == '+' || ch == '-' {
				lx.lexeme += string(ch)
			} else {
				lx.position--
				lx.column--
				return lx.constantOrError(RCONST)
			}

		case "INSTRING":
			if ch == '"' {
				if len(lx.stack) == 0 || lx.stack[len(lx.stack)-1] != '"' {
					return LexItem{
						token:  ERR,
						lexeme: "Mismatched quotes",
						line:   lx.line,
						column: lx.column,
						id:     lx.nextID(),
					}
				}
				lx.stack = lx.stack[:len(lx.stack)-1]
				return LexItem{token: SCONST, lexeme: lx.lexeme, line: lx.line, column: lx.column, id: lx.nextID()}
			}
			lx.lexeme += string(ch)

		case "COMMENT":
			if ch == '\n' {
				lx.line++
				lx.column = 0
				state = "START"
			}

		case "SIGN":
			if lx.lexeme == ":" && ch == ':' {
				lx.lexeme += string(ch)
				return LexItem{token: DCOLON, lexeme: lx.lexeme, line: lx.line, column: lx.column, id: lx.nextID()}
			} else if lx.lexeme == "=" && ch == '=' {
				lx.lexeme += string(ch)
				return LexItem{token: EQ, lexeme: lx.lexeme, line: lx.line, column: lx.column, id: lx.nextID()}
			} else if lx.lexeme == "<" && ch == '=' {
				lx.lexeme += string(ch)
				return LexItem{token: LEQ, lexeme: lx.lexeme, line: lx.line, column: lx.column, id: lx.nextID()}
			} else if lx.lexeme == ">" && ch == '=' {
				lx.lexeme += string(ch)
				return LexItem{token: GEQ, lexeme: lx.lexeme, line: lx.line, column: lx.column, id: lx.nextID()}
			} else if lx.lexeme == "/" && ch == '=' {
				lx.lexeme += string(ch)
				return LexItem{token: NEQ, lexeme: lx.lexeme, line: lx.line, column: lx.column, id: lx.nextID()}
			} else {
				lx.position--
				lx.column--
				return lx.operatorOrError()
			}
		}
	}

	if len(lx.stack) > 0 {
		unmatched := lx.stack[len(lx.stack)-1]
		return LexItem{
			token:  ERR,
			lexeme: fmt.Sprintf("Unmatched delimiter: '%c'", unmatched),
			line:   lx.line,
			column: lx.column,
			id:     lx.nextID(),
		}
	}

	return LexItem{token: DONE, lexeme: "", line: lx.line, column: lx.column, id: lx.nextID()}
}

func (lx *Lexer) identOrKeyword() LexItem {
	isValid, errMsg := isValidIdentifier(lx.lexeme)
	if !isValid {
		return LexItem{
			token:  ERR,
			lexeme: fmt.Sprintf("Invalid identifier: %s (%s)", lx.lexeme, errMsg),
			line:   lx.line,
			column: lx.column,
			id:     lx.nextID(),
		}
	}

	keywords := map[string]Token{
		"IF":        IF,
		"ELSE":      ELSE,
		"PRINT":     PRINT,
		"INTEGER":   INTEGER,
		"REAL":      REAL,
		"CHARACTER": CHARACTER,
		"END":       END,
		"THEN":      THEN,
		"PROGRAM":   PROGRAM,
		"USE":       USE,
		"MODULE":    MODULE,
		"TYPE":      TYPE,
		"CANCEL":    CALL,
		"IMPLICIT":  IMPLICIT,
		"NONE":      NONE,
		"FUNCTION":  FUNCTION,
		"SUBROUTINE": SUBROUTINE,
		"DO":        DO,
		"WHILE":     WHILE,
		"SELECT":    SELECT,
		"CASE":      CASE,
	}
	token, found := keywords[strings.ToUpper(lx.lexeme)]
	if found {
		return LexItem{
			token:  token,
			lexeme: lx.lexeme,
			line:   lx.line,
			column: lx.column,
			id:     lx.nextID(),
		}
	}

	return lx.assignID(IDENT, lx.lexeme)
}

func (lx *Lexer) constantOrError(tokenType Token) LexItem {
	if tokenType == ICONST {
		return lx.assignID(ICONST, lx.lexeme)
	}
	return lx.assignID(RCONST, lx.lexeme)
}

func (lx *Lexer) assignID(token Token, lexeme string) LexItem {
	if id, found := lx.idMap[lexeme]; found {
		return LexItem{token: token, lexeme: lexeme, line: lx.line, column: lx.column, id: id}
	}
	lx.constID++
	lx.idMap[lexeme] = lx.constID
	return LexItem{token: token, lexeme: lexeme, line: lx.line, column: lx.column, id: lx.constID}
}

func (lx *Lexer) operatorOrError() LexItem {
	operators := map[string]Token{
		"+":  PLUS,
		"-":  MINUS,
		"*":  MULT,
		"/":  DIV,
		"=":  ASSOP,
		"<":  LTHAN,
		">":  GTHAN,
		",":  COMMA,
		".":  DOT,
		"[":  LBRACKET,
		"]":  RBRACKET,
		"%":  PERCENT,
		"::": DCOLON,
	}
	token, found := operators[lx.lexeme]
	if found {
		return LexItem{token: token, lexeme: lx.lexeme, line: lx.line, column: lx.column, id: lx.nextID()}
	}
	return LexItem{
		token:  ERR,
		lexeme: fmt.Sprintf("Unknown operator: %s", lx.lexeme),
		line:   lx.line,
		column: lx.column,
		id:     lx.nextID(),
	}
}

func (lx *Lexer) nextID() int {
	lx.constID++
	return lx.constID
}

func isValidIdentifier(lexeme string) (bool, string) {
	if len(lexeme) == 0 {
		return false, "Identifier is empty"
	}

	if unicode.IsDigit(rune(lexeme[0])) {
		return false, "Identifier cannot start with a digit"
	}

	for _, ch := range lexeme {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			return false, fmt.Sprintf("Invalid character '%c' in identifier", ch)
		}
	}

	return true, "" 
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: lexer <input file> <output file>")
		return
	}

	inputFileName := os.Args[1]
	outputFileName := os.Args[2]

	inputFile, err := os.ReadFile(inputFileName)
	if err != nil {
		fmt.Printf("Error reading input file: %s\n", err)
		return
	}

	lexer := NewLexer(string(inputFile))
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Printf("Error creating output file: %s\n", err)
		return
	}
	defer outputFile.Close()

	for {
		token := lexer.NextToken()
		if token.token == DONE {
			break
		}
		if token.token == ERR {
			fmt.Fprintf(outputFile, "Error: %s at line %d, column %d, ID: %d\n", token.lexeme, token.line, token.column, token.id)
			return 
		}
		fmt.Fprintf(outputFile, "Token: %-10s Lexeme: %-10s Line: %d, Column: %d, ID: %d\n", tokenMap[token.token], token.lexeme, token.line, token.column, token.id)
	}
}