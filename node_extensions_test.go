package tplinator_test

import (
	"testing"

	"github.com/bmdelacruz/tplinator"
	"golang.org/x/net/html"
)

func TestNodeExtension_Conditional(t *testing.T) {
	extdep := tplinator.NewDefaultExtensionDependencies()
	params := make(map[string]interface{})

	// test if-only conditional element
	h1Node := tplinator.CreateNode(html.ElementNode, "h1", []html.Attribute{
		{Key: "go-if", Val: "shouldRender"},
	}, false)
	tplinator.ConditionalExtensionNodeProcessor(h1Node)

	params["shouldRender"] = true
	finalH1Node, _, err := h1Node.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if finalH1Node == nil {
		t.Errorf("shouldRender is true but finalH1Node is nil")
	}
	params["shouldRender"] = false
	finalH1Node, _, err = h1Node.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if finalH1Node != nil {
		t.Errorf("shouldRender is false but finalH1Node is not nil")
	}
	delete(params, "shouldRender")
	_, _, err = h1Node.ApplyExtensions(extdep, params)
	if err == nil {
		t.Errorf("expecting an error")
	}

	h1Node = tplinator.CreateNode(html.ElementNode, "h1", []html.Attribute{
		{Key: "go-if", Val: "shouldRende]r"},
	}, false)
	tplinator.ConditionalExtensionNodeProcessor(h1Node)

	params["shouldRender"] = false
	_, _, err = h1Node.ApplyExtensions(extdep, params)
	if err == nil {
		t.Errorf("expecting an error")
	}

	// test branching conditional elements
	divNode := tplinator.CreateNode(html.ElementNode, "div", nil, false)
	h1Node = tplinator.CreateNode(html.ElementNode, "h1", []html.Attribute{
		{Key: "go-if", Val: "hasOne"},
	}, false)
	h2Node := tplinator.CreateNode(html.ElementNode, "h2", []html.Attribute{
		{Key: "go-elif", Val: "hasTwo"},
	}, false)
	h3Node := tplinator.CreateNode(html.ElementNode, "h3", []html.Attribute{
		{Key: "go-else-if", Val: "hasThree"},
	}, false)
	h4Node := tplinator.CreateNode(html.ElementNode, "h4", []html.Attribute{
		{Key: "go-else", Val: ""},
	}, false)

	divNode.AppendChild(h1Node)
	divNode.AppendChild(h2Node)
	divNode.AppendChild(h3Node)
	divNode.AppendChild(h4Node)
	tplinator.ConditionalExtensionNodeProcessor(h1Node)

	params["hasOne"] = false
	params["hasTwo"] = true
	params["hasThree"] = false
	finalH1Node, _, err = h1Node.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if finalH1Node != h2Node {
		t.Errorf("finalH1Node should be equal to h2Node")
	}
	params["hasOne"] = false
	params["hasTwo"] = false
	params["hasThree"] = true
	finalH1Node, _, err = h1Node.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if finalH1Node != h3Node {
		t.Errorf("finalH1Node should be equal to h3Node")
	}
	params["hasOne"] = false
	params["hasTwo"] = false
	params["hasThree"] = false
	finalH1Node, _, err = h1Node.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if finalH1Node != h4Node {
		t.Errorf("finalH1Node should be equal to h4Node")
	}

	// test absent else-if and else branch conditional elements
	divNode = tplinator.CreateNode(html.ElementNode, "div", nil, false)
	h1Node = tplinator.CreateNode(html.ElementNode, "h1", []html.Attribute{
		{Key: "go-if", Val: "hasOne"},
	}, false)
	h2Node = tplinator.CreateNode(html.ElementNode, "h2", nil, false)

	divNode.AppendChild(h1Node)
	divNode.AppendChild(h2Node)
	tplinator.ConditionalExtensionNodeProcessor(h1Node)

	params["hasOne"] = false
	finalH1Node, _, err = h1Node.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if finalH1Node != nil {
		t.Errorf("finalH1Node should be nil")
	}
}

