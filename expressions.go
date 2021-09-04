package condition

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	errExpectedArrayInput = "expected array input got %q"
	errExpectedNArguments = "expected %d argument(s), got %d"
	errExpectedNumber     = "expected number as an argument"
)

func newExpression(name string, handler ExpressionFunc) ExpressionFunc {
	return func(e *Evaluator, n *Node) (interface{}, error) {
		res, err := handler(e, n)
		if err != nil {
			return nil, fmt.Errorf("%s expression: %s", name, err.Error())
		}
		return res, nil
	}
}

func OrExpressionHandler(e *Evaluator, n *Node) (interface{}, error) {
	child := n.Children[0]
	if child.Type != NodeTypeArray {
		return nil, fmt.Errorf(errExpectedArrayInput, child.Type)
	}

	for _, n := range child.Children {
		res, err := e.evaluateNode(n)
		if err != nil {
			return nil, err
		}

		asBool := castToBool(res)
		if asBool == true {
			return true, nil
		}

	}

	return false, nil
}

func AndExpressionHandler(e *Evaluator, n *Node) (interface{}, error) {
	child := n.Children[0]
	if child.Type != NodeTypeArray {
		return nil, fmt.Errorf(errExpectedArrayInput, child.Type)
	}

	for _, n := range child.Children {
		res, err := e.evaluateNode(n)
		if err != nil {
			return nil, err
		}

		asBool := castToBool(res)
		if asBool == false {
			return false, nil
		}

	}

	return true, nil
}

func EqExpressionHandler(e *Evaluator, n *Node) (interface{}, error) {
	if len(n.Children) != 1 && n.Children[0].Type != NodeTypeArray {
		return false, nil
	}

	params := n.Children[0]
	if len(params.Children) != 2 {
		return false, nil
	}

	resA, err := e.evaluateNode(params.Children[0])
	if err != nil {
		return nil, err
	}

	resB, err := e.evaluateNode(params.Children[1])
	if err != nil {
		return nil, err
	}

	return resA == resB, nil
}

func NotExpressionHandler(e *Evaluator, n *Node) (interface{}, error) {
	if len(n.Children) != 1 {
		return nil, fmt.Errorf(errExpectedNArguments, 1, len(n.Children))
	}
	res, err := e.evaluateNode(n.Children[0])
	if err != nil {
		return nil, err
	}

	return !castToBool(res), nil
}

func GtExpressionHandler(e *Evaluator, n *Node) (interface{}, error) {
	params := n.Children[0]

	if len(params.Children) != 2 {
		return nil, fmt.Errorf(errExpectedNArguments, 2, len(params.Children))
	}

	resA, err := e.evaluateNode(params.Children[0])
	if err != nil {
		return nil, err
	}

	resB, err := e.evaluateNode(params.Children[1])
	if err != nil {
		return nil, err
	}

	floatA, ok := resA.(float64)
	if !ok {
		return nil, fmt.Errorf(errExpectedNumber)
	}

	floatB, ok := resB.(float64)
	if !ok {
		return nil, fmt.Errorf(errExpectedNumber)
	}

	return floatA > floatB, nil
}

func GteExpressionHandler(e *Evaluator, n *Node) (interface{}, error) {
	params := n.Children[0]

	if len(params.Children) != 2 {
		return nil, fmt.Errorf(errExpectedNArguments, 2, len(params.Children))
	}

	resA, err := e.evaluateNode(params.Children[0])
	if err != nil {
		return nil, err
	}

	resB, err := e.evaluateNode(params.Children[1])
	if err != nil {
		return nil, err
	}

	floatA, ok := resA.(float64)
	if !ok {
		return nil, fmt.Errorf(errExpectedNumber)
	}

	floatB, ok := resB.(float64)
	if !ok {
		return nil, fmt.Errorf(errExpectedNumber)
	}

	return floatA >= floatB, nil
}

func LtExpressionHandler(e *Evaluator, n *Node) (interface{}, error) {
	params := n.Children[0]

	if len(params.Children) != 2 {
		return nil, fmt.Errorf(errExpectedNArguments, 2, len(params.Children))
	}

	resA, err := e.evaluateNode(params.Children[0])
	if err != nil {
		return nil, err
	}

	resB, err := e.evaluateNode(params.Children[1])
	if err != nil {
		return nil, err
	}

	floatA, ok := resA.(float64)
	if !ok {
		return nil, fmt.Errorf(errExpectedNumber)
	}

	floatB, ok := resB.(float64)
	if !ok {
		return nil, fmt.Errorf(errExpectedNumber)
	}

	return floatA < floatB, nil
}

