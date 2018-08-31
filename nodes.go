package tplinator

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type node struct {
	nodeType html.NodeType
	data     string

	attributes map[string]string
	extensions []nodeExtension

	parent      *node
	firstChild  *node
	lastChild   *node
	prevSibling *node
	nextSibling *node
}

func createNode(srcNode *html.Node) *node {
	return &node{} // TODO
}

func (n *node) Execute(evaluator evaluator) (string, error) {
	switch n.nodeType {
	case html.TextNode:
		data, err := evaluator(n.data)
		if err != nil {
			return "", err
		}
		return data, nil
	case html.DoctypeNode:
		return "<!DOCTYPE html>", nil
	case html.ElementNode:
		// apply the node's extensions to create the node that's
		// going to be used
		fnode := n // final node
		for _, ext := range n.extensions {
			fnode, err := ext.Apply(*fnode, evaluator)
			if err != nil {
				return "", err
			} else if fnode == nil {
				return "", nil
			}
		}

		// evaluate string interpolatons in attributes
		attributes := copyAttributes(fnode.attributes)
		for ak, av := range attributes {
			evaluatedAv, err := evaluator(av)
			if err != nil {
				return "", err
			}
			attributes[ak] = evaluatedAv
		}

		// create string representation of the html element
		var sb strings.Builder

		sb.WriteString(fmt.Sprint("<", fnode.data))
		for attrKey, attrVal := range attributes {
			sb.WriteString(fmt.Sprintf(" %s=\"%s\"", attrKey, attrVal))
		}
		sb.WriteString(">")
		for cn := fnode.firstChild; cn != nil; cn = cn.nextSibling {
			cnstr, err := cn.Execute(evaluator)
			if err != nil {
				return "", err
			}
			sb.WriteString(cnstr)
		}
		sb.WriteString(fmt.Sprint("</", fnode.data, ">"))

		return sb.String(), nil
	default:
		return "", nil
	}
}

func (n *node) Attributes() map[string]string {
	return copyAttributes(n.attributes)
}

func (n *node) AddAttribute(key, value string) {
	n.attributes[key] = value
}

func (n *node) RemoveAttribute(key string) {
	delete(n.attributes, key)
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

func copyAttributes(attrs map[string]string) map[string]string {
	attrsCopy := make(map[string]string)
	for k, v := range attrs {
		attrsCopy[k] = v
	}
	return attrsCopy
}
