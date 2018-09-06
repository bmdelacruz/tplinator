package tplinator_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bmdelacruz/tplinator"
	"golang.org/x/net/html"
)

func TestCreateNode(t *testing.T) {
	testCases := []struct {
		name string

		// arguments
		data          string
		isSelfClosing bool
		nodeType      html.NodeType
		attributes    []html.Attribute

		// expected tags
		expectedStartTag string
		expectedEndTag   string
	}{
		{
			name: "doctype node test",

			data:          "html",
			isSelfClosing: false,
			nodeType:      html.DoctypeNode,
			attributes:    nil,

			expectedStartTag: "<!DOCTYPE html>",
			expectedEndTag:   "",
		},
		{
			name: "text node test",

			data:          "Hello, world!",
			isSelfClosing: false,
			nodeType:      html.TextNode,
			attributes:    nil,

			expectedStartTag: "Hello, world!",
			expectedEndTag:   "",
		},
		{
			name: "element (div) node test",

			data:          "div",
			isSelfClosing: false,
			nodeType:      html.ElementNode,
			attributes: []html.Attribute{
				{Key: "class", Val: "container"},
				{Key: "hidden", Val: ""},
			},

			expectedStartTag: `<div class="container" hidden>`,
			expectedEndTag:   "</div>",
		},
		{
			name: "self-closing element (img) node test",

			data:          "img",
			isSelfClosing: true,
			nodeType:      html.ElementNode,
			attributes: []html.Attribute{
				{Key: "src", Val: "/static/images/cat.png"},
			},

			expectedStartTag: `<img src="/static/images/cat.png"/>`,
			expectedEndTag:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node := tplinator.CreateNode(tc.nodeType, tc.data, tc.attributes, tc.isSelfClosing)

			startTag, endTag := node.Tags()
			if startTag != tc.expectedStartTag || endTag != tc.expectedEndTag {
				t.Errorf("unexpected tag/s. startTag: '%v' endTag: '%v'. test case: %+v", startTag, endTag, tc)
			}
		})
	}
}

func TestNode_Parent(t *testing.T) {
	divNode := tplinator.CreateNode(html.ElementNode, "div", nil, false)
	formNode := tplinator.CreateNode(html.ElementNode, "form", nil, false)

	divNode.AppendChild(formNode)

	if formNode.Parent() != divNode {
		t.Errorf("divNode should be the parent of formNode")
	}
}

func TestNode_FirstChild(t *testing.T) {
	divNode := tplinator.CreateNode(html.ElementNode, "div", nil, false)
	formNode := tplinator.CreateNode(html.ElementNode, "form", nil, false)
	imgNode := tplinator.CreateNode(html.ElementNode, "img", nil, true)

	divNode.AppendChild(formNode)
	divNode.AppendChild(imgNode)

	if divNode.FirstChild() != formNode {
		t.Errorf("formNode should be the first child of divNode")
	}
}

func TestNode_LastChild(t *testing.T) {
	divNode := tplinator.CreateNode(html.ElementNode, "div", nil, false)
	formNode := tplinator.CreateNode(html.ElementNode, "form", nil, false)
	imgNode := tplinator.CreateNode(html.ElementNode, "img", nil, true)

	divNode.AppendChild(formNode)
	divNode.AppendChild(imgNode)

	if divNode.LastChild() != imgNode {
		t.Errorf("formNode should be the last child of imgNode")
	}
}

func TestNode_NextSibling(t *testing.T) {
	divNode := tplinator.CreateNode(html.ElementNode, "div", nil, false)
	formNode := tplinator.CreateNode(html.ElementNode, "form", nil, false)
	imgNode := tplinator.CreateNode(html.ElementNode, "img", nil, true)
	h1Node := tplinator.CreateNode(html.ElementNode, "h1", nil, false)

	divNode.AppendChild(formNode)
	divNode.AppendChild(imgNode)
	divNode.AppendChild(h1Node)

	if imgNode.NextSibling() != h1Node {
		t.Errorf("h1Node should be the next sibling of imgNode")
	}
}

func TestNode_PreviousSibling(t *testing.T) {
	divNode := tplinator.CreateNode(html.ElementNode, "div", nil, false)
	formNode := tplinator.CreateNode(html.ElementNode, "form", nil, false)
	imgNode := tplinator.CreateNode(html.ElementNode, "img", nil, true)
	h1Node := tplinator.CreateNode(html.ElementNode, "h1", nil, false)

	divNode.AppendChild(formNode)
	divNode.AppendChild(imgNode)
	divNode.AppendChild(h1Node)

	if imgNode.PreviousSibling() != formNode {
		t.Errorf("formNode should be the next sibling of imgNode")
	}
}

