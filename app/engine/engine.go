package engine

import (
	"engine/app/engine/compiler"
	"engine/app/engine/executor"
	"log"
)

type Engine struct {
	scanner *compiler.Scanner
	parser  *compiler.Parser
	builder *compiler.Builder
	node    *executor.Node
}

func NewEngine(s string) *Engine {
	scanner := compiler.NewScanner(s)
	tokens, err := scanner.Lexer()
	if err != nil {
		log.Println(err)
		return nil
	}
	parser := compiler.NewParser(tokens)
	parser.Print()
	err = parser.CheckBalance()
	if err != nil {
		log.Println(err)
		return nil
	}
	err = parser.ParseSyntax()
	if err != nil {
		log.Println(err)
		return nil
	}
	bulider := compiler.NewBuilder(parser)
	node, err_ := bulider.Build()
	if err_ != nil {
		log.Println(err)
		return nil
	}
	return &Engine{
		scanner: scanner,
		parser:  parser,
		builder: bulider,
		node:    node,
	}
}

func (e *Engine) Calculate(parameters map[string]interface{}) error {
	err := e.node.Eval(parameters)
	return err
}

func (e *Engine) GetVal() (interface{}, executor.TypeFlags) {
	return e.node.GetVal()
}
