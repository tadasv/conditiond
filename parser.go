package condition

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type TokenType int

const (
	TokenTypeBracketOpen TokenType = iota
	TokenTypeBracketClose
	TokenTypeBraceOpen
	TokenTypeBraceClose
	TokenTypeLiteral
)

type LiteralType int

const (
	LiteralTypeNumber LiteralType = iota
	LiteralTypeString
	LiteralTypeBool
	LiteralTypeNull
)

type Token struct {
	Type        TokenType
	LiteralType LiteralType
	Value       interface{}
}

func (t Token) String() string {
	switch t.Type {
	case TokenTypeBraceOpen:
		return "BRACE_OPEN"
	case TokenTypeBraceClose:
		return "BRACE_CLOSE"
	case TokenTypeBracketOpen:
		return "BRACKET_OPEN"
	case TokenTypeBracketClose:
		return "BRACKET_CLOSE"
	case TokenTypeLiteral:
		switch t.LiteralType {
		case LiteralTypeBool:
			return fmt.Sprintf("LITERAL<bool::%v>", t.Value)
		case LiteralTypeString:
			return fmt.Sprintf("LITERAL<string::%v>", t.Value)
		case LiteralTypeNumber:
			return fmt.Sprintf("LITERAL<number::%v>", t.Value)
		case LiteralTypeNull:
			return "LITERAL<null>"
		}
		return fmt.Sprintf("LITERAL<unknown::%v>", t.Value)
	}

	return "UNKNOWN"
}

func tokenize(value string) ([]Token, error) {
	decoder := json.NewDecoder(strings.NewReader(value))

	tokens := []Token{}

	for {
		jsonToken, err := decoder.Token()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		switch v := jsonToken.(type) {
		case json.Delim:
			switch v {
			case json.Delim('['):
				tokens = append(tokens, Token{Type: TokenTypeBracketOpen})
			case json.Delim(']'):
				tokens = append(tokens, Token{Type: TokenTypeBracketClose})
			case json.Delim('{'):
				tokens = append(tokens, Token{Type: TokenTypeBraceOpen})
			case json.Delim('}'):
				tokens = append(tokens, Token{Type: TokenTypeBraceClose})
			}
		case bool:
			tokens = append(tokens, Token{
				Type:        TokenTypeLiteral,
				Value:       v,
				LiteralType: LiteralTypeBool,
			})
		case string:
			tokens = append(tokens, Token{
				Type:        TokenTypeLiteral,
				Value:       v,
				LiteralType: LiteralTypeString,
			})
		case float64:
			tokens = append(tokens, Token{
				Type:        TokenTypeLiteral,
				Value:       v,
				LiteralType: LiteralTypeNumber,
			})
		case nil:
			tokens = append(tokens, Token{
				Type:        TokenTypeLiteral,
				LiteralType: LiteralTypeNull,
			})
		}
	}

	return tokens, nil
}
