package condition

import (
	"testing"
)

func TestNot(t *testing.T) {
	evaluator := NewEvaluator()
	evaluator.AddHandler("not", NotExpressionHandler)

	testCases := []struct {
		in  string
		out bool
	}{
		{
			in:  `{"not": true}`,
			out: false,
		},
		{
			in:  `{"not": false}`,
			out: true,
		},
	}

	for _, test := range testCases {
		root, err := Parse(test.in)
		if err != nil {
			t.Errorf("%q got an error: %s", test.in, err.Error())
		} else {
			res, err := evaluator.Evaluate(nil, root)
			if err != nil {
				t.Errorf("%q got an error: %s", test.in, err.Error())
			} else {
				resAsBool := res.(bool)
				if resAsBool != test.out {
					t.Errorf("%q expected %t got %t", test.in, test.out, resAsBool)
				}
			}
		}
	}
}

func TestAnd(t *testing.T) {
	evaluator := NewEvaluator()
	evaluator.AddHandler("and", AndExpressionHandler)

	testCases := []struct {
		in  string
		out bool
	}{
		{
			in:  `{"and": [true, false]}`,
			out: false,
		},
		{
			in:  `{"and": [true, true]}`,
			out: true,
		},
	}

	for _, test := range testCases {
		root, err := Parse(test.in)
		if err != nil {
			t.Errorf("%q got an error: %s", test.in, err.Error())
		} else {
			res, err := evaluator.Evaluate(nil, root)
			if err != nil {
				t.Errorf("%q got an error: %s", test.in, err.Error())
			} else {
				resAsBool := res.(bool)
				if resAsBool != test.out {
					t.Errorf("%q expected %t got %t", test.in, test.out, resAsBool)
				}
			}
		}
	}
}

func TestOr(t *testing.T) {
	evaluator := NewEvaluator()
	evaluator.AddHandler("or", OrExpressionHandler)

	testCases := []struct {
		in  string
		out bool
	}{
		{
			in:  `{"or": [true, false]}`,
			out: true,
		},
		{
			in:  `{"or": [false, false]}`,
			out: false,
		},
	}

	for _, test := range testCases {
		root, err := Parse(test.in)
		if err != nil {
			t.Errorf("%q got an error: %s", test.in, err.Error())
		} else {
			res, err := evaluator.Evaluate(nil, root)
			if err != nil {
				t.Errorf("%q got an error: %s", test.in, err.Error())
			} else {
				resAsBool := res.(bool)
				if resAsBool != test.out {
					t.Errorf("%q expected %t got %t", test.in, test.out, resAsBool)
				}
			}
		}
	}
}

func TestGte(t *testing.T) {
	evaluator := NewEvaluator()
	evaluator.AddHandler("gte", GteExpressionHandler)

	testCases := []struct {
		in  string
		out bool
	}{
		{
			in:  `{"gte": [1, 0]}`,
			out: true,
		},
		{
			in:  `{"gte": [1.5, 1]}`,
			out: true,
		},
		{
			in:  `{"gte": [1.5, 2]}`,
			out: false,
		},
		{
			in:  `{"gte": [1.5, 1.5]}`,
			out: true,
		},
		{
			in:  `{"gte": [5, 5]}`,
			out: true,
		},
	}

	for _, test := range testCases {
		root, err := Parse(test.in)
		if err != nil {
			t.Errorf("%q got an error: %s", test.in, err.Error())
		} else {
			res, err := evaluator.Evaluate(nil, root)
			if err != nil {
				t.Errorf("%q got an error: %s", test.in, err.Error())
			} else {
				resAsBool := res.(bool)
				if resAsBool != test.out {
					t.Errorf("%q expected %t got %t", test.in, test.out, resAsBool)
				}
			}
		}
	}
}

