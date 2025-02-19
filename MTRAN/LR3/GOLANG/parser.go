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
	CONTAINS
	COLON
	DEFAULT
)

var tokenMap = map[Token]string{
	IF:         "IF",
	ELSE:       "ELSE",
	PRINT:      "PRINT",
	INTEGER:    "INTEGER",
	REAL:       "REAL",
	CHARACTER:  "CHARACTER",
	END:        "END",
	THEN:       "THEN",
	PROGRAM:    "PROGRAM",
	USE:        "USE",
	MODULE:     "MODULE",
	TYPE:       "TYPE",
	CALL:       "CALL",
	IMPLICIT:   "IMPLICIT",
	NONE:       "NONE",
	FUNCTION:   "FUNCTION",
	SUBROUTINE: "SUBROUTINE",
	DO:         "DO",
	WHILE:      "WHILE",
	SELECT:     "SELECT",
	CASE:       "CASE",
	IDENT:      "IDENT",
	ICONST:     "ICONST",
	RCONST:     "RCONST",
	SCONST:     "SCONST",
	PLUS:       "PLUS",
	MINUS:      "MINUS",
	MULT:       "MULT",
	DIV:        "DIV",
	ASSOP:      "ASSOP",
	EQ:         "EQ",
	POW:        "POW",
	GTHAN:      "GTHAN",
	LTHAN:      "LTHAN",
	GEQ:        "GEQ",
	LEQ:        "LEQ",
	NEQ:        "NEQ",
	CAT:        "CAT",
	COMMA:      "COMMA",
	LPAREN:     "LPAREN",
	RPAREN:     "RPAREN",
	LBRACKET:   "LBRACKET",
	RBRACKET:   "RBRACKET",
	DOT:        "DOT",
	DCOLON:     "DCOLON",
	COLON:      "COLON",
	PERCENT:    "PERCENT",
	ERR:        "ERR",
	DONE:       "DONE",
	CONTAINS:   "CONTAINS",
	DEFAULT:    "DEFAULT",
}

type LexItem struct {
	token  Token
	lexeme string
	line   int
	column int
	id     int
}

type Lexer struct {
	input    string
	position int
	line     int
	column   int
	lexeme   string
	stack    []rune
	constID  int
	idMap    map[string]int
	tokens   []LexItem
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		line:   1,
		column: 0,
		stack:  []rune{},
		idMap:  make(map[string]int),
		tokens: []LexItem{},
	}
}

