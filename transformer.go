package tplinator

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

func Tplinate(documentReader io.Reader) (*PrecompiledTemplate, error) {
	documentNode, err := html.Parse(documentReader)
	if err != nil {
		return nil, err
	}
	ptNode, err := precompileToNode(documentNode)
	if err != nil {
		return nil, err
	}
	return &PrecompiledTemplate{
		documentNode: ptNode,
	}, nil
}

func precompileToNode(documentNode *html.Node) (*node, error) {
	cleanTextNodes(documentNode)

	return nil, nil
}

func cleanTextNodes(node *html.Node) {
	var toBeRemoved []*html.Node

	// check if node has TextNode children.
	for cn := node.FirstChild; cn != nil; cn = cn.NextSibling {
		if cn.Type == html.TextNode {
			// if the string on the node becomes an empty string,
			// remove it from the tree
			cn.Data = strings.TrimSpace(cn.Data)
			if len(cn.Data) == 0 {
				toBeRemoved = append(toBeRemoved, cn)
			}
		}
	}
	for _, cn := range toBeRemoved {
		cn.Parent.RemoveChild(cn)
	}

	for cn := node.FirstChild; cn != nil; cn = cn.NextSibling {
		cleanTextNodes(cn)
	}
}
