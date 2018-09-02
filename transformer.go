package tplinator

import (
	"errors"
	"fmt"
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
		} else if cn.Type == html.CommentNode {
			toBeRemoved = append(toBeRemoved, cn)
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
	fmt.Println("el proc", srcNode.Data)

	extensions := createExtensions(srcNode)
	noode := createNode(srcNode)
	noode.extensions = extensions

	for cn := srcNode.FirstChild; cn != nil; cn = cn.NextSibling {
		if childNode := generateEquivalentNode(cn); childNode != nil {
			noode.AppendChild(childNode)
		}
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
	if hasIfCond, ifCond := tryExtractAttribute(srcNode, ifAttr()); hasIfCond {
		cne := &conditionalNodeExtension{}
		cne.conditions = append(cne.conditions, &cneCondition{
			condition: ifCond,
			isSelf:    true,
		})

		// find siblings that needs to be removed from parent
		var elseSibling *html.Node
		var elifSiblings []*html.Node
		for sib := srcNode.NextSibling; sib != nil; sib = sib.NextSibling {
			if hasElseIfCond, elseIfCond := hasAttributeKey(sib, elseIfAttr()); hasElseIfCond {
				fmt.Println("elif sib proc", sib.Data)

				cne.conditions = append(cne.conditions, &cneCondition{
					condition: elseIfCond,
					node:      generateEquivalentNode(sib),
				})

				elifSiblings = append(elifSiblings, sib)
				continue
			} else if hasElse, _ := hasAttributeKey(sib, elseAttr()); hasElse {
				if elseSibling != nil {
					panic(errors.New("found an extraneous else element. please remove it"))
				}

				fmt.Println("else sib proc", sib.Data)
				cne.elseNode = generateEquivalentNode(sib)

				elseSibling = sib
				continue
			}
			// nextNonCondSibling = sib
			break
		}

		// remove siblings that are going to be part of the conditional
		// extension since they will no longer needed to be attached to
		// the document node
		for _, elifSibling := range elifSiblings {
			fmt.Println("rm elif sib", elifSibling.Data)
			elifSibling.Parent.RemoveChild(elifSibling)
		}
		if elseSibling != nil {
			fmt.Println("rm else sib", elseSibling.Data)
			elseSibling.Parent.RemoveChild(elseSibling)
		}

		return cne
	}

	return nil
}

func ifAttr() func(string) bool {
	return equalToOneOf("go-if")
}

func elseIfAttr() func(string) bool {
	return equalToOneOf("go-elif", "go-else-if")
}

func elseAttr() func(string) bool {
	return equalToOneOf("go-else")
}

func hasAttributeKey(node *html.Node, checker func(string) bool) (bool, string) {
	for _, attr := range node.Attr {
		if checker(attr.Key) {
			return true, attr.Val
		}
	}
	return false, ""
}

func tryExtractAttribute(node *html.Node, checker func(string) bool) (bool, string) {
	for attrIdx, attr := range node.Attr {
		if checker(attr.Key) {
			node.Attr = append(node.Attr[:attrIdx], node.Attr[attrIdx+1:]...)
			return true, attr.Val
		}
	}
	return false, ""
}

func equalToOneOf(strs ...string) func(string) bool {
	return func(input string) bool {
		for _, str := range strs {
			if input == str {
				return true
			}
		}
		return false
	}
}