// GenerateTokens - Генерация всех токенов
func (lx *Lexer) GenerateTokens() {
	for {
		token := lx.NextToken()
		lx.tokens = append(lx.tokens, token) // Сохраняем токен в срез
		if token.token == DONE || token.token == ERR {
			break
		}
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
				return lx.assignID(LPAREN, "(", lx.line, lx.column)
			} else if ch == ')' {
				if len(lx.stack) > 0 && lx.stack[len(lx.stack)-1] == '(' {
					lx.stack = lx.stack[:len(lx.stack)-1]
				} else {
					return lx.errorLexeme(fmt.Sprintf("Unmatched closing parenthesis: ')'", ch), lx.line, lx.column)
				}
				return lx.assignID(RPAREN, ")", lx.line, lx.column)
			} else if ch == '[' {
				lx.stack = append(lx.stack, '[')
				return lx.assignID(LBRACKET, "[", lx.line, lx.column)
			} else if ch == ']' {
				if len(lx.stack) > 0 && lx.stack[len(lx.stack)-1] == '[' {
					lx.stack = lx.stack[:len(lx.stack)-1]
				} else {
					return lx.errorLexeme("Unmatched closing bracket: ']' at line", lx.line, lx.column)
				}
				return lx.assignID(RBRACKET, "]", lx.line, lx.column)
			} else if ch == '!' {
				state = "COMMENT"
			} else if strings.ContainsRune("+-*/=<>.,:()%$", rune(ch)) {
				lx.lexeme = string(ch)
				state = "SIGN"
			} else {
				return lx.errorLexeme(fmt.Sprintf("Unexpected character: '%c'", ch), lx.line, lx.column)
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
				return lx.assignID(ICONST, lx.lexeme, lx.line, lx.column)
			}

		case "INREAL":
			if unicode.IsDigit(rune(ch)) || ch == 'E' || ch == '+' || ch == '-' {
				lx.lexeme += string(ch)
			} else {
				lx.position--
				lx.column--
				return lx.assignID(RCONST, lx.lexeme, lx.line, lx.column)
			}

		case "INSTRING":
			if ch == '"' {
				if len(lx.stack) == 0 || lx.stack[len(lx.stack)-1] != '"' {
					return lx.errorLexeme("Mismatched quotes", lx.line, lx.column)
				}
				lx.stack = lx.stack[:len(lx.stack)-1]
				return lx.assignID(SCONST, lx.lexeme, lx.line, lx.column)
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
				return lx.assignID(DCOLON, lx.lexeme, lx.line, lx.column)
			} else if lx.lexeme == "=" && ch == '=' {
				lx.lexeme += string(ch)
				return lx.assignID(EQ, lx.lexeme, lx.line, lx.column)
			} else if lx.lexeme == "<" && ch == '=' {
				lx.lexeme += string(ch)
				return lx.assignID(LEQ, lx.lexeme, lx.line, lx.column)
			} else if lx.lexeme == ">" && ch == '=' {
				lx.lexeme += string(ch)
				return lx.assignID(GEQ, lx.lexeme, lx.line, lx.column)
			} else if lx.lexeme == "/" && ch == '=' {
				lx.lexeme += string(ch)
				return lx.assignID(NEQ, lx.lexeme, lx.line, lx.column)
			} else {
				lx.position--
				lx.column--
				return lx.operatorOrError()
			}
		}
	}

	if len(lx.stack) > 0 {
		unmatched := lx.stack[len(lx.stack)-1]
		return lx.errorLexeme(fmt.Sprintf("Unmatched delimiter: '%c'", unmatched), lx.line, lx.column)
	}

	return LexItem{token: DONE, lexeme: "", line: lx.line, column: lx.column, id: lx.nextID()}
}

func (lx *Lexer) identOrKeyword() LexItem {
	isValid, errMsg := isValidIdentifier(lx.lexeme)
	if !isValid {
		return lx.errorLexeme(fmt.Sprintf("Invalid identifier: %s (%s)", lx.lexeme, errMsg), lx.line, lx.column)
	}

	keywords := map[string]Token{
		"IF":         IF,
		"ELSE":       ELSE,
		"PRINT":      PRINT,
		"INTEGER":    INTEGER,
		"REAL":       REAL,
		"CHARACTER":  CHARACTER,
		"END":        END,
		"THEN":       THEN,
		"PROGRAM":    PROGRAM,
		"USE":        USE,
		"MODULE":     MODULE,
		"TYPE":       TYPE,
		"CANCEL":     CALL,
		"IMPLICIT":   IMPLICIT,
		"NONE":       NONE,
		"FUNCTION":   FUNCTION,
		"SUBROUTINE": SUBROUTINE,
		"DO":         DO,
		"WHILE":      WHILE,
		"SELECT":     SELECT,
		"CASE":       CASE,
		"CONTAINS":   CONTAINS,
		"DEFAULT":    DEFAULT,
	}
	token, found := keywords[strings.ToUpper(lx.lexeme)]
	if found {
		return lx.assignID(token, lx.lexeme, lx.line, lx.column)
	}

	return lx.assignID(IDENT, lx.lexeme, lx.line, lx.column)
}

func (lx *Lexer) assignID(token Token, lexeme string, line, column int) LexItem {
	if id, found := lx.idMap[lexeme]; found {
		return LexItem{token: token, lexeme: lexeme, line: line, column: column, id: id}
	}
	lx.constID++
	lx.idMap[lexeme] = lx.constID
	return LexItem{token: token, lexeme: lexeme, line: line, column: column, id: lx.constID}
}

func (lx *Lexer) errorLexeme(message string, line, column int) LexItem {
	return LexItem{
		token:  ERR,
		lexeme: message,
		line:   line,
		column: column,
		id:     lx.nextID(),
	}
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
		":":  COLON,
	}
	token, found := operators[lx.lexeme]
	if found {
		return lx.assignID(token, lx.lexeme, lx.line, lx.column)
	}
	return lx.errorLexeme(fmt.Sprintf("Unknown operator: %s", lx.lexeme), lx.line, lx.column)
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