func TestGt(t *testing.T) {
	evaluator := NewEvaluator()
	evaluator.AddHandler("gt", GtExpressionHandler)

	testCases := []struct {
		in  string
		out bool
	}{
		{
			in:  `{"gt": [1, 0]}`,
			out: true,
		},
		{
			in:  `{"gt": [1.5, 1]}`,
			out: true,
		},
		{
			in:  `{"gt": [1.5, 2]}`,
			out: false,
		},
		{
			in:  `{"gt": [1.5, 1.5]}`,
			out: false,
		},
		{
			in:  `{"gt": [5, 5]}`,
			out: false,
		},
	}

	for _, test := range testCases {
		root, err := Parse(test.in)
		if err != nil {
			t.Errorf("%q got an error: %s", test.in, err.Error())
		} else {
			res, err := evaluator.Evaluate(nil, root)
			if err != nil {
				t.Errorf("%q got an error: %s", test.in, err.Error())
			} else {
				resAsBool := res.(bool)
				if resAsBool != test.out {
					t.Errorf("%q expected %t got %t", test.in, test.out, resAsBool)
				}
			}
		}
	}
}

func TestLte(t *testing.T) {
	evaluator := NewEvaluator()
	evaluator.AddHandler("lte", LteExpressionHandler)

	testCases := []struct {
		in  string
		out bool
	}{
		{
			in:  `{"lte": [1, 0]}`,
			out: false,
		},
		{
			in:  `{"lte": [1.5, 1]}`,
			out: false,
		},
		{
			in:  `{"lte": [1.5, 2]}`,
			out: true,
		},
		{
			in:  `{"lte": [1.5, 1.5]}`,
			out: true,
		},
		{
			in:  `{"lte": [5, 5]}`,
			out: true,
		},
	}

	for _, test := range testCases {
		root, err := Parse(test.in)
		if err != nil {
			t.Errorf("%q got an error: %s", test.in, err.Error())
		} else {
			res, err := evaluator.Evaluate(nil, root)
			if err != nil {
				t.Errorf("%q got an error: %s", test.in, err.Error())
			} else {
				resAsBool := res.(bool)
				if resAsBool != test.out {
					t.Errorf("%q expected %t got %t", test.in, test.out, resAsBool)
				}
			}
		}
	}
}

func TestLt(t *testing.T) {
	evaluator := NewEvaluator()
	evaluator.AddHandler("lt", LtExpressionHandler)

	testCases := []struct {
		in  string
		out bool
	}{
		{
			in:  `{"lt": [1, 0]}`,
			out: false,
		},
		{
			in:  `{"lt": [1.5, 1]}`,
			out: false,
		},
		{
			in:  `{"lt": [1.5, 2]}`,
			out: true,
		},
		{
			in:  `{"lt": [1.5, 1.5]}`,
			out: false,
		},
		{
			in:  `{"lt": [5, 5]}`,
			out: false,
		},
	}

	for _, test := range testCases {
		root, err := Parse(test.in)
		if err != nil {
			t.Errorf("%q got an error: %s", test.in, err.Error())
		} else {
			res, err := evaluator.Evaluate(nil, root)
			if err != nil {
				t.Errorf("%q got an error: %s", test.in, err.Error())
			} else {
				resAsBool := res.(bool)
				if resAsBool != test.out {
					t.Errorf("%q expected %t got %t", test.in, test.out, resAsBool)
				}
			}
		}
	}
}

func TestEq(t *testing.T) {
	evaluator := NewEvaluator()
	evaluator.AddHandler("eq", EqExpressionHandler)
	evaluator.AddHandler("or", OrExpressionHandler)
	evaluator.AddHandler("and", AndExpressionHandler)

	testCases := []struct {
		in  string
		out bool
	}{
		{
			in:  `{"eq": [true, true]}`,
			out: true,
		},
		{
			in:  `{"eq": [true, false]}`,
			out: false,
		},
		{
			in:  `{"eq": [1, 1]}`,
			out: true,
		},
		{
			in:  `{"eq": ["a", "a"]}`,
			out: true,
		},
		{
			in:  `{"eq": ["a", "ab"]}`,
			out: false,
		},
		{
			in:  `{"eq": [null, null]}`,
			out: true,
		},
		{
			in:  `{"eq": [{"or": [true]}, {"and": [true]}]}`,
			out: true,
		},
	}

	for _, test := range testCases {
		root, err := Parse(test.in)
		if err != nil {
			t.Errorf("%q got an error: %s", test.in, err.Error())
		} else {
			res, err := evaluator.Evaluate(nil, root)
			if err != nil {
				t.Errorf("%q got an error: %s", test.in, err.Error())
			} else {
				resAsBool := res.(bool)
				if resAsBool != test.out {
					t.Errorf("%q expected %t got %t", test.in, test.out, resAsBool)
				}
			}
		}
	}
}