func TestNode_Children(t *testing.T) {
	childNodes := []*tplinator.Node{
		tplinator.CreateNode(html.ElementNode, "img", nil, true),
		tplinator.CreateNode(html.ElementNode, "form", nil, false),
	}

	node := tplinator.CreateNode(html.ElementNode, "div", nil, false)
	for _, childNode := range childNodes {
		node.AppendChild(childNode)
	}

	actualChildNodes := make([]*tplinator.Node, 0)
	node.Children(func(_ int, childNode *tplinator.Node) bool {
		actualChildNodes = append(actualChildNodes, childNode)
		return true
	})

	if len(childNodes) != len(actualChildNodes) {
		t.Error("actual and expected child nodes does not have the same count")
	} else {
		for i := 0; i < len(childNodes); i++ {
			if childNodes[i] != actualChildNodes[i] {
				t.Error("the order of the actual and expected child nodes are not the same")
			}
		}
	}
}

func TestNode_ChildrenBreak(t *testing.T) {
	imgNode := tplinator.CreateNode(html.ElementNode, "img", nil, true)
	formNode := tplinator.CreateNode(html.ElementNode, "form", nil, false)
	h1Node := tplinator.CreateNode(html.ElementNode, "h1", nil, false)

	childNodes := []*tplinator.Node{imgNode, formNode, h1Node}

	node := tplinator.CreateNode(html.ElementNode, "div", nil, false)
	for _, childNode := range childNodes {
		node.AppendChild(childNode)
	}

	expectedChildNodes := []*tplinator.Node{imgNode, formNode}
	actualChildNodes := make([]*tplinator.Node, 0)
	node.Children(func(_ int, childNode *tplinator.Node) bool {
		actualChildNodes = append(actualChildNodes, childNode)
		return childNode != formNode
	})

	if len(expectedChildNodes) != len(actualChildNodes) {
		t.Errorf("actual and expected child nodes does not have the same count. want %+v got %+v", expectedChildNodes, actualChildNodes)
	} else {
		for i := 0; i < len(expectedChildNodes); i++ {
			if expectedChildNodes[i] != actualChildNodes[i] {
				t.Errorf("the order of the actual and expected child nodes are not the same. want %+v got %+v", expectedChildNodes, actualChildNodes)
			}
		}
	}
}

func TestNode_NextSiblings(t *testing.T) {
	imgNode := tplinator.CreateNode(html.ElementNode, "img", nil, true)
	formNode := tplinator.CreateNode(html.ElementNode, "form", nil, false)
	h1Node := tplinator.CreateNode(html.ElementNode, "h1", nil, false)

	childNodes := []*tplinator.Node{imgNode, formNode, h1Node}
	siblingNodes := []*tplinator.Node{formNode, h1Node}

	node := tplinator.CreateNode(html.ElementNode, "div", nil, false)
	for _, childNode := range childNodes {
		node.AppendChild(childNode)
	}

	actualSiblingNodes := make([]*tplinator.Node, 0)
	imgNode.NextSiblings(func(sibling *tplinator.Node) bool {
		actualSiblingNodes = append(actualSiblingNodes, sibling)
		return true
	})

	if len(siblingNodes) != len(actualSiblingNodes) {
		t.Error("actual and expected sibling nodes does not have the same count")
	} else {
		for i := 0; i < len(siblingNodes); i++ {
			if siblingNodes[i] != actualSiblingNodes[i] {
				t.Error("the order of the actual and expected child sibling are not the same")
			}
		}
	}
}

func TestNode_NextSiblingsBreak(t *testing.T) {
	imgNode := tplinator.CreateNode(html.ElementNode, "img", nil, true)
	formNode := tplinator.CreateNode(html.ElementNode, "form", nil, false)
	h1Node := tplinator.CreateNode(html.ElementNode, "h1", nil, false)

	childNodes := []*tplinator.Node{imgNode, formNode, h1Node}
	siblingNodes := []*tplinator.Node{formNode}

	node := tplinator.CreateNode(html.ElementNode, "div", nil, false)
	for _, childNode := range childNodes {
		node.AppendChild(childNode)
	}

	actualSiblingNodes := make([]*tplinator.Node, 0)
	imgNode.NextSiblings(func(sibling *tplinator.Node) bool {
		actualSiblingNodes = append(actualSiblingNodes, sibling)
		return sibling != formNode
	})

	if len(siblingNodes) != len(actualSiblingNodes) {
		t.Error("actual and expected sibling nodes does not have the same count")
	} else {
		for i := 0; i < len(siblingNodes); i++ {
			if siblingNodes[i] != actualSiblingNodes[i] {
				t.Error("the order of the actual and expected child sibling are not the same")
			}
		}
	}
}