// ASTNode - Узел синтаксического дерева
type ASTNode struct {
	NodeType string
	Value    string
	Children []*ASTNode
}

// Создание нового узла AST
func NewASTNode(nodeType, value string) *ASTNode {
	return &ASTNode{
		NodeType: nodeType,
		Value:    value,
		Children: []*ASTNode{},
	}
}

// Добавление дочернего узла
func (n *ASTNode) AddChild(child *ASTNode) {
	n.Children = append(n.Children, child)
}

// PrintTree - Красивый вывод AST в виде дерева
func (n *ASTNode) PrintTree(prefix string, isLast bool) {
	fmt.Print(prefix)

	if isLast {
		fmt.Print("└─")
		prefix += "  "
	} else {
		fmt.Print("├─")
		prefix += "│ "
	}

	fmt.Println(n.NodeType, "-", n.Value)

	for i, child := range n.Children {
		child.PrintTree(prefix, i == len(n.Children)-1)
	}
}

// Parser - Парсер для анализа последовательности токенов
type Parser struct {
	tokens []LexItem // Массив токенов
	curTok LexItem   // Текущий токен
	pos    int       // Текущая позиция в массиве токенов
}

// Создание парсера
func NewParser(tokens []LexItem) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    -1, // Начинаем до первого токена
	}
}

// Получение следующего токена
func (p *Parser) nextToken() {
	p.pos++
	if p.pos < len(p.tokens) {
		if p.curTok.token == ERR {

			p.curTok = LexItem{token: DONE, lexeme: p.curTok.lexeme, line: p.curTok.line, column: p.curTok.column, id: p.curTok.id}
		}
		p.curTok = p.tokens[p.pos]
	} else {
		p.curTok = LexItem{token: DONE, lexeme: "", line: 0, column: 0, id: 0}
	}
}

// Главная функция парсера
func (p *Parser) Parse() *ASTNode {
	for _, item := range p.tokens {
		if item.token == ERR {
			fmt.Printf("Error: %s at line %d, column %d, ID: %d\n", item.lexeme, item.line, item.column, item.id)
			return nil

		}
	}

	p.nextToken() // Получаем первый токен
	return p.parseProgram()
}

// Парсинг программы
func (p *Parser) parseProgram() *ASTNode {
	root := NewASTNode("Program", "")

	for p.curTok.token != DONE {
		switch p.curTok.token {
		case MODULE:
			root.AddChild(p.parseModule())
		case FUNCTION:
			root.AddChild(p.parseFunction())
		case SUBROUTINE:
			root.AddChild(p.parseSubroutine())
		case TYPE:
			root.AddChild(p.parseTypeDefinition())
		case IF:
			root.AddChild(p.parseIfStatement())
		case DO:
			root.AddChild(p.parseDoLoop())
		case WHILE:
			root.AddChild(p.parseWhileLoop())
		case INTEGER, REAL, CHARACTER:
			root.AddChild(p.parseVariableDeclaration())
		case IDENT:
			root.AddChild(p.parseAssignmentOrExpression())
		case ERR:
			root.AddChild(NewASTNode(string(p.curTok.token), p.curTok.lexeme))
			break
		default:
			p.nextToken()
		}
	}

	return root
}
func (p *Parser) parseModule() *ASTNode {
	moduleNode := NewASTNode("Module", p.curTok.lexeme)
	p.nextToken() // Пропускаем "MODULE"

	if p.curTok.token == IDENT {
		moduleNode.Value = p.curTok.lexeme
		p.nextToken()
	}

	for p.curTok.token != END && p.curTok.token != DONE {
		switch p.curTok.token {
		case IMPLICIT:
			moduleNode.AddChild(p.parseImplicit())
		case TYPE:
			moduleNode.AddChild(p.parseTypeDefinition())
		case FUNCTION:
			moduleNode.AddChild(p.parseFunction())
		case SUBROUTINE:
			moduleNode.AddChild(p.parseSubroutine())
		case CONTAINS:
			p.nextToken()
			containsNode := NewASTNode("Contains", "Functions declaration")
			for p.curTok.token != END && p.curTok.token != DONE {
				if p.curTok.token == FUNCTION {
					containsNode.AddChild(p.parseFunction())
				} else if p.curTok.token == SUBROUTINE {
					containsNode.AddChild(p.parseSubroutine())
				} else {
					p.nextToken()
				}
			}
			moduleNode.AddChild(containsNode)
		default:
			p.nextToken()
		}
	}

	p.nextToken() // Пропускаем "END MODULE"
	p.nextToken() // Пропускаем "END MODULE"
	p.nextToken() // Пропускаем "END MODULE"

	return moduleNode
}

