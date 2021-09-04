package condition

import (
	"errors"
	"fmt"
)

type ExpressionFunc func(*Evaluator, *Node) (interface{}, error)

type Evaluator struct {
	funcs   map[string]ExpressionFunc
	context interface{}
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		funcs:   map[string]ExpressionFunc{},
		context: nil,
	}
}

func (e Evaluator) AddHandler(name string, handler ExpressionFunc) {
	e.funcs[name] = newExpression(name, handler)
}

func (e *Evaluator) Evaluate(ctx interface{}, root *Node) (interface{}, error) {
	e.context = ctx
	return e.evaluateNode(root)
}

func (e *Evaluator) evaluateNode(n *Node) (interface{}, error) {
	if n == nil {
		return nil, errors.New("received nil AST node as an input to evaluator")
	}

	switch n.Type {
	case NodeTypeLiteral:
		return n.Token.Value, nil
	case NodeTypeFunction:
		funcName := n.Token.Value.(string)
		f, ok := e.funcs[funcName]
		if !ok {
			return nil, fmt.Errorf("no expression handler bound to %q", funcName)
		}
		return f(e, n)
	case NodeTypeArray:
		res := make([]interface{}, len(n.Children))
		for i, child := range n.Children {
			val, err := e.evaluateNode(child)
			if err != nil {
				return nil, err
			}

			res[i] = val
		}
		return res, nil
	}
	return nil, nil
}