func TestHashmod(t *testing.T) {
	evaluator := NewEvaluator()
	evaluator.AddHandler("sha1mod", Sha1modExpressionHandler)

	testCases := []struct {
		in  string
		out interface{}
	}{
		{
			in:  `{"sha1mod": ["value", 100]}`,
			out: float64(3),
		},
		{
			in:  `{"sha1mod": ["", 100]}`,
			out: float64(74),
		},
		{
			in:  `{"sha1mod": [null, 100]}`,
			out: float64(24),
		},
		{
			in:  `{"sha1mod": [100, 100]}`,
			out: float64(41),
		},
	}

	for _, test := range testCases {
		root, err := Parse(test.in)
		if err != nil {
			t.Errorf("%q got an error: %s", test.in, err.Error())
		} else {
			res, err := evaluator.Evaluate(nil, root)
			if err != nil {
				t.Errorf("%q got an error: %s", test.in, err.Error())
			} else {
				if res != test.out {
					t.Errorf("%q expected %t got %t", test.in, test.out, res)
				}
			}
		}
	}
}

func TestIf(t *testing.T) {
	evaluator := Evaluator{
		funcs: map[string]ExpressionFunc{
			"if": IfExpressionHandler,
		},
	}

	testCases := []struct {
		in  string
		out interface{}
	}{
		{
			in:  `{"if": [true, true]}`,
			out: true,
		},
		{
			in:  `{"if": [true, false]}`,
			out: false,
		},
		{
			in:  `{"if": [false, false, true]}`,
			out: true,
		},
		{
			in:  `{"if": [false, false]}`,
			out: nil,
		},
	}

	for _, test := range testCases {
		root, err := Parse(test.in)
		if err != nil {
			t.Errorf("%q got an error: %s", test.in, err.Error())
		} else {
			res, err := evaluator.Evaluate(nil, root)
			if err != nil {
				t.Errorf("%q got an error: %s", test.in, err.Error())
			} else {
				if res != test.out {
					t.Errorf("%q expected %t got %t", test.in, test.out, res)
				}
			}
		}
	}
}

func TestContext(t *testing.T) {
	context := `{
			"key": "value",
			"key2": 123,
			"key3": {
				"subkey1": 1,
				"subkey2": [1, 2]
			},
			"key4": true,
			"key5": null
		}`
	evaluator := NewEvaluator()
	evaluator.AddHandler("ctx", ContextExpressionHandler)

	testCases := []struct {
		in  string
		out interface{}
	}{
		{
			in:  `{"ctx": ["key"]}`,
			out: "value",
		},
		{
			in:  `{"ctx": ["key2"]}`,
			out: float64(123),
		},
		{
			in:  `{"ctx": ["key3", "subkey1"]}`,
			out: float64(1),
		},
		{
			in:  `{"ctx": ["key3", "subkey2"]}`,
			out: []float64{1, 2},
		},
		{
			in:  `{"ctx": ["key3", "subkey2", 1]}`,
			out: float64(2),
		},
		{
			in:  `{"ctx": ["key4"]}`,
			out: true,
		},
		{
			in:  `{"ctx": ["key5"]}`,
			out: nil,
		},
	}

	for _, test := range testCases {
		root, err := Parse(test.in)
		if err != nil {
			t.Errorf("%q got an error: %s", test.in, err.Error())
		} else {
			res, err := evaluator.Evaluate(context, root)
			if err != nil {
				t.Errorf("%q got an error: %s", test.in, err.Error())
			} else {
				if asFloatArray, ok := test.out.([]float64); ok {
					resAsFloatArray := res.([]interface{})
					for i, v := range asFloatArray {
						if v != resAsFloatArray[i] {
							t.Errorf("%q expected %t got %t", test.in, test.out, res)
							break
						}
					}
				} else {
					if res != test.out {
						t.Errorf("%q expected %t got %t", test.in, test.out, res)
					}
				}
			}
		}
	}
}