func (p *Parser) parseImplicit() *ASTNode {
	implicitNode := NewASTNode("Implicit", p.curTok.lexeme)
	p.nextToken()
	if p.curTok.token == NONE {
		implicitNode.Value = p.curTok.lexeme
		p.nextToken()
	}
	return implicitNode
}

func (p *Parser) parseFunction() *ASTNode {
	funcNode := NewASTNode("Function", "")
	p.nextToken()

	if p.curTok.token == IDENT {
		funcNode.Value = p.curTok.lexeme
		p.nextToken()
	}

	// Аргументы
	argsNode := NewASTNode("Arguments", "")
	if p.curTok.token == LPAREN {
		p.nextToken()
		for p.curTok.token != RPAREN && p.curTok.token != DONE {
			if p.curTok.token == IDENT {
				argNode := NewASTNode("Argument", p.curTok.lexeme)
				p.nextToken()
				if p.curTok.token == COMMA {
					p.nextToken()
				}
				argsNode.AddChild(argNode)
			}
		}
		p.nextToken() // Пропускаем ")"
	}
	funcNode.AddChild(argsNode)
	if p.curTok.lexeme == "result" {
		returnNodes := NewASTNode("Return Parameters", "")
		p.nextToken()
		p.nextToken()
		for p.curTok.token != RPAREN && p.curTok.token != DONE {
			if p.curTok.token == IDENT {
				returnNode := NewASTNode("Return parameter", p.curTok.lexeme)
				p.nextToken()
				if p.curTok.token == COMMA {
					p.nextToken()
				}
				returnNodes.AddChild(returnNode)
			}
		}
		p.nextToken() // Пропускаем ")"

		funcNode.AddChild(returnNodes)
	}

	// Тело функции
	bodyNode := NewASTNode("FunctionBody", "")
	for p.curTok.token != END && p.curTok.token != DONE {
		bodyNode.AddChild(p.parseStatement())
	}
	funcNode.AddChild(bodyNode)

	p.nextToken() // Пропускаем "END FUNCTION"
	p.nextToken() // Пропускаем "END FUNCTION"
	p.nextToken() // Пропускаем "END FUNCTION"

	return funcNode
}

func (p *Parser) parseSubroutine() *ASTNode {
	subNode := NewASTNode("Subroutine", "")

	p.nextToken()
	if p.curTok.token == IDENT {
		subNode.Value = p.curTok.lexeme
		p.nextToken()
	}

	// Аргументы
	argsNode := NewASTNode("Arguments", "")
	if p.curTok.token == LPAREN {
		p.nextToken()
		for p.curTok.token != RPAREN && p.curTok.token != DONE {
			if p.curTok.token == IDENT {
				argsNode.AddChild(NewASTNode("Argument", p.curTok.lexeme))
				p.nextToken()
			}
			if p.curTok.token == COMMA {
				p.nextToken()
			}
		}
		p.nextToken() // Пропускаем ")"
	}
	subNode.AddChild(argsNode)

	// Тело подпрограммы
	bodyNode := NewASTNode("SubroutineBody", "")
	for p.curTok.token != END && p.curTok.token != DONE && p.curTok.token != ERR {

		bodyNode.AddChild(p.parseStatement())
	}
	if p.curTok.token == ERR {
		fmt.Printf("Error: %s at line %d, column %d, ID: %d\n", p.curTok.lexeme, p.curTok.line, p.curTok.column, p.curTok.id)
		os.Exit(1)
	}
	subNode.AddChild(bodyNode)

	p.nextToken()
	p.nextToken()

	return subNode
}