func TestNode_Attributes(t *testing.T) {
	node := tplinator.CreateNode(html.ElementNode, "img", []html.Attribute{
		html.Attribute{Key: "src", Val: "/static/images/cat.png"},
		html.Attribute{Key: "class", Val: "pictures animal"},
		html.Attribute{Key: "hidden", Val: ""},
	}, true)

	attributes := []tplinator.Attribute{
		tplinator.Attribute{Key: "src", Value: "/static/images/cat.png"},
		tplinator.Attribute{Key: "class", Value: "pictures animal"},
		tplinator.Attribute{Key: "hidden", Value: "", KeyOnly: true},
	}
	actualAttributes := node.Attributes()

	if len(actualAttributes) != len(attributes) {
		t.Errorf("actual and expected attributes does not have the same length. actual: %+v, expected: %+v", actualAttributes, attributes)
	} else {
		for attrIdx, attr := range node.Attributes() {
			expectedAttribute := attributes[attrIdx]
			if attr.Key != expectedAttribute.Key || attr.Value != expectedAttribute.Value || attr.KeyOnly != expectedAttribute.KeyOnly {
				t.Errorf("actual and expected attribute does not have the same data. actual: %+v, expected: %+v", attr, expectedAttribute)
			}
		}
	}
}

func TestNode_HasAttribute(t *testing.T) {
	testCases := []struct {
		name string

		// arguments
		key string

		// expected values
		hasAttribute   bool
		attributeIndex int
		attributeValue string
	}{
		{
			name: "present attribute (src)",

			key: "src",

			hasAttribute:   true,
			attributeIndex: 0,
			attributeValue: "/static/images/cat.png",
		},
		{
			name: "absent attribute (id)",

			key: "id",

			hasAttribute:   false,
			attributeIndex: -1,
			attributeValue: "",
		},
	}

	node := tplinator.CreateNode(html.ElementNode, "img", []html.Attribute{
		html.Attribute{Key: "src", Val: "/static/images/cat.png"},
	}, true)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hasAttribute, attributeIndex, attributeValue := node.HasAttribute(tc.key)
			if hasAttribute != tc.hasAttribute ||
				attributeIndex != tc.attributeIndex ||
				attributeValue != tc.attributeValue {

				t.Errorf("attribute assertion failed. result: (hasAttribute=%v, index=%v, value=%v) test case: %+v",
					hasAttribute, attributeIndex, attributeValue, tc)
			}
		})
	}
}

func TestNode_HasAttributes(t *testing.T) {
	testCases := []struct {
		name string

		// arguments
		nodeMaker func() *tplinator.Node
		testFunc  func(tplinator.Attribute) bool

		// expected values
		checkerFunc func([]tplinator.Attribute) error
	}{
		{
			name: "present attributes",

			nodeMaker: func() *tplinator.Node {
				return tplinator.CreateNode(html.ElementNode, "img", []html.Attribute{
					html.Attribute{Key: "id", Val: "#animal002"},
					html.Attribute{Key: "src", Val: "/static/images/cat.png"},
					html.Attribute{Key: "go-if", Val: "hasImage"},
					html.Attribute{Key: "go-if-class-cat", Val: "isCat"},
				}, true)
			},
			testFunc: func(attr tplinator.Attribute) bool {
				return attr.Key == "src" || attr.Key == "go-if" ||
					strings.HasPrefix(attr.Key, "go-if-class-")
			},

			checkerFunc: func(attrs []tplinator.Attribute) error {
				if len(attrs) != 3 {
					return fmt.Errorf("expecting 3 matching attributes")
				} else if attrs[0].Key != "src" {
					return fmt.Errorf("expecting `src` as 1st matching attribute")
				} else if attrs[1].Key != "go-if" {
					return fmt.Errorf("expecting `go-if` as 2nd matching attribute")
				} else if attrs[2].Key != "go-if-class-cat" {
					return fmt.Errorf("expecting `go-if-class-cat` as 3rd matching attribute")
				}
				return nil
			},
		},
	}

	tc := testCases[0]
	matches := tc.nodeMaker().HasAttributes(tc.testFunc)
	err := tc.checkerFunc(matches)
	if err != nil {
		t.Error(err, matches)
	}
}

