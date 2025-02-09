package main

import (
	"fmt"
	"os"
)

// Узлы AST
type Node interface {
	Print(level int)
}

// Узел для модуля
type ModuleNode struct {
	Name      string
	Body      []Node
}

func (m *ModuleNode) Print(level int) {
	fmt.Printf("%sModule: %s\n", indent(level), m.Name)
	for _, stmt := range m.Body {
		stmt.Print(level + 1)
	}
}

// Узел для объявления типа
type TypeNode struct {
	Name   string
	Fields []VarDeclNode
}

func (t *TypeNode) Print(level int) {
	fmt.Printf("%sType: %s\n", indent(level), t.Name)
	for _, field := range t.Fields {
		field.Print(level + 1)
	}
}

// Узел для объявления переменной
type VarDeclNode struct {
	Type string
	Name string
}

func (v *VarDeclNode) Print(level int) {
	fmt.Printf("%sVar: %s %s\n", indent(level), v.Type, v.Name)
}

// Узел для функций
type FunctionNode struct {
	Name       string
	Parameters []VarDeclNode
	Body       []Node
}

func (f *FunctionNode) Print(level int) {
	fmt.Printf("%sFunction: %s\n", indent(level), f.Name)
	for _, param := range f.Parameters {
		param.Print(level + 1)
	}
	for _, stmt := range f.Body {
		stmt.Print(level + 1)
	}
}

// Узел для подпрограмм
type SubroutineNode struct {
	Name       string
	Parameters []VarDeclNode
	Body       []Node
}

func (s *SubroutineNode) Print(level int) {
	fmt.Printf("%sSubroutine: %s\n", indent(level), s.Name)
	for _, param := range s.Parameters {
		param.Print(level + 1)
	}
	for _, stmt := range s.Body {
		stmt.Print(level + 1)
	}
}

// Операторы
type PrintNode struct {
	Message string
}

func (p *PrintNode) Print(level int) {
	fmt.Printf("%sPrint: \"%s\"\n", indent(level), p.Message)
}

type AssignNode struct {
	Variable string
	Value    string
}

func (a *AssignNode) Print(level int) {
	fmt.Printf("%sAssign: %s = %s\n", indent(level), a.Variable, a.Value)
}

// Вспомогательная функция для форматирования вывода
func indent(level int) string {
	return string([]rune("  ")[0:level*2])
}