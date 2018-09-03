package tplinator

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type attribute struct {
	key   string
	value string
}

type node struct {
	nodeType html.NodeType
	data     string

	attributes []attribute
	// attributeMap map[string]*attribute

	extensions []nodeExtension

	parent      *node
	firstChild  *node
	lastChild   *node
	prevSibling *node
	nextSibling *node
}

func createNode(srcNode *html.Node) *node {
	attrs := make([]attribute, len(srcNode.Attr))
	for attrIdx, attr := range srcNode.Attr {
		attrs[attrIdx] = attribute{
			key:   attr.Key,
			value: attr.Val,
		}
	}
	return &node{
		nodeType: srcNode.Type,
		data:     srcNode.Data,

		attributes: attrs,
	}
}

func (n *node) Execute(interpolator interpolator, evaluator evaluator) (string, error) {
	switch n.nodeType {
	case html.TextNode:
		data, err := interpolator(n.data)
		if err != nil {
			return "", err
		}
		return data, nil
	case html.DoctypeNode:
		return "<!DOCTYPE html>", nil
	case html.ElementNode:
		var err error

		// apply the node's extensions to create the node that's
		// going to be used
		fnode := n // final node
		for _, ext := range n.extensions {
			fnode, err = ext.Apply(*fnode, interpolator, evaluator)
			if err != nil {
				return "", err
			} else if fnode == nil {
				return "", nil
			}
		}

		// evaluate string interpolatons in attributes
		attributes := make([]attribute, len(fnode.attributes))
		copy(attributes, fnode.attributes)

		for ak, av := range attributes {
			evaluatedAv, err := interpolator(av.value)
			if err != nil {
				return "", err
			}
			attributes[ak].value = evaluatedAv
		}

		// create string representation of the html element
		var sb strings.Builder

		sb.WriteString(fmt.Sprintf("<%s", fnode.data))
		for _, attr := range attributes {
			if attr.value == "" {
				sb.WriteString(fmt.Sprintf(" %s", attr.key))
			} else {
				sb.WriteString(fmt.Sprintf(" %s=\"%s\"", attr.key, attr.value))
			}
		}

		if isSelfClosing(fnode.data) {
			if fnode.firstChild != nil {
				panic("self closing element has children")
			}
			sb.WriteString("/>")
		} else {
			sb.WriteString(">")
			for cn := fnode.firstChild; cn != nil; cn = cn.nextSibling {
				cnstr, err := cn.Execute(interpolator, evaluator)
				if err != nil {
					return "", err
				}
				sb.WriteString(cnstr)
			}
			sb.WriteString(fmt.Sprint("</", fnode.data, ">"))
		}

		return sb.String(), nil
	default:
		var sb strings.Builder

		for cn := n.firstChild; cn != nil; cn = cn.nextSibling {
			cnstr, err := cn.Execute(interpolator, evaluator)
			if err != nil {
				return "", err
			}
			sb.WriteString(cnstr)
		}

		return sb.String(), nil
	}
}

func (n *node) Attributes() []attribute {
	attributesCopy := make([]attribute, len(n.attributes))
	copy(attributesCopy, n.attributes)

	return attributesCopy
}

func (n *node) HasAttribute(key string) (bool, int, string) {
	for attrIdx, attr := range n.attributes {
		if attr.key == key {
			return true, attrIdx, attr.value
		}
	}
	return false, -1, ""
}

func (n *node) AddAttribute(key, value string) {
	n.attributes = append(n.attributes, attribute{key: key, value: value})
}

func (n *node) ReplaceAttribute(key string, value string) {
	targetIdx := -1
	for attrIdx, attr := range n.attributes {
		if attr.key == key {
			targetIdx = attrIdx
		}
	}
	if targetIdx >= 0 {
		n.attributes[targetIdx] = attribute{key: key, value: value}
	}
}

func (n *node) RemoveAttribute(key string) {
	targetIdx := -1
	for attrIdx, attr := range n.attributes {
		if attr.key == key {
			targetIdx = attrIdx
		}
	}
	if targetIdx >= 0 {
		n.attributes = append(n.attributes[:targetIdx], n.attributes[targetIdx+1:]...)
	}
}

func (n *node) Insert(newChildNode *node, beforeChildNode *node) {
	if newChildNode.parent != nil || newChildNode.prevSibling != nil || newChildNode.nextSibling != nil {
		panic("the node is already a child of another node")
	}

	var prev, next *node
	if beforeChildNode != nil {
		prev, next = beforeChildNode.prevSibling, beforeChildNode
	} else {
		prev = n.lastChild
	}
	if prev != nil {
		prev.nextSibling = newChildNode
	} else {
		n.firstChild = newChildNode
	}
	if next != nil {
		next.prevSibling = newChildNode
	} else {
		n.lastChild = newChildNode
	}

	newChildNode.parent = n
	newChildNode.prevSibling = prev
	newChildNode.nextSibling = next
}

func (n *node) AppendChild(node *node) {
	if node.parent != nil || node.prevSibling != nil || node.nextSibling != nil {
		panic("the node is already a child of another node")
	}

	last := n.lastChild
	if last != nil {
		last.nextSibling = node
	} else {
		n.firstChild = node
	}

	n.lastChild = node
	node.parent = n
	node.prevSibling = last
}

func isSelfClosing(tag string) bool {
	return equalToOneOf(
		"area", "base", "br", "col", "command", "embed",
		"hr", "img", "input", "keygen", "link", "menuitem",
		"meta", "param", "source", "track", "wbr",
	)(tag)
}
