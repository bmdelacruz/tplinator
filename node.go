package tplinator

import (
	"errors"

	"github.com/alediaferia/stackgo"
	"golang.org/x/net/html"
)

type Node struct {
	Data string
	Type html.NodeType

	isSelfClosing bool

	attributes []Attribute
	extensions []Extension

	contextParams EvaluatorParams
	parentECS     EvaluatorContextSource

	parent      *Node
	firstChild  *Node
	lastChild   *Node
	prevSibling *Node
	nextSibling *Node
}

func CreateNode(
	nodeType html.NodeType, data string,
	attributes []html.Attribute, isSelfClosing bool,
) *Node {
	attrs := make([]Attribute, len(attributes))
	for attrIdx, attr := range attributes {
		attrs[attrIdx] = Attribute{
			Key:     attr.Key,
			Value:   attr.Val,
			KeyOnly: attr.Val == "",
		}
	}
	return &Node{
		Data: data,
		Type: nodeType,

		isSelfClosing: isSelfClosing,
		attributes:    attrs,
	}
}

func CopyNode(node *Node) *Node {
	copyNode := func(n *Node) *Node {
		nodeCopy := &Node{
			Data:          n.Data,
			Type:          n.Type,
			isSelfClosing: n.isSelfClosing,
			contextParams: n.contextParams,
		}
		for _, attr := range n.attributes {
			nodeCopy.attributes = append(nodeCopy.attributes, attr)
		}
		for _, ext := range n.extensions {
			nodeCopy.extensions = append(nodeCopy.extensions, ext)
		}
		return nodeCopy
	}

	type root struct {
		nodeCopy *Node
		origNode *Node
	}
	type nonRoot struct {
		nodeCopy  *Node
		origNode  *Node
		newParent *Node
	}

	nodeCopy := copyNode(node)

	copierStack := stackgo.NewStack()
	copierStack.Push(root{
		origNode: node,
		nodeCopy: nodeCopy,
	})

	for copierStack.Top() != nil {
		switch stackItem := copierStack.Pop().(type) {
		case root:
			origNode := stackItem.origNode
			newParent := stackItem.nodeCopy

			var children []*Node
			origNode.Children(func(_ int, child *Node) bool {
				children = append(children, child)
				return true
			})
			for i := len(children) - 1; i >= 0; i-- {
				copierStack.Push(nonRoot{
					newParent: newParent,
					origNode:  children[i],
					nodeCopy:  copyNode(children[i]),
				})
			}
		case nonRoot:
			origNode := stackItem.origNode
			newParent := stackItem.newParent
			newParent.AppendChild(stackItem.nodeCopy)

			var children []*Node
			origNode.Children(func(_ int, child *Node) bool {
				children = append(children, child)
				return true
			})
			for i := len(children) - 1; i >= 0; i-- {
				copierStack.Push(nonRoot{
					newParent: stackItem.nodeCopy,
					origNode:  children[i],
					nodeCopy:  copyNode(children[i]),
				})
			}
		}
	}

	return nodeCopy
}

func (n Node) ContextParams() EvaluatorParams {
	return n.contextParams
}

func (n *Node) SetContextParams(contextParams EvaluatorParams) {
	n.contextParams = contextParams
}

func (n *Node) GetContextParams() []EvaluatorParams {
	var evaluatorParams []EvaluatorParams

	if n.parentECS != nil {
		if n.contextParams != nil {
			evaluatorParams = append(evaluatorParams, n.contextParams)
		}
		if parentCP := n.parentECS.GetContextParams(); len(parentCP) > 0 {
			evaluatorParams = append(evaluatorParams, parentCP...)
		}
		return evaluatorParams
	}

	epStack := stackgo.NewStack()
	epStack.Push(n)

	for epStack.Top() != nil {
		currentNode := epStack.Pop().(*Node)
		if currentNode.contextParams != nil {
			evaluatorParams = append(
				evaluatorParams, currentNode.contextParams,
			)
		}
		if parentNode := currentNode.Parent(); parentNode != nil {
			epStack.Push(parentNode)
		}
	}

	return evaluatorParams
}

func (n *Node) SetParentEvaluatorContextSource(ecs EvaluatorContextSource) {
	n.parentECS = ecs
}

func (n Node) Parent() *Node {
	return n.parent
}

func (n Node) FirstChild() *Node {
	return n.firstChild
}

func (n Node) LastChild() *Node {
	return n.lastChild
}

func (n Node) PreviousSibling() *Node {
	return n.prevSibling
}

func (n Node) NextSibling() *Node {
	return n.nextSibling
}

