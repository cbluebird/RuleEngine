package compiler

import (
	executor2 "engine/app/engine/executor"
	"engine/app/engine/token"
	"errors"
	"fmt"
)

// planner
type planner func(builder *Builder, curPre *precedence) (*executor2.Node, error)

func planValue(builder *Builder, curPre *precedence) (*executor2.Node, error) {
	if !builder.parser.hasNext() {
		return nil, nil
	}

	tok := builder.parser.next()
	switch tok.Kind {
	case token.OpenParen:
		// 最高优先级
		ret, err := builder.Build()
		if err != nil {
			return nil, err
		}

		// advance past the CloseParen token. We know that it's a CLAUSE_CLOSE, because at parse-time we check for unbalanced parens.
		builder.parser.next()

		// the stage we got represents all of the logic contained within the parens
		// but for technical reasons, we need to wrap this stage in a "noop" stage which breaks long chains of precedence.
		// see github #33.
		node := executor2.NewNode(nil, ret, executor2.NOOP, nil)
		return node, nil
	case token.Identifier:
		node := executor2.NewNode(nil, nil, executor2.VALUE, tok.Value)
		return node, nil
	case token.IntegerLiteral:
		node := executor2.NewNodeWithType(nil, nil, executor2.LITERAL, tok.Value, executor2.TypeInteger)
		return node, nil
	case token.FloatLiteral:
		node := executor2.NewNodeWithType(nil, nil, executor2.LITERAL, tok.Value, executor2.TypeFloat)
		return node, nil
	case token.BoolLiteral:
		node := executor2.NewNodeWithType(nil, nil, executor2.LITERAL, tok.Value, executor2.TypeBool)
		return node, nil
	case token.StringLiteral:
		node := executor2.NewNodeWithType(nil, nil, executor2.LITERAL, tok.Value, executor2.TypeString)
		return node, nil
	case token.Subtraction:
		ret, err := curPre.plan(builder)
		if err != nil {
			return nil, err
		}
		node := executor2.NewNodeWithPrefixFix(ret, executor2.NEGATIVE, nil)
		return node, nil
	case token.Addition: // token.Not,
		ret, err := curPre.plan(builder)
		if err != nil {
			return nil, err
		}
		node := executor2.NewNodeWithPrefixFix(ret, executor2.POSITIVE, nil)
		return node, nil
	case token.Not:
		builder.parser.rewind()
		return nil, nil
	default:
		errorMsg := fmt.Sprintf("Unable to plan token kind: '%s', value: '%v'", tok.Kind.String(), tok.Value)
		return nil, errors.New(errorMsg)
	}
}

type precedence struct {
	validKindsToSymbols map[token.Kind]executor2.Symbol // 当前优先级的token类型
	nextPrecedence      *precedence                     // 更高优先级的
	planner             planner
}

var (
	// * / %
	multiplicative = &precedence{
		validKindsToSymbols: executor2.MultiKindsToSymbol,
		planner:             planValue,
	}

	// + -
	additive = &precedence{
		validKindsToSymbols: executor2.AddKindsToSymbol,
		nextPrecedence:      multiplicative,
	}

	// > >= < <= == !=
	comparator = &precedence{
		validKindsToSymbols: executor2.CompareKindsToSymbol,
		nextPrecedence:      additive,
	}

	//!
	logicalNot = &precedence{
		validKindsToSymbols: executor2.NotKindsToSymbol,
		nextPrecedence:      comparator,
	}

	// &&
	logicalAnd = &precedence{
		validKindsToSymbols: executor2.AndKindsToSymbol,
		nextPrecedence:      logicalNot,
	}

	// ||
	logicalOr = &precedence{
		validKindsToSymbols: executor2.OrKindsToSymbol,
		nextPrecedence:      logicalAnd,
	}

	lowestPrecedence = logicalOr
)

func (p *precedence) plan(builder *Builder) (*executor2.Node, error) {
	var err error
	var leftNode, rightNode *executor2.Node

	if p.nextPrecedence != nil {
		leftNode, err = p.nextPrecedence.plan(builder)
		if err != nil {
			return nil, err
		}
	} else if p.planner != nil {
		leftNode, err = p.planner(builder, p)
		if err != nil {
			return nil, err
		}
	}

	for builder.parser.hasNext() {
		tok := builder.parser.next()

		if tok.Kind.IsEof() {
			break
		}

		symbol, exist := p.validKindsToSymbols[token.Kind(tok.Kind)]
		if !exist {
			break
		}

		rightNode, err = p.plan(builder)
		if err != nil {
			return nil, err
		}

		node := executor2.NewNode(leftNode, rightNode, symbol, nil)
		return node, nil
	}
	builder.parser.rewind()
	return leftNode, nil
}

type Builder struct {
	rootPlanner *precedence
	parser      *Parser
}

func NewBuilder(p *Parser) *Builder {
	return &Builder{
		rootPlanner: lowestPrecedence,
		parser:      p,
	}
}

func (b *Builder) Build() (*executor2.Node, error) {
	if b.parser == nil {
		return nil, errors.New("parse is nil")
	}

	if b.rootPlanner != nil {
		// TODO 树优化
		return b.rootPlanner.plan(b)
	}

	return nil, errors.New("build failed")
}