func (p *Parser) parseAssignmentOrExpression() *ASTNode {
	identNode := NewASTNode("Identifier", p.curTok.lexeme)
	p.nextToken()

	if p.curTok.token == PERCENT { // Разыменование компонента типа person%name
		p.nextToken()
		if p.curTok.token == IDENT {
			fieldNode := NewASTNode("FieldAccess", p.curTok.lexeme)
			identNode.AddChild(fieldNode)
			p.nextToken()
			if p.curTok.token == ASSOP {
				p.nextToken()
				exprNode := p.parseExpression()
				assignNode := NewASTNode("Assignment", "")
				assignNode.AddChild(identNode)
				assignNode.AddChild(exprNode)
				return assignNode
			}
			return identNode
		}
	}

	if p.curTok.token == ASSOP {
		p.nextToken()
		exprNode := p.parseExpression()
		assignNode := NewASTNode("Assignment", "")
		assignNode.AddChild(identNode)
		assignNode.AddChild(exprNode)
		return assignNode
	}

	return identNode
}
func (p *Parser) parseDoWhileStatement() *ASTNode {
	doNode := NewASTNode("DoWhileStatement", "")
	p.nextToken()

	bodyNode := NewASTNode("DoBody", "")
	for p.curTok.token != WHILE && p.curTok.token != END && p.curTok.token != DONE {
		bodyNode.AddChild(p.parseStatement())
	}
	doNode.AddChild(bodyNode)

	if p.curTok.token == WHILE {
		p.nextToken()
		conditionNode := p.parseExpression()
		doNode.AddChild(conditionNode)
	}

	p.nextToken() // Пропускаем "END DO"
	p.nextToken() // Пропускаем "END DO"

	return doNode
}

func (p *Parser) parseDoLoop() *ASTNode {
	doNode := NewASTNode("DoLoop", "")
	p.nextToken()
	if p.curTok.token == WHILE {
		p.nextToken()
	}
	if p.curTok.token == LPAREN {
		p.nextToken()
	}
	// Инициализация счетчика
	conditionNode := NewASTNode("DoConditionStatement", "")
	initNode := p.parseAssignmentOrExpression()
	conditionNode.AddChild(initNode)
	doNode.AddChild(conditionNode)
	if p.curTok.token == COMMA {
		condition := NewASTNode("Condition", "<")
		conditionNode.AddChild(condition)
		p.nextToken()

		exprNode := p.parseExpression()
		conditionNode.AddChild(exprNode)
	} else {
		condition := NewASTNode("Condition", p.curTok.lexeme)
		conditionNode.AddChild(condition)
		p.nextToken()
		exprNode := p.parseExpression()
		conditionNode.AddChild(exprNode)
		p.nextToken()

	}
	bodyNode := NewASTNode("DoBody", "")
	for p.curTok.token != END && p.curTok.token != DONE {
		bodyNode.AddChild(p.parseStatement())
	}
	doNode.AddChild(bodyNode)

	p.nextToken() // Пропускаем "END DO"
	p.nextToken() // Пропускаем "END DO"

	return doNode
}

func (p *Parser) parsePrintStatement() *ASTNode {
	printNode := NewASTNode("Print", "")
	p.nextToken()

	argsNode := NewASTNode("PrintArguments", "")
	if p.curTok.token != MULT {
		fmt.Printf("Error: %s at line %d, %s, column %d, ID: %d\n", p.curTok.lexeme, p.curTok.line, "* expected", p.curTok.column, p.curTok.id)
		os.Exit(1)
	} else {
		p.nextToken()
		p.nextToken()
		lastToken := p.curTok.token // Keep track of the last token

		for p.curTok.token != RPAREN && p.curTok.token != DONE && p.curTok.token != END {
			argsNode.AddChild(p.parseExpression())

			if p.curTok.token == COMMA {
				p.nextToken()
				lastToken = p.curTok.token
			} else {
				if lastToken != COMMA && p.curTok.token != COMMA {
					break
				}
			}

			lastToken = p.curTok.token
		}
	}
	printNode.AddChild(argsNode)
	return printNode
}

