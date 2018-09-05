package tplinator

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/alediaferia/stackgo"
)

type Evaluator interface {
	EvaluateBool(input string) (bool, error)
	EvaluateString(input string) (string, error)
	Evaluate(input string) (interface{}, error)
}

type Template struct {
	rootNodes []*Node

	extDeps ExtensionDependencies // TODO
}

func CreateTemplateFromString(tplStr string, parserOptions ...ParserOptionFunc) (*Template, error) {
	return CreateTemplateFromReader(strings.NewReader(tplStr), parserOptions...)
}

func CreateTemplateFromBytes(tplBytes []byte, parserOptions ...ParserOptionFunc) (*Template, error) {
	return CreateTemplateFromReader(bytes.NewReader(tplBytes), parserOptions...)
}

func CreateTemplateFromFile(tplFilePath string, parserOptions ...ParserOptionFunc) (*Template, error) {
	tplFile, err := os.Open(tplFilePath)
	if err != nil {
		return nil, err
	}
	return CreateTemplateFromReader(tplFile, parserOptions...)
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

func (tpl *Template) Execute(params map[string]interface{}) ([]byte, error) {
	var bb bytes.Buffer

	tagStack := stackgo.NewStack()
	pushNode := func(node *Node, deps ExtensionDependencies) error {
		node, err := node.ApplyExtensions(deps)
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

	for i := len(tpl.rootNodes) - 1; i >= 0; i-- {
		err := pushNode(tpl.rootNodes[i], tpl.extDeps)
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
					err := pushNode(children[i], tpl.extDeps)
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
