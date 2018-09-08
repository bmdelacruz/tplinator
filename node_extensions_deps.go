package tplinator

type DependencyKey string

const (
	EvaluatorExtDepKey DependencyKey = "evaluator"
)

type ExtensionDependencies interface {
	Get(dependencyKey DependencyKey) interface{}
}

type compoundExtensionDependencies struct {
	extDeps       []ExtensionDependencies
	defaultExtDep ExtensionDependencies
}

func (ed *compoundExtensionDependencies) Get(dependencyKey DependencyKey) interface{} {
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

func (ed *DefaultExtensionDependencies) Get(dependencyKey DependencyKey) interface{} {
	switch dependencyKey {
	case EvaluatorExtDepKey:
		return ed.evaluator
	default:
		return nil
	}
}