func TestNodeExtension_ConditionalClass(t *testing.T) {
	extdep := tplinator.NewDefaultExtensionDependencies()
	params := make(map[string]interface{})

	divNode := tplinator.CreateNode(html.ElementNode, "div", []html.Attribute{
		{Key: "go-if-class-animal", Val: "isAnAnimal"},
	}, false)
	tplinator.ConditionalClassExtensionNodeProcessor(divNode)

	params["isAnAnimal"] = true
	finalDivNode, _, err := divNode.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if finalDivNode == nil {
		t.Errorf("finalDivNode should not be nil")
	} else {
		hasClass, _, classVal := finalDivNode.HasAttribute("class")
		if !hasClass {
			t.Errorf("finalDivNode should have a class attribute")
		} else if classVal != "animal" {
			t.Errorf("finalDivNode's class attribute must have `animal` as its value")
		}
	}

	params["isAnAnimal"] = false
	finalDivNode, _, err = divNode.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if finalDivNode == nil {
		t.Errorf("finalDivNode should not be nil")
	} else {
		hasClass, _, _ := finalDivNode.HasAttribute("class")
		if hasClass {
			t.Errorf("finalDivNode should not have a class attribute")
		}
	}

	delete(params, "isAnAnimal")
	finalDivNode, _, err = divNode.ApplyExtensions(extdep, params)
	if err == nil {
		t.Errorf("expecting an error")
	}

	divNode = tplinator.CreateNode(html.ElementNode, "div", []html.Attribute{
		{Key: "go-if-class-animal", Val: "isAnAnima]l"},
	}, false)
	tplinator.ConditionalClassExtensionNodeProcessor(divNode)

	params["isAnAnimal"] = false
	finalDivNode, _, err = divNode.ApplyExtensions(extdep, params)
	if err == nil {
		t.Errorf("expecting an error")
	}

	divNode = tplinator.CreateNode(html.ElementNode, "div", []html.Attribute{
		{Key: "class", Val: "photo-entry"},
		{Key: "go-if-class-animal", Val: "isAnAnimal"},
	}, false)
	tplinator.ConditionalClassExtensionNodeProcessor(divNode)

	params["isAnAnimal"] = true
	finalDivNode, _, err = divNode.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if finalDivNode == nil {
		t.Errorf("finalDivNode should not be nil")
	} else {
		hasClass, _, classVal := finalDivNode.HasAttribute("class")
		if !hasClass {
			t.Errorf("finalDivNode should have a class attribute")
		} else if classVal != "photo-entry animal" {
			t.Errorf("finalDivNode's class attribute must have `photo-entry animal` as its value")
		}
	}

	params["isAnAnimal"] = false
	finalDivNode, _, err = divNode.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if finalDivNode == nil {
		t.Errorf("finalDivNode should not be nil")
	} else {
		hasClass, _, classVal := finalDivNode.HasAttribute("class")
		if !hasClass {
			t.Errorf("finalDivNode should have a class attribute")
		} else if classVal != "photo-entry" {
			t.Errorf("finalDivNode's class attribute must have `photo-entry` as its value")
		}
	}
}

func TestNodeExtension_Range(t *testing.T) {
	extdep := tplinator.NewDefaultExtensionDependencies()
	params := make(map[string]interface{})

	divNode := tplinator.CreateNode(html.ElementNode, "div", []html.Attribute{
		{Key: "go-range", Val: "pets"},
	}, false)
	tplinator.RangeExtensionNodeProcessor(divNode)

	params["pets"] = tplinator.RangeParams(
		tplinator.EvaluatorParams{
			"name": "catdog",
			"age":  "1",
		},
		tplinator.EvaluatorParams{
			"name": "felycat",
			"age":  "2",
		},
	)
	_, nodes, err := divNode.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("failed to apply range extension. cause: %v", err)
		return
	} else if len(nodes) != 2 {
		t.Errorf("expected 2 new nodes")
		return
	}

	divNode = tplinator.CreateNode(html.ElementNode, "div", []html.Attribute{
		{Key: "go-range", Val: "pets"},
	}, false)
	tplinator.RangeExtensionNodeProcessor(divNode)

	params["pets"] = tplinator.EvaluatorParams{
		"name": "catdog",
		"age":  "1",
	}
	_, _, err = divNode.ApplyExtensions(extdep, params)
	if err == nil {
		t.Error("expecting an error")
	}

	divNode = tplinator.CreateNode(html.ElementNode, "div", []html.Attribute{
		{Key: "go-range", Val: "pet]s"},
	}, false)
	tplinator.RangeExtensionNodeProcessor(divNode)

	params["pets"] = tplinator.RangeParams(
		tplinator.EvaluatorParams{
			"name": "catdog",
			"age":  "1",
		},
		tplinator.EvaluatorParams{
			"name": "felycat",
			"age":  "2",
		},
	)
	_, _, err = divNode.ApplyExtensions(extdep, params)
	if err == nil {
		t.Error("expecting an error")
	}

	divNode = tplinator.CreateNode(html.ElementNode, "div", []html.Attribute{
		{Key: "go-range", Val: "pets"},
		{Key: "value", Val: "{{go:name}}"},
	}, false)

	tplinator.RangeExtensionNodeProcessor(divNode)
	tplinator.StringInterpolationNodeProcessor(divNode)

	params["pets"] = tplinator.RangeParams(
		tplinator.EvaluatorParams{
			"age": "1",
		},
		tplinator.EvaluatorParams{
			"age": "2",
		},
	)
	_, _, err = divNode.ApplyExtensions(extdep, params)
	if err == nil {
		t.Error("expecting an error")
	}
}

