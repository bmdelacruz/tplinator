package tplinator_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/bmdelacruz/tplinator"
)

func TestCreateTemplateFromReader(t *testing.T) {
	_, err := tplinator.CreateTemplateFromReader(
		strings.NewReader(`<div></div>`),
	)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	_, err = tplinator.CreateTemplateFromReader(
		strings.NewReader(`<div></div`),
	)
	if err == nil {
		t.Error("unexpecting an error")
	}
}

func TestTemplate_RenderBytes(t *testing.T) {
	tpl, err := tplinator.CreateTemplateFromReader(
		strings.NewReader(`<div>{{go:number}}</div>`),
		tplinator.NodeProcessorsParserOption(
			tplinator.StringInterpolationNodeProcessor,
		),
	)
	if err != nil {
		t.Error("unexpected error:", err)
		return
	}
	tpl.AddExtensionDependencies(tplinator.NewDefaultExtensionDependencies())

	actualBytes, err := tpl.RenderBytes(tplinator.EvaluatorParams{
		"number": "1",
	})
	if err != nil {
		t.Error("unexpected error:", err)
		return
	}
	expectedBytes := []byte(`<div>1</div>`)
	if !reflect.DeepEqual(actualBytes, expectedBytes) {
		t.Error("unexpected result. actual:", actualBytes,
			"expected:", expectedBytes)
	}

	_, err = tpl.RenderBytes(tplinator.EvaluatorParams{})
	if err == nil {
		t.Error("expecting an error due to evaluator")
	} else {
		t.Log("expecting error due to evaluator:", err)
	}
}

func TestTemplate_RenderString(t *testing.T) {
	tpl, err := tplinator.CreateTemplateFromReader(
		strings.NewReader(`<div>{{go:number}}</div>`),
		tplinator.NodeProcessorsParserOption(
			tplinator.StringInterpolationNodeProcessor,
		),
	)
	if err != nil {
		t.Error("unexpected error:", err)
		return
	}
	tpl.AddExtensionDependencies(tplinator.NewDefaultExtensionDependencies())

	actualString, err := tpl.RenderString(tplinator.EvaluatorParams{
		"number": "1",
	})
	if err != nil {
		t.Error("unexpected error:", err)
		return
	}

	expectedString := `<div>1</div>`
	if !reflect.DeepEqual(actualString, expectedString) {
		t.Error("unexpected result. actual:", actualString,
			"expected:", expectedString)
	}

	_, err = tpl.RenderString(tplinator.EvaluatorParams{})
	if err == nil {
		t.Error("expecting an error due to evaluator")
	} else {
		t.Log("expecting error due to evaluator:", err)
	}
}
