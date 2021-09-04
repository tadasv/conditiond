package condition

// NewDefaultEvaluator create a new evaluator with handler for all expressions
// in the expression registry.
func NewDefaultEvaluator() *Evaluator {
	evaluator := NewEvaluator()

	for key, handler := range ExpressionRegistry {
		evaluator.AddHandler(key, handler)
	}

	/*
		evaluator.AddHandler("or", OrExpressionHandler)
		evaluator.AddHandler("not", NotExpressionHandler)
		evaluator.AddHandler("if", IfExpressionHandler)
		evaluator.AddHandler("context", ContextExpressionHandler)
		evaluator.AddHandler("gt", GtExpressionHandler)
		evaluator.AddHandler("lt", LtExpressionHandler)
		evaluator.AddHandler("gte", GteExpressionHandler)
		evaluator.AddHandler("lte", LteExpressionHandler)
		evaluator.AddHandler("eq", EqExpressionHandler)
		evaluator.AddHandler("sha1mod", Sha1modExpressionHandler)
	*/
	return evaluator
}
