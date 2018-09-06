package tplinator

import (
	"bytes"
	"io"
	"strings"

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

func (tpl *Template) RenderBytes(params EvaluatorParams) ([]byte, error) {
	var bb bytes.Buffer
	if err := tpl.Render(params, func(str string) {
		bb.WriteString(str)
	}); err != nil {
		return nil, err
	}
	return bb.Bytes(), nil
}

func (tpl *Template) RenderString(params EvaluatorParams) (string, error) {
	var sb strings.Builder
	if err := tpl.Render(params, func(str string) {
		sb.WriteString(str)
	}); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func (tpl *Template) Render(params EvaluatorParams, writerFunc func(string)) error {
	tagStack := stackgo.NewStack()

	pushNode := func(node *Node) {
		if node == nil {
			return
		}
		st, et := node.Tags()
		if et != "" {
			tagStack.Push(tplEndTag{tag: et})
		}
		tagStack.Push(tplStartTag{node: node, tag: st})
	}
	applyNodeExts := func(node *Node) error {
		node, sibs, err := node.ApplyExtensions(&tpl.extDeps, params)
		if err != nil {
			return err
		}
		for i := len(sibs) - 1; i >= 0; i-- {
			pushNode(sibs[i])
		}
		if node != nil {
			pushNode(node)
		}
		return nil
	}

	for _, rootNode := range tpl.rootNodes {
		err := applyNodeExts(rootNode)
		if err != nil {
			return err
		}
		for tagStack.Top() != nil {
			switch tag := tagStack.Pop().(type) {
			case tplStartTag:
				writerFunc(tag.tag)

				var children []*Node
				tag.node.Children(func(_ int, child *Node) bool {
					children = append(children, child)
					return true
				})
				for i := len(children) - 1; i >= 0; i-- {
					err := applyNodeExts(children[i])
					if err != nil {
						return err
					}
				}
			case tplEndTag:
				writerFunc(tag.tag)
			}
		}
	}

	return nil
}

type tplStartTag struct {
	node *Node
	tag  string
}

type tplEndTag struct {
	tag string
}
