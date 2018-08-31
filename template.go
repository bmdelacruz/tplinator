package tplinator

type PrecompiledTemplate struct {
	documentNode *node
}

func (pt *PrecompiledTemplate) Execute(data map[string]interface{}) ([]byte, error) {
	return nil, nil // TODO
}

func (pt *PrecompiledTemplate) ExecuteStrict(data map[string]interface{}) ([]byte, error) {
	return nil, nil // TODO
}

type evaluator func(string) (string, error)

func evaluatorFunc(string) (string, error) {
	return "", nil // TODO
}

func strictEvaluatorFunc(string) (string, error) {
	return "", nil // TODO
}