func LteExpressionHandler(e *Evaluator, n *Node) (interface{}, error) {
	params := n.Children[0]

	if len(params.Children) != 2 {
		return nil, fmt.Errorf(errExpectedNArguments, 2, len(params.Children))
	}

	resA, err := e.evaluateNode(params.Children[0])
	if err != nil {
		return nil, err
	}

	resB, err := e.evaluateNode(params.Children[1])
	if err != nil {
		return nil, err
	}

	floatA, ok := resA.(float64)
	if !ok {
		return nil, fmt.Errorf(errExpectedNumber)
	}

	floatB, ok := resB.(float64)
	if !ok {
		return nil, fmt.Errorf(errExpectedNumber)
	}

	return floatA <= floatB, nil
}

func IfExpressionHandler(e *Evaluator, n *Node) (interface{}, error) {
	if n.Children[0].Type != NodeTypeArray {
		return nil, fmt.Errorf(errExpectedArrayInput, n.Children[0].Type)
	}

	params := n.Children[0].Children
	if len(params) < 2 {
		return nil, fmt.Errorf(errExpectedNArguments, 2, len(params))
	} else if len(params) > 3 {
		return nil, fmt.Errorf(errExpectedNArguments, 3, len(params))
	}

	predicateRes, err := e.evaluateNode(params[0])
	if err != nil {
		return nil, err
	}

	resAsBool := castToBool(predicateRes)
	if resAsBool {
		return e.evaluateNode(params[1])
	} else if len(params) > 2 {
		return e.evaluateNode(params[2])
	}

	return nil, nil
}

func Sha1modExpressionHandler(e *Evaluator, n *Node) (interface{}, error) {
	resA, err := e.evaluateNode(n.Children[0].Children[0])
	if err != nil {
		return nil, err
	}
	resB, err := e.evaluateNode(n.Children[0].Children[1])
	if err != nil {
		return nil, err
	}

	resBuint := uint64(resB.(float64))

	keyValue, err := json.Marshal(resA)
	if err != nil {
		return nil, err
	}

	hash := sha1.Sum(keyValue)
	value := binary.BigEndian.Uint64(hash[:8])

	result := value % resBuint
	return float64(result), nil
}

func ContextExpressionHandler(e *Evaluator, n *Node) (interface{}, error) {
	params := n.Children[0]

	if params.Type != NodeTypeArray {
		return nil, fmt.Errorf(errExpectedArrayInput, params.Type)
	}

	evaluatedPath := []interface{}{}
	for _, child := range params.Children {
		res, err := e.evaluateNode(child)
		if err != nil {
			return nil, err
		}
		evaluatedPath = append(evaluatedPath, res)
	}

	ctx := ""
	var decodedData interface{}

	switch t := e.context.(type) {
	case string:
		ctx = t
	case json.RawMessage:
		ctx = string(t)
	default:
		decodedData = t
	}

	if decodedData == nil {
		if err := json.NewDecoder(strings.NewReader(ctx)).Decode(&decodedData); err != nil {
			return nil, err
		}
	}

	val := recursiveGet(decodedData, evaluatedPath)
	switch val.(type) {
	case pathTypeMismatch:
		return nil, fmt.Errorf("only strings and integers supported as input values")
	case unknownPathType:
		return nil, fmt.Errorf("only strings and integers supported as input values")
	case notFound:
		return nil, nil
	}

	return val, nil
}

type pathTypeMismatch struct{}
type notFound struct{}
type unknownPathType struct{}

func recursiveGet(data interface{}, path []interface{}) interface{} {
	if len(path) == 0 {
		switch data.(type) {
		case string:
			return data.(string)
		case float64:
			return data.(float64)
		case bool:
			return data.(bool)
		case nil:
			return nil
		case []interface{}:
			return data
		case map[string]interface{}:
			return data
		}
	}

	switch path[0].(type) {
	case string:
		switch data.(type) {
		case map[string]interface{}:
			for k, v := range data.(map[string]interface{}) {
				if k == path[0].(string) {
					return recursiveGet(v, path[1:])
				}
			}
			return notFound{}
		default:
			return pathTypeMismatch{}
		}
	case float64:
		// When parsing JSON we do not get ints, only floats. This is why we're
		// casting float64 to int
		switch data.(type) {
		case []interface{}:
			for i, v := range data.([]interface{}) {
				if i == int(path[0].(float64)) {
					return recursiveGet(v, path[1:])
				}
			}
			return notFound{}
		default:
			return pathTypeMismatch{}
		}
	}

	return unknownPathType{}
}

func castToBool(a interface{}) bool {
	switch v := a.(type) {
	case nil:
		return false
	/*case string:
		if len(v) == 0 {
			return false
		}
		return true
	case float32:
		if v == 0 {
			return false
		}
		return true
	case float64:
		if v == 0 {
			return false
		}
		return true
	case int32:
		if v == 0 {
			return false
		}
		return true
	case int64:
		if v == 0 {
			return false
		}
		return true
	case int:
		if v == 0 {
			return false
		}
		return true
	*/
	case bool:
		return v
	}
	return true
}