func TestNode_AddAttribute(t *testing.T) {
	testCases := []struct {
		name string

		modifyNodeFunc func(*tplinator.Node)

		// arguments
		key   string
		value string

		// expected values
		hasAttribute   bool
		attributeIndex int
		attributeValue string
	}{
		{
			name: "add absent attribute (id)",

			modifyNodeFunc: func(node *tplinator.Node) {
				node.AddAttribute("id", "cat_tubby")
			},

			key:   "id",
			value: "cat_tubby",

			hasAttribute:   true,
			attributeIndex: 1,
			attributeValue: "cat_tubby",
		},
		{
			name: "add present attribute (src)",

			modifyNodeFunc: func(node *tplinator.Node) {
				node.AddAttribute("src", "/static/images/dog.png")
			},

			key:   "src",
			value: "/static/images/dog.png",

			hasAttribute:   true,
			attributeIndex: 0,
			attributeValue: "/static/images/dog.png",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node := tplinator.CreateNode(html.ElementNode, "img", []html.Attribute{
				html.Attribute{Key: "src", Val: "/static/images/cat.png"},
			}, true)

			tc.modifyNodeFunc(node)

			hasAttribute, attributeIndex, attributeValue := node.HasAttribute(tc.key)
			if hasAttribute != tc.hasAttribute ||
				attributeIndex != tc.attributeIndex ||
				attributeValue != tc.attributeValue {

				t.Errorf("attribute assertion failed. result: (hasAttribute=%v, index=%v, value=%v) test case: %+v",
					hasAttribute, attributeIndex, attributeValue, tc)
			}
		})
	}
}

func TestNode_RemoveAttribute(t *testing.T) {
	testCases := []struct {
		name string

		modifyNodeFunc func(*tplinator.Node)

		// arguments
		key string

		// expected values
		hasAttribute   bool
		attributeIndex int
		attributeValue string
	}{
		{
			name: "remove absent attribute (id)",

			modifyNodeFunc: func(node *tplinator.Node) {
				node.RemoveAttribute("id")
			},

			key: "id",

			hasAttribute:   false,
			attributeIndex: -1,
			attributeValue: "",
		},
		{
			name: "remove present attribute (src)",

			modifyNodeFunc: func(node *tplinator.Node) {
				node.RemoveAttribute("src")
			},

			key: "src",

			hasAttribute:   false,
			attributeIndex: -1,
			attributeValue: "",
		},
		{
			name: "remove present attribute (src) then check attribute (class)",

			modifyNodeFunc: func(node *tplinator.Node) {
				node.RemoveAttribute("src")
			},

			key: "class",

			hasAttribute:   true,
			attributeIndex: 0,
			attributeValue: "pictures animal",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node := tplinator.CreateNode(html.ElementNode, "img", []html.Attribute{
				html.Attribute{Key: "src", Val: "/static/images/cat.png"},
				html.Attribute{Key: "class", Val: "pictures animal"},
			}, true)

			tc.modifyNodeFunc(node)

			hasAttribute, attributeIndex, attributeValue := node.HasAttribute(tc.key)
			if hasAttribute != tc.hasAttribute ||
				attributeIndex != tc.attributeIndex ||
				attributeValue != tc.attributeValue {

				t.Errorf("attribute assertion failed. result: (hasAttribute=%v, index=%v, value=%v) test case: %+v",
					hasAttribute, attributeIndex, attributeValue, tc)
			}
		})
	}
}