func (p *Parser) parseTypeDefinition() *ASTNode {
	typeNode := NewASTNode("TypeDefinition", p.curTok.lexeme)
	p.nextToken() // Пропускаем "TYPE"
	if p.curTok.token == LPAREN {
		typeNode = p.parseTypeDeclaration()
		return typeNode
	} else {
		p.nextToken() // Пропускаем "::"

		if p.curTok.token == IDENT {
			typeNode.Value = p.curTok.lexeme
			p.nextToken()
		}

		for p.curTok.token != END && p.curTok.token != DONE {
			typeNode.AddChild(p.parseVariableDeclaration())
		}

		p.nextToken() // Пропускаем "END TYPE"
		p.nextToken() // Пропускаем "END TYPE"
		p.nextToken() // Пропускаем "END TYPE"
		return typeNode
	}
}

func (p *Parser) parseTypeDeclaration() *ASTNode {
	p.nextToken() // Пропускаем "TYPE"
	typeNode := NewASTNode("TypeDeclaration", p.curTok.lexeme)
	if p.curTok.token == IDENT {
		typeNode.Value = p.curTok.lexeme
		p.nextToken()
	}
	p.nextToken() // Пропускаем "TYPE
	p.nextToken() // Пропускаем "TYPE"

	typeName := NewASTNode("Variable", p.curTok.lexeme)
	typeNode.AddChild(typeName)
	p.nextToken() // Пропускаем "TYPE"

	return typeNode
}
func (p *Parser) parseVariableDeclaration() *ASTNode {
	varType := p.curTok.lexeme
	p.nextToken()
	if p.curTok.token != DCOLON && p.curTok.token != COMMA {
		p.tokens[p.pos+1] = LexItem{token: ERR, lexeme: "DCOLON OR COMMA EXPECTED", line: p.curTok.line, column: p.curTok.column, id: p.curTok.id}
	}
	p.nextToken()

	if p.curTok.lexeme == "intent" || p.curTok.lexeme == "INTENT" {
		p.nextToken()
		p.nextToken()
		if p.curTok.lexeme == "in" {
			varType = string(varType) + " readonly"
		}
		p.nextToken()
		p.nextToken()
		p.nextToken()

	}
	varDecl := NewASTNode("Variable Declaration", varType)

	for p.curTok.token == IDENT {
		varDecl.AddChild(NewASTNode("Variable", p.curTok.lexeme))
		p.nextToken()
		if p.curTok.token == COMMA {
			p.nextToken()
		} else {
			break
		}
	}

	return varDecl
}
func (p *Parser) parseIfStatement() *ASTNode {
	ifNode := NewASTNode("If Statement", "")
	p.nextToken()
	p.nextToken()

	// Условие
	conditionNode := p.parseExpression()
	ifNode.AddChild(conditionNode)
	p.nextToken()

	if p.curTok.token == THEN {
		p.nextToken()

	}

	// Тело IF
	bodyNode := NewASTNode("If Body", "")
	for p.curTok.token != ELSE && p.curTok.token != END && p.curTok.token != DONE {
		bodyNode.AddChild(p.parseStatement())
	}
	p.nextToken()
	p.nextToken()

	ifNode.AddChild(bodyNode)

	// ELSE IF
	for p.curTok.token == ELSE {
		p.nextToken()
		elseNode := NewASTNode("Else Block", "")
		for p.curTok.token != END && p.curTok.token != DONE {
			elseNode.AddChild(p.parseStatement())
		}
		ifNode.AddChild(elseNode)
	}

	return ifNode
}

func (p *Parser) parseWhileLoop() *ASTNode {
	whileNode := NewASTNode("While Loop", "")
	p.nextToken()

	// Условие
	conditionNode := p.parseExpression()
	whileNode.AddChild(conditionNode)

	// Тело
	bodyNode := NewASTNode("While Body", "")
	for p.curTok.token != END && p.curTok.token != DONE {
		bodyNode.AddChild(p.parseStatement())
	}
	whileNode.AddChild(bodyNode)

	p.nextToken() // Пропускаем "END WHILE"
	return whileNode
}

