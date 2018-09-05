package tplinator

type Evaluator interface {
	EvaluateBool(input string) (bool, error)
	EvaluateString(input string) (string, error)
	Evaluate(input string) (interface{}, error)
}
