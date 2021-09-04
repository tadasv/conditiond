package condition

/*

literal := number | string | null | float ;

list := '[', *(literal | function) ,']' ;
function := '{', string, ':', list | literal | function, '}' ;

expression := function | list | literal ;

*/

import (
	"fmt"
	"io"
)

type NodeType int

const (
	NodeTypeLiteral NodeType = iota
	NodeTypeArray
	NodeTypeFunction
	NodeTypeExpression
)

type consumeFunc func(*parser) (consumeFunc, error)

type Node struct {
	Type     NodeType
	Token    Token
	Parent   *Node
	Children []*Node
}

func (n *Node) appendChild(child *Node) {
	if n.Children == nil {
		n.Children = []*Node{
			child,
		}
	} else {
		n.Children = append(n.Children, child)
	}
}

func newLiteralNode(parent *Node, token Token) *Node {
	return &Node{
		Type:   NodeTypeLiteral,
		Token:  token,
		Parent: parent,
	}
}

func newFunctionNode(parent *Node, token Token) *Node {
	return &Node{
		Type:   NodeTypeFunction,
		Token:  token,
		Parent: parent,
	}
}

func newArrayNode(parent *Node) *Node {
	return &Node{
		Type:   NodeTypeArray,
		Parent: parent,
	}
}

func Parse(expression string) (*Node, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		return nil, err
	}

	p := &parser{
		tokens: tokens,
	}

	return buildAST(p)
}

type parser struct {
	tokens []Token
	pos    int

	lastNode *Node
}

func (p *parser) consume() (Token, error) {
	if p.pos >= len(p.tokens) {
		return Token{}, io.EOF
	}
	token := p.tokens[p.pos]
	p.pos++
	return token, nil
}

func (p *parser) peek() (Token, error) {
	if p.pos >= len(p.tokens) {
		return Token{}, io.EOF
	}
	return p.tokens[p.pos], nil
}

func (p *parser) getRootNode() *Node {
	currentNode := p.lastNode
	for currentNode != nil {
		if currentNode.Parent == nil {
			break
		}

		currentNode = currentNode.Parent
	}

	return currentNode
}

func buildAST(p *parser) (*Node, error) {

	next, err := consumeExpression(p)
	if err != nil {
		return nil, err
	}

	for {
		if next != nil {
			next, err = next(p)
			if err == io.EOF {
				break
			}

			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return p.getRootNode(), nil
}

func consumeLiteral(p *parser) (consumeFunc, error) {
	token, err := p.consume()
	if err != nil {
		return nil, err
	}

	if token.Type != TokenTypeLiteral {
		return nil, fmt.Errorf("expected literal, got %s instead", token.String())
	}

	newNode := &Node{
		Type:   NodeTypeLiteral,
		Token:  token,
		Parent: p.lastNode,
	}

	if p.lastNode != nil {
		p.lastNode.appendChild(newNode)

		switch p.lastNode.Type {
		case NodeTypeArray:
			return consumeArrayValue, nil
		case NodeTypeFunction:
			return consumeFunctionEnd, nil
		}
	} else {
		p.lastNode = newNode
	}

	return nil, nil
}

func consumeArrayStart(p *parser) (consumeFunc, error) {
	// consume [
	token, err := p.consume()
	if err != nil {
		return nil, err
	}

	if token.Type != TokenTypeBracketOpen {
		return nil, fmt.Errorf("expected [, got %s", token.String())
	}

	arrayNode := newArrayNode(p.lastNode)
	if p.lastNode != nil {
		p.lastNode.appendChild(arrayNode)
	}
	p.lastNode = arrayNode

	return consumeArrayValue, nil
}

func consumeArrayEnd(p *parser) (consumeFunc, error) {
	// consume ]
	token, err := p.consume()
	if err != nil {
		return nil, err
	}

	if token.Type != TokenTypeBracketClose {
		return nil, fmt.Errorf("expected ], got %s", token.String())
	}

	if p.lastNode.Parent == nil {
		return nil, io.EOF
	}

	switch p.lastNode.Parent.Type {
	case NodeTypeFunction:
		p.lastNode = p.lastNode.Parent
		return consumeFunctionEnd, nil
	}

	return nil, fmt.Errorf("unsupported parent node for array: %#v\n", p.lastNode.Parent.Type)
}

func consumeArrayValue(p *parser) (consumeFunc, error) {
	valueToken, err := p.peek()
	if err != nil {
		return nil, err
	}

	switch valueToken.Type {
	case TokenTypeBracketClose:
		return consumeArrayEnd, nil
	case TokenTypeLiteral:
		return consumeLiteral, nil
	case TokenTypeBraceOpen:
		return consumeFunctionStart, nil
	}

	return nil, fmt.Errorf("unexpected token as array value: %s", valueToken.String())
}

func consumeFunctionStart(p *parser) (consumeFunc, error) {
	// consume {
	token, err := p.consume()
	if err != nil {
		return nil, err
	}

	if token.Type != TokenTypeBraceOpen {
		return nil, fmt.Errorf("expected {, got %s", token.String())
	}

	functionNameToken, err := p.consume()
	if err != nil {
		return nil, err
	}

	functionNode := newFunctionNode(p.lastNode, functionNameToken)
	if p.lastNode != nil {
		p.lastNode.appendChild(functionNode)
	}
	p.lastNode = functionNode

	return consumeFunctionValue, nil
}

func consumeFunctionEnd(p *parser) (consumeFunc, error) {
	// consume }
	token, err := p.consume()
	if err != nil {
		return nil, err
	}

	if token.Type != TokenTypeBraceClose {
		return nil, fmt.Errorf("expected }, got %s", token.String())
	}

	if p.lastNode.Parent == nil {
		return nil, io.EOF
	}

	switch p.lastNode.Parent.Type {
	case NodeTypeArray:
		p.lastNode = p.lastNode.Parent
		return consumeArrayValue, nil
	case NodeTypeFunction:
		p.lastNode = p.lastNode.Parent
		return consumeFunctionEnd, nil
	}

	return nil, nil
}

func consumeFunctionValue(p *parser) (consumeFunc, error) {
	valueStartToken, err := p.peek()
	if err != nil {
		return nil, err
	}

	switch valueStartToken.Type {
	case TokenTypeLiteral:
		return consumeLiteral, nil
	case TokenTypeBraceOpen:
		return consumeFunctionStart, nil
	case TokenTypeBracketOpen:
		return consumeArrayStart, nil
	}

	return nil, fmt.Errorf("expected literal, function or array as a value for function expression, got %s instead", valueStartToken.String())
}

func consumeExpression(p *parser) (consumeFunc, error) {
	token, err := p.peek()
	if err != nil {
		return nil, err
	}

	switch token.Type {
	case TokenTypeBraceOpen:
		return consumeFunctionStart, nil
	case TokenTypeBracketOpen:
		return consumeArrayStart, nil
	case TokenTypeLiteral:
		return consumeLiteral, nil
	default:
		return nil, fmt.Errorf("unexpected token: %s", token.String())
	}
}
