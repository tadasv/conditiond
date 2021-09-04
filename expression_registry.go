package condition

var ExpressionRegistry map[string]ExpressionFunc

func init() {
	ExpressionRegistry = map[string]ExpressionFunc{
		"and":     AndExpressionHandler,
		"or":      OrExpressionHandler,
		"not":     NotExpressionHandler,
		"if":      IfExpressionHandler,
		"context": ContextExpressionHandler,
		"gt":      GtExpressionHandler,
		"lt":      LtExpressionHandler,
		"gte":     GteExpressionHandler,
		"lte":     LteExpressionHandler,
		"eq":      EqExpressionHandler,
		"sha1mod": Sha1modExpressionHandler,
	}
}
