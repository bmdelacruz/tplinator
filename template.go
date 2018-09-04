package tplinator

import (
	"bytes"
	"io"

	"github.com/alediaferia/stackgo"
)

type Template struct {
	rootNodes []*Node
	extDeps   compoundExtensionDependencies
}

func CreateTemplateFromReader(reader io.Reader, parserOptions ...ParserOptionFunc) (*Template, error) {
	rootNodes, err := ParseNodes(reader, parserOptions...)
	if err != nil {
		return nil, err
	}
	return &Template{
		rootNodes: rootNodes,
	}, nil
}

func (tpl *Template) AddExtensionDependencies(extDeps ...ExtensionDependencies) {
	tpl.extDeps.extDeps = append(tpl.extDeps.extDeps, extDeps...)
}

func (tpl *Template) Render(params EvaluatorParams) ([]byte, error) {
	var bb bytes.Buffer

	tagStack := stackgo.NewStack()
	pushNode := func(node *Node) error {
		node, err := node.ApplyExtensions(&tpl.extDeps, params)
		if err != nil {
			return err
		} else if node == nil {
			return nil
		}

		st, et := node.Tags()
		if et != "" {
			tagStack.Push(tplEndTag{tag: et})
		}
		tagStack.Push(tplStartTag{node: node, tag: st})

		return nil
	}

	for _, rootNode := range tpl.rootNodes {
		err := pushNode(rootNode)
		if err != nil {
			return nil, err
		}
		for tagStack.Top() != nil {
			switch tag := tagStack.Pop().(type) {
			case tplStartTag:
				bb.WriteString(tag.tag)

				var children []*Node
				tag.node.Children(func(_ int, child *Node) bool {
					children = append(children, child)
					return true
				})
				for i := len(children) - 1; i >= 0; i-- {
					err := pushNode(children[i])
					if err != nil {
						return nil, err
					}
				}
			case tplEndTag:
				bb.WriteString(tag.tag)
			}
		}
	}

	return bb.Bytes(), nil
}

type tplStartTag struct {
	node *Node
	tag  string
}

type tplEndTag struct {
	tag string
}
