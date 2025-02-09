package main

import (
	"fmt"
)

type Parser struct {
	lexer  *Lexer
	tokens []LexItem
	pos    int
}

func NewParser(lexer *Lexer) *Parser {
	tokens := []LexItem{}

	for {
		tok := lexer.NextToken()
		if tok.token == DONE {
			break
		}
		tokens = append(tokens, tok)
	}

	return &Parser{
		lexer:  lexer,
		tokens: tokens,
		pos:    0,
	}
}

func (p *Parser) parseModule() *ModuleNode {
	p.expect(PROGRAM)
	moduleName := p.expect(IDENT).lexeme
	p.expect(IMPLICIT)
	p.expect(NONE)

	body := []Node{}
	for !p.match(END) {
		if p.match(TYPE) {
			body = append(body, p.parseType())
		} else if p.match(FUNCTION) {
			body = append(body, p.parseFunction())
		} else if p.match(SUBROUTINE) {
			body = append(body, p.parseSubroutine())
		} else {
			p.error("Unexpected token in module")
		}
	}

	p.expect(END)
	p.expect(MODULE)
	return &ModuleNode{Name: moduleName, Body: body}
}

func (p *Parser) parseType() *TypeNode {
	p.expect(TYPE)
	typeName := p.expect(IDENT).lexeme
	p.expect(END)
	p.expect(TYPE)
	return &TypeNode{Name: typeName}
}

func (p *Parser) parseFunction() *FunctionNode {
	p.expect(FUNCTION)
	funcName := p.expect(IDENT).lexeme
	p.expect(RESULT)
	p.expect(IDENT)

	p.expect(IMPLICIT)
	p.expect(NONE)

	body := []Node{}
	for !p.match(END) {
		if p.match(PRINT) {
			body = append(body, p.parsePrint())
		} else {
			p.error("Unexpected token inside function")
		}
	}

	p.expect(END)
	p.expect(FUNCTION)

	return &FunctionNode{Name: funcName, Body: body}
}

func (p *Parser) parseSubroutine() *SubroutineNode {
	p.expect(SUBROUTINE)
	subName := p.expect(IDENT).lexeme

	p.expect(IMPLICIT)
	p.expect(NONE)

	body := []Node{}
	for !p.match(END) {
		if p.match(PRINT) {
			body = append(body, p.parsePrint())
		} else {
			p.error("Unexpected token inside subroutine")
		}
	}

	p.expect(END)
	p.expect(SUBROUTINE)

	return &SubroutineNode{Name: subName, Body: body}
}

func (p *Parser) parsePrint() *PrintNode {
	p.expect(PRINT)
	p.expect(MULT)
	message := p.expect(SCONST).lexeme
	return &PrintNode{Message: message}
}

// Вспомогательные методы
func (p *Parser) expect(expected Token) LexItem {
	if p.pos >= len(p.tokens) {
		p.error("Unexpected end of input")
	}
	token := p.tokens[p.pos]
	if token.token != expected {
		p.error(fmt.Sprintf("Expected %s, got %s", tokenMap[expected], tokenMap[token.token]))
	}
	p.pos++
	return token
}

func (p *Parser) match(expected Token) bool {
	if p.pos < len(p.tokens) && p.tokens[p.pos].token == expected {
		p.pos++
		return true
	}
	return false
}

func (p *Parser) error(message string) {
	fmt.Printf("Syntax Error: %s\n", message)
	os.Exit(1)
}