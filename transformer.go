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
	return generateEquivalentNode(documentNode)
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

func generateEquivalentNode(srcNode *html.Node) *node {
	extensions := createExtensions(srcNode)

	noode := createNode(srcNode)
	noode.extensions = extensions

	for cn := srcNode.FirstChild; cn != nil; cn = cn.NextSibling {
		noode.AppendChild(generateEquivalentNode(cn)) // FIXME
	}

	return noode
}

func createExtensions(srcNode *html.Node) []nodeExtension {
	var exts []nodeExtension

	// check if a child of the srcNode has conditional attribute
	if ext := tryCreateConditionalNodeExtension(srcNode); ext != nil {
		exts = append(exts, ext)
	}

	return exts
}

func tryCreateConditionalNodeExtension(srcNode *html.Node) nodeExtension {
	if ifCond, hasIfCond := tryExtractIfCondition(srcNode); hasIfCond {
		cne := &conditionalNodeExtension{}
		cne.conditions = append(cne.conditions, &cneCondition{
			condition: ifCond,
			node:      generateEquivalentNode(srcNode),
		})

		var toBeRemoved []*html.Node
		for sib := srcNode.NextSibling; sib != nil; sib = sib.NextSibling {
			if elseIfCond, hasElseIfCond := tryExtractElseIfCondition(sib); hasElseIfCond {
				toBeRemoved = append(toBeRemoved, sib)
				cne.conditions = append(cne.conditions, &cneCondition{
					condition: elseIfCond,
					node:      generateEquivalentNode(sib),
				})
				continue
			} else if hasElse := tryExtractElse(sib); hasElse {
				toBeRemoved = append(toBeRemoved, sib)
				cne.elseNode = generateEquivalentNode(sib)
			}
			break
		}

		srcNode.Parent.RemoveChild(srcNode)
		for _, tbr := range toBeRemoved {
			tbr.Parent.RemoveChild(tbr)
		}

		return cne
	}
	return nil
}

func tryExtractIfCondition(node *html.Node) (string, bool) {
	for attrIdx, attr := range node.Attr {
		if attr.Key == "go-if" {
			node.Attr = append(node.Attr[:attrIdx], node.Attr[attrIdx+1:]...)
			return attr.Val, true
		}
	}
	return "", false
}

func tryExtractElseIfCondition(node *html.Node) (string, bool) {
	for attrIdx, attr := range node.Attr {
		if attr.Key == "go-elif" || attr.Key == "go-else-if" {
			node.Attr = append(node.Attr[:attrIdx], node.Attr[attrIdx+1:]...)
			return attr.Val, true
		}
	}
	return "", false
}

func tryExtractElse(node *html.Node) bool {
	for attrIdx, attr := range node.Attr {
		if attr.Key == "go-else" {
			node.Attr = append(node.Attr[:attrIdx], node.Attr[attrIdx+1:]...)
			return true
		}
	}
	return false
}