func (p *Parser) parseSelectCaseStatement() *ASTNode {
	// Создаем узел для SELECT CASE
	selectNode := NewASTNode("SelectCase", "")
	p.nextToken()
	p.nextToken()
	p.nextToken()

	exprNode := p.parseExpression()
	selectNode.AddChild(exprNode)
	p.nextToken()

	for p.curTok.token == CASE {
		p.nextToken()
		if p.curTok.token != DEFAULT {
			p.nextToken()
		}

		caseExprNode := p.parseCaseExpression()
		caseNode := NewASTNode("Case", "")
		caseNode.AddChild(caseExprNode)

		p.nextToken()
		bodyNode := NewASTNode("CaseBody", "")
		for p.curTok.token != CASE && p.curTok.token != END && p.curTok.token != DONE {
			bodyNode.AddChild(p.parseStatement())
		}

		caseNode.AddChild(bodyNode)
		selectNode.AddChild(caseNode)
	}

	p.nextToken() // Пропускаем "END SELECT"
	p.nextToken() // Пропускаем "END SELECT"

	return selectNode
}

func (p *Parser) parseCaseExpression() *ASTNode {
	if p.curTok.token == DEFAULT {
		return NewASTNode("CaseDefault", "default")
	}

	exprNode := NewASTNode("CaseExpression", "")
	for p.curTok.token != COLON && p.curTok.token != DONE {
		exprNode.AddChild(p.parseExpression())
		if p.curTok.token == COLON { // Диапазоны вида 1:5
			p.nextToken() // Пропускаем ":"
			exprNode.AddChild(p.parseExpression())
			break
		}
	}
	return exprNode
}

// Парсинг выражений
func (p *Parser) parseExpression() *ASTNode {
	expr := NewASTNode("Expression", p.curTok.lexeme)
	p.nextToken()
	return expr
}

func (p *Parser) parseStatement() *ASTNode {
	switch p.curTok.token {
	case IF:
		return p.parseIfStatement()
	case DO:
		return p.parseDoLoop()
	case SELECT:
		return p.parseSelectCaseStatement()
	case CALL:
		return p.parseCallStatement()
	case PRINT:
		return p.parsePrintStatement()
	case IDENT:
		return p.parseAssignmentOrExpression()
	case TYPE:
		return p.parseTypeDefinition()
	case INTEGER, REAL, CHARACTER:
		return p.parseVariableDeclaration()
	case IMPLICIT:
		return p.parseImplicit()
	default:
		p.nextToken()
		return nil
	}
}

func (p *Parser) parseCallStatement() *ASTNode {
	callNode := NewASTNode("Call", p.curTok.lexeme)
	p.nextToken() // Пропускаем "CALL"

	// Имя подпрограммы
	if p.curTok.token == IDENT {
		callNode.AddChild(NewASTNode("SubroutineName", p.curTok.lexeme))
		p.nextToken()
	}

	// Список аргументов
	if p.curTok.token == LPAREN {
		p.nextToken() // Пропускаем "("
		argsNode := NewASTNode("Arguments", "")
		for p.curTok.token != RPAREN && p.curTok.token != DONE {
			argsNode.AddChild(p.parseExpression())
			if p.curTok.token == COMMA {
				p.nextToken() // Пропускаем ","
			}
		}
		callNode.AddChild(argsNode)
		p.nextToken() // Пропускаем ")"
	}
	return callNode
}

func main() {
	inputFileName := "/home/marat/6_SEM/MTRAN/LR3/GOLANG/INPUT.TXT"
	inputFile, err := os.ReadFile(inputFileName)
	if err != nil {
		fmt.Printf("Error reading input file: %s\n", err)
		return
	}

	lexer := NewLexer(string(inputFile))
	lexer.GenerateTokens() // Генерация всех токенов

	// fmt.Println("Tokens:")
	// for _, token := range lexer.tokens {
	// 	if token.token == DONE {
	// 		break
	// 	}
	// 	if token.token == ERR {
	// 		fmt.Printf("Error: %s at line %d, column %d, ID: %d\n", token.lexeme, token.line, token.column, token.id)
	// 	}
	// 	fmt.Printf("Token: %-10s Lexeme: %-10s Line: %d, Column: %d, ID: %d\n", tokenMap[token.token], token.lexeme, token.line, token.column, token.id)
	// }

	parser := NewParser(lexer.tokens)
	ast := parser.Parse()
	if ast != nil {
		fmt.Println("Abstract Syntax Tree:")
		ast.PrintTree("", true)
	}
}
