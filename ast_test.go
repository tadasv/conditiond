package condition

import (
	"testing"
)

func TestAST(t *testing.T) {
	expected := `FUNCTION<functionB>
 \_ ARRAY
    |__ LITERAL<number::555>
    |__ LITERAL<number::66.6>
    |__ LITERAL<null>
    |__ LITERAL<string::value1>
    |__ FUNCTION<functionC>
    |    \_ FUNCTION<functionD>
    |        \_ ARRAY
    |           |__ LITERAL<string::value>
    |            \_ LITERAL<number::55>
    |__ LITERAL<bool::false>
     \_ FUNCTION<functionA>
         \_ LITERAL<string::input>
`

	root, err := Parse(`
{
	"functionB":
	[
		555,
		66.6,
		null,
		"value1",
		{"functionC": {"functionD": ["value", 55] } },
		false,
		{"functionA": "input"}
	]
}
	`)

	if err != nil {
		t.Errorf("%s\n", err.Error())
	} else {
		res := Stringify(root)
		if res != expected {
			t.Errorf("Expected: \n%s\n\nGot: \n%s\n\n", expected, res)
		}
	}
}
