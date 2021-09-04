package condition

// NewDefaultEvaluator create a new evaluator with handler for all expressions
// in the expression registry.
func NewDefaultEvaluator() *Evaluator {
	evaluator := NewEvaluator()

	for key, handler := range ExpressionRegistry {
		evaluator.AddHandler(key, handler)
	}

	return evaluator
}
