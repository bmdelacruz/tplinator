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
		html.Attribute{Key: "go-if", Val: "shouldRender"},
	}, false)
	tplinator.ConditionalExtensionNodeProcessor(h1Node)

	params["shouldRender"] = true
	finalH1Node, err := h1Node.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if finalH1Node == nil {
		t.Errorf("shouldRender is true but finalH1Node is nil")
	}
	params["shouldRender"] = false
	finalH1Node, err = h1Node.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if finalH1Node != nil {
		t.Errorf("shouldRender is false but finalH1Node is not nil")
	}

	// test branching conditional elements
	divNode := tplinator.CreateNode(html.ElementNode, "div", nil, false)
	h1Node = tplinator.CreateNode(html.ElementNode, "h1", []html.Attribute{
		html.Attribute{Key: "go-if", Val: "hasOne"},
	}, false)
	h2Node := tplinator.CreateNode(html.ElementNode, "h2", []html.Attribute{
		html.Attribute{Key: "go-elif", Val: "hasTwo"},
	}, false)
	h3Node := tplinator.CreateNode(html.ElementNode, "h3", []html.Attribute{
		html.Attribute{Key: "go-else-if", Val: "hasThree"},
	}, false)
	h4Node := tplinator.CreateNode(html.ElementNode, "h4", []html.Attribute{
		html.Attribute{Key: "go-else", Val: ""},
	}, false)

	divNode.AppendChild(h1Node)
	divNode.AppendChild(h2Node)
	divNode.AppendChild(h3Node)
	divNode.AppendChild(h4Node)
	tplinator.ConditionalExtensionNodeProcessor(h1Node)

	params["hasOne"] = false
	params["hasTwo"] = true
	params["hasThree"] = false
	finalH1Node, err = h1Node.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if finalH1Node != h2Node {
		t.Errorf("finalH1Node should be equal to h2Node")
	}
	params["hasOne"] = false
	params["hasTwo"] = false
	params["hasThree"] = true
	finalH1Node, err = h1Node.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if finalH1Node != h3Node {
		t.Errorf("finalH1Node should be equal to h3Node")
	}
	params["hasOne"] = false
	params["hasTwo"] = false
	params["hasThree"] = false
	finalH1Node, err = h1Node.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if finalH1Node != h4Node {
		t.Errorf("finalH1Node should be equal to h4Node")
	}

	// test absent else-if and else branch conditional elements
	divNode = tplinator.CreateNode(html.ElementNode, "div", nil, false)
	h1Node = tplinator.CreateNode(html.ElementNode, "h1", []html.Attribute{
		html.Attribute{Key: "go-if", Val: "hasOne"},
	}, false)
	h2Node = tplinator.CreateNode(html.ElementNode, "h2", nil, false)

	divNode.AppendChild(h1Node)
	divNode.AppendChild(h2Node)
	tplinator.ConditionalExtensionNodeProcessor(h1Node)

	params["hasOne"] = false
	finalH1Node, err = h1Node.ApplyExtensions(extdep, params)
	if err != nil {
		t.Errorf("encountered an unexpected error")
	} else if finalH1Node != nil {
		t.Errorf("finalH1Node should be nil")
	}
}