func TestNode_Insert(t *testing.T) {
	testCases := []struct {
		name     string
		nodeFunc func(imgNode, formNode, h1Node *tplinator.Node) ([]*tplinator.Node, func(*tplinator.Node))
	}{
		{
			name: "h1-img-form order",
			nodeFunc: func(imgNode, formNode, h1Node *tplinator.Node) ([]*tplinator.Node, func(*tplinator.Node)) {
				return []*tplinator.Node{h1Node, imgNode, formNode}, func(node *tplinator.Node) {
					node.AppendChild(imgNode)
					node.AppendChild(formNode)
					node.Insert(h1Node, imgNode)
				}
			},
		},
		{
			name: "img-h1-form order",
			nodeFunc: func(imgNode, formNode, h1Node *tplinator.Node) ([]*tplinator.Node, func(*tplinator.Node)) {
				return []*tplinator.Node{imgNode, h1Node, formNode}, func(node *tplinator.Node) {
					node.AppendChild(imgNode)
					node.AppendChild(formNode)
					node.Insert(h1Node, formNode)
				}
			},
		},
		{
			name: "img-form-h1 order",
			nodeFunc: func(imgNode, formNode, h1Node *tplinator.Node) ([]*tplinator.Node, func(*tplinator.Node)) {
				return []*tplinator.Node{imgNode, formNode, h1Node}, func(node *tplinator.Node) {
					node.AppendChild(imgNode)
					node.AppendChild(formNode)
					node.Insert(h1Node, nil)
				}
			},
		},
	}
	testInsert := func(t *testing.T, modifyNodeFunc func(*tplinator.Node), childNodes []*tplinator.Node) {
		node := tplinator.CreateNode(html.ElementNode, "div", nil, false)
		modifyNodeFunc(node)

		actualChildNodes := make([]*tplinator.Node, 0)
		node.Children(func(_ int, childNode *tplinator.Node) bool {
			actualChildNodes = append(actualChildNodes, childNode)
			return true
		})

		if len(childNodes) != len(actualChildNodes) {
			t.Error("actual and expected child nodes does not have the same count")
		}
		for i := 0; i < len(childNodes); i++ {
			if childNodes[i] != actualChildNodes[i] {
				t.Error("the order of the actual and expected child nodes are not the same")
			}
		}
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			imgNode := tplinator.CreateNode(html.ElementNode, "img", nil, true)
			formNode := tplinator.CreateNode(html.ElementNode, "form", nil, false)
			h1Node := tplinator.CreateNode(html.ElementNode, "h1", nil, false)

			childNodes, modifyNodeFunc := tc.nodeFunc(imgNode, formNode, h1Node)
			testInsert(t, modifyNodeFunc, childNodes)
		})
	}
}

func TestNode_RemoveChild(t *testing.T) {
	testCases := []struct {
		name     string
		nodeFunc func(imgNode, formNode, h1Node *tplinator.Node) ([]*tplinator.Node, func(*tplinator.Node))
	}{
		{
			name: "img-h1 order",
			nodeFunc: func(imgNode, formNode, h1Node *tplinator.Node) ([]*tplinator.Node, func(*tplinator.Node)) {
				return []*tplinator.Node{imgNode, h1Node}, func(node *tplinator.Node) {
					node.AppendChild(imgNode)
					node.AppendChild(formNode)
					node.AppendChild(h1Node)
					node.RemoveChild(formNode)
				}
			},
		},
		{
			name: "form-h1 order",
			nodeFunc: func(imgNode, formNode, h1Node *tplinator.Node) ([]*tplinator.Node, func(*tplinator.Node)) {
				return []*tplinator.Node{formNode, h1Node}, func(node *tplinator.Node) {
					node.AppendChild(imgNode)
					node.AppendChild(formNode)
					node.AppendChild(h1Node)
					node.RemoveChild(imgNode)
				}
			},
		},
		{
			name: "img-form order",
			nodeFunc: func(imgNode, formNode, h1Node *tplinator.Node) ([]*tplinator.Node, func(*tplinator.Node)) {
				return []*tplinator.Node{imgNode, formNode}, func(node *tplinator.Node) {
					node.AppendChild(imgNode)
					node.AppendChild(formNode)
					node.AppendChild(h1Node)
					node.RemoveChild(h1Node)
				}
			},
		},
	}
	testRemove := func(t *testing.T, modifyNodeFunc func(*tplinator.Node), childNodes []*tplinator.Node) {
		node := tplinator.CreateNode(html.ElementNode, "div", nil, false)
		modifyNodeFunc(node)

		actualChildNodes := make([]*tplinator.Node, 0)
		node.Children(func(_ int, childNode *tplinator.Node) bool {
			actualChildNodes = append(actualChildNodes, childNode)
			return true
		})

		if len(childNodes) != len(actualChildNodes) {
			t.Error("actual and expected child nodes does not have the same count")
		}
		for i := 0; i < len(childNodes); i++ {
			if childNodes[i] != actualChildNodes[i] {
				t.Error("the order of the actual and expected child nodes are not the same")
			}
		}
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			imgNode := tplinator.CreateNode(html.ElementNode, "img", nil, true)
			formNode := tplinator.CreateNode(html.ElementNode, "form", nil, false)
			h1Node := tplinator.CreateNode(html.ElementNode, "h1", nil, false)

			childNodes, modifyNodeFunc := tc.nodeFunc(imgNode, formNode, h1Node)
			testRemove(t, modifyNodeFunc, childNodes)
		})
	}
}
