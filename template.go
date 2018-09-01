package tplinator

type PrecompiledTemplate struct {
	documentNode *node
}

func (pt *PrecompiledTemplate) Execute(data map[string]interface{}) ([]byte, error) {
	docStr, err := pt.documentNode.Execute(
		createEvaluatorFunc(data),
		createBoolEvaluatorFunc(data),
	)
	if err != nil {
		return nil, err
	}
	return []byte(docStr), nil
}

func (pt *PrecompiledTemplate) ExecuteStrict(data map[string]interface{}) ([]byte, error) {
	docStr, err := pt.documentNode.Execute(
		createStrictEvaluatorFunc(data),
		createStrictBoolEvaluatorFunc(data),
	)
	if err != nil {
		return nil, err
	}
	return []byte(docStr), nil
}

type evaluator func(string) (string, error)
type boolEvaluator func(string) (bool, error)

func createEvaluatorFunc(params map[string]interface{}) evaluator {
	return func(data string) (string, error) {
		return data, nil // TODO
	}
}

func createStrictEvaluatorFunc(params map[string]interface{}) evaluator {
	return func(data string) (string, error) {
		return data, nil // TODO
	}
}

func createBoolEvaluatorFunc(params map[string]interface{}) boolEvaluator {
	return func(data string) (bool, error) {
		return false, nil // TODO
	}
}

func createStrictBoolEvaluatorFunc(params map[string]interface{}) boolEvaluator {
	return func(data string) (bool, error) {
		return false, nil // TODO
	}
}