func (n Node) Tags() (string, string) {
	switch n.Type {
	case html.DoctypeNode:
		return "<!DOCTYPE " + n.Data + ">", ""
	case html.TextNode:
		return n.Data, ""
	case html.ElementNode:
		startTag := "<" + n.Data
		for _, attr := range n.attributes {
			startTag += " " + attr.String()
		}
		if n.isSelfClosing {
			return startTag + "/>", ""
		}
		return startTag + ">", "</" + n.Data + ">"
	case html.ErrorNode, html.CommentNode, html.DocumentNode:
		fallthrough
	default:
		panic(errors.New("assertion error"))
	}
}

func (n Node) Children(nodeFunc func(int, *Node) bool) {
	index := 0
	for child := n.firstChild; child != nil; child = child.nextSibling {
		if shouldContinue := nodeFunc(index, child); !shouldContinue {
			break
		}
		index++
	}
}

func (n Node) NextSiblings(nodeFunc func(*Node) bool) {
	for sibling := n.nextSibling; sibling != nil; sibling = sibling.nextSibling {
		if shouldContinue := nodeFunc(sibling); !shouldContinue {
			break
		}
	}
}

func (n *Node) AddExtension(extension Extension) {
	if extension == nil {
		panic(errors.New("assertion error"))
	}
	n.extensions = append(n.extensions, extension)
}

func (n *Node) ApplyExtensions(dependencies ExtensionDependencies, params EvaluatorParams) (*Node, []*Node, error) {
	siblings := make([]*Node, 0)
	currentNode := n
	for _, extension := range n.extensions {
		newNode, newSibs, err := extension.Apply(currentNode, dependencies, params)
		if err != nil {
			return nil, nil, err
		} else if newNode != nil {
			currentNode = newNode
		} else {
			currentNode = nil
			siblings = append(siblings, newSibs...)
			break
		}
	}
	return currentNode, siblings, nil
}

func (n Node) Attributes() []Attribute {
	attributesCopy := make([]Attribute, len(n.attributes))
	copy(attributesCopy, n.attributes)

	return attributesCopy
}

func (n Node) HasAttribute(key string) (bool, int, string) {
	for attrIdx, attr := range n.attributes {
		if attr.Key == key {
			return true, attrIdx, attr.Value
		}
	}
	return false, -1, ""
}

func (n Node) HasAttributes(testFunc func(Attribute) bool) []Attribute {
	var matches []Attribute
	for _, attr := range n.attributes {
		if testFunc(attr) {
			matches = append(matches, attr)
		}
	}
	return matches
}

func (n *Node) AddAttribute(key, value string) {
	if attrAlreadyExists, _, _ := n.HasAttribute(key); attrAlreadyExists {
		n.ReplaceAttribute(key, value)
	} else {
		n.attributes = append(n.attributes, Attribute{Key: key, Value: value})
	}
}

func (n *Node) ReplaceAttribute(key string, value string) {
	targetIdx := -1
	for attrIdx, attr := range n.attributes {
		if attr.Key == key {
			targetIdx = attrIdx
		}
	}
	if targetIdx >= 0 {
		n.attributes[targetIdx] = Attribute{Key: key, Value: value}
	}
}

func (n *Node) RemoveAttribute(key string) {
	targetIdx := -1
	for attrIdx, attr := range n.attributes {
		if attr.Key == key {
			targetIdx = attrIdx
		}
	}
	if targetIdx >= 0 {
		n.attributes = append(n.attributes[:targetIdx], n.attributes[targetIdx+1:]...)
	}
}

func (n *Node) Insert(newChildNode *Node, beforeChildNode *Node) {
	if newChildNode.parent != nil || newChildNode.prevSibling != nil || newChildNode.nextSibling != nil {
		panic("the node is already a child of another node")
	}

	var prev, next *Node
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

func (n *Node) AppendChild(node *Node) {
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

func (n *Node) RemoveChild(child *Node) {
	if child.parent != n {
		panic("this node is not the parent of the specified child")
	}
	if n.firstChild == child {
		n.firstChild = child.nextSibling
	}
	if child.nextSibling != nil {
		child.nextSibling.prevSibling = child.prevSibling
	}
	if n.lastChild == child {
		n.lastChild = child.prevSibling
	}
	if child.prevSibling != nil {
		child.prevSibling.nextSibling = child.nextSibling
	}
	child.parent = nil
	child.prevSibling = nil
	child.nextSibling = nil
}

type Attribute struct {
	Key     string
	Value   string
	KeyOnly bool
}

func (a Attribute) String() string {
	if a.KeyOnly {
		return a.Key
	}
	return a.Key + "=\"" + a.Value + "\""
}
