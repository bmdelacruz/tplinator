package tplinator_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bmdelacruz/tplinator"
)

func TestDefaultExtensionDependencies(t *testing.T) {
	deps := tplinator.NewDefaultExtensionDependencies()

	evaluator, isEvaluator := deps.Get(tplinator.EvaluatorExtDepKey).(tplinator.Evaluator)
	if evaluator == nil || !isEvaluator {
		t.Error("expected to get an evaluator")
	}
	modifier := deps.Get("modifier")
	if modifier != nil {
		t.Error("expected modifier to be nil")
	}
}

func TestCompoundExtensionDependencies(t *testing.T) {
	tpl, err := tplinator.Tplinate(
		strings.NewReader(`<div></div>`),
		tplinator.NodeProcessorsParserOption(
			func(node *tplinator.Node) {
				node.AddExtension(&someExt{})
			},
		),
	)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	tpl.AddExtensionDependencies(&someExtDef{})

	err = tpl.Render(tplinator.EvaluatorParams{}, func(_ string) {})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

type someExt struct {
}

func (se *someExt) Apply(
	node *tplinator.Node,
	dependencies tplinator.ExtensionDependencies,
	params tplinator.EvaluatorParams,
) (*tplinator.Node, []*tplinator.Node, error) {
	someDependency, isBool := dependencies.Get("someDependency").(bool)
	if !isBool || !someDependency {
		return nil, nil, fmt.Errorf("failed to resolve dependency `someDependency`")
	}
	return node, nil, nil
}

type someExtDef struct {
}

func (ed *someExtDef) Get(depKey tplinator.DependencyKey) interface{} {
	if depKey == "someDependency" {
		return true
	}
	return nil
}