func TestNodeExtension_StringInterpolation(t *testing.T) {
	extdep := tplinator.NewDefaultExtensionDependencies()
	params := make(map[string]interface{})

	textNode := tplinator.CreateNode(html.TextNode, "{{go:description}}", nil, false)
	tplinator.StringInterpolationNodeProcessor(textNode)

	params["description"] = "Lorem ipsum dolor sit amet"
	newTextNode, _, err := textNode.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if newTextNode == nil {
		t.Errorf("newTextNode should not be nil")
	} else if newTextNode.Data != "Lorem ipsum dolor sit amet" {
		t.Errorf("wanted `Lorem ipsum dolor sit amet`, got `%v`", newTextNode.Data)
	}

	h1Node := tplinator.CreateNode(html.ElementNode, "h1", nil, false)
	h1Node.SetContextParams(tplinator.EvaluatorParams{
		"description": "The big brown fox",
	})
	textNode = tplinator.CreateNode(html.TextNode, "{{go:description}}", nil, false)
	h1Node.AppendChild(textNode)

	tplinator.StringInterpolationNodeProcessor(h1Node)
	tplinator.StringInterpolationNodeProcessor(textNode)

	params["description"] = "Lorem ipsum dolor sit amet"
	newTextNode, _, err = textNode.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if newTextNode == nil {
		t.Errorf("newTextNode should not be nil")
	} else if newTextNode.Data != "The big brown fox" {
		t.Errorf("wanted `The big brown fox`, got `%v`", newTextNode.Data)
	}

	aNode := tplinator.CreateNode(html.ElementNode, "h1", []html.Attribute{
		{Key: "href", Val: "/users/{{go:uid}}"},
	}, false)
	tplinator.StringInterpolationNodeProcessor(aNode)

	params["uid"] = "10000284736283"
	newANode, _, err := aNode.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if newANode == nil {
		t.Errorf("newANode should not be nil")
	} else if hasHref, _, hrefVal := newANode.HasAttribute("href"); hasHref {
		if hrefVal != "/users/10000284736283" {
			t.Errorf("wanted `/users/10000284736283`, got `%v`", hrefVal)
		}
	} else {
		t.Errorf("newANode should have an href attribute")
	}

	aNode = tplinator.CreateNode(html.ElementNode, "h1", []html.Attribute{
		{Key: "href", Val: "/users/{{go:uid}}"},
	}, false)
	aNode.SetContextParams(tplinator.EvaluatorParams{
		"uid": "41728897352322",
	})
	tplinator.StringInterpolationNodeProcessor(aNode)

	params["uid"] = "10000284736283"
	newANode, _, err = aNode.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if newANode == nil {
		t.Errorf("newANode should not be nil")
	} else if hasHref, _, hrefVal := newANode.HasAttribute("href"); hasHref {
		if hrefVal != "/users/41728897352322" {
			t.Errorf("wanted `/users/41728897352322`, got `%v`", hrefVal)
		}
	} else {
		t.Errorf("newANode should have an href attribute")
	}
}
