package tplinator

const (
	evaluatorExtDepKey = "evaluator"
)

type ExtensionDependencies interface {
	Get(dependencyKey string) interface{}
}

type compoundExtensionDependencies struct {
	extDeps       []ExtensionDependencies
	defaultExtDep ExtensionDependencies
}

func (ed *compoundExtensionDependencies) Get(dependencyKey string) interface{} {
	for _, extDep := range ed.extDeps {
		if dep := extDep.Get(dependencyKey); dep != nil {
			return dep
		}
	}
	return ed.defaultExtDep.Get(dependencyKey)
}

type DefaultExtensionDependencies struct {
	evaluator Evaluator
}

func NewDefaultExtensionDependencies() ExtensionDependencies {
	return &DefaultExtensionDependencies{
		evaluator: &govaluator{},
	}
}

func (ed *DefaultExtensionDependencies) Get(dependencyKey string) interface{} {
	switch dependencyKey {
	case evaluatorExtDepKey:
		return ed.evaluator
	default:
		return nil
	}
}
