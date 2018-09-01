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
	return &PrecompiledTemplate{
		documentNode: precompileToNode(documentNode),
	}, nil
}

func precompileToNode(documentNode *html.Node) *node {
	cleanTextNodes(documentNode)
	return convertToNode(documentNode)
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

func convertToNode(srcNode *html.Node) *node {
	extensions := createExtensions(srcNode)

	noode := createNode(srcNode)
	noode.extensions = extensions

	for cn := srcNode.FirstChild; cn != nil; cn = cn.NextSibling {
		noode.AppendChild(convertToNode(cn))
	}

	return noode
}

func createExtensions(srcNode *html.Node) []nodeExtension {
	// TODO
	// - check if a child of the srcNode has conditional attribute
	// - etc

	return nil
}
