package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	docBytes, err := ioutil.ReadFile("tester.html")
	if err != nil {
		panic(err)
	}
	rootNode, err := html.Parse(bytes.NewReader(docBytes))
	if err != nil {
		panic(err)
	}

	cleanTextNodes(rootNode)
	// printNode(rootNode)

	var buffer bytes.Buffer
	html.Render(&buffer, rootNode)
	fmt.Printf("%s\n", buffer.Bytes())
}

func cleanTextNodes(node *html.Node) {
	var toBeRemoved []*html.Node

	for cn := node.FirstChild; cn != nil; cn = cn.NextSibling {
		if cn.Type == html.TextNode {
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

func printNode(node *html.Node) {
	printNodeLevel(node, 0)
}

func printNodeLevel(node *html.Node, level int) {
	switch node.Type {
	case html.ElementNode:
		fmt.Print(strings.Repeat(" ", level)+"|<", node.Data, ">{")
	case html.TextNode:
		fmt.Print(strings.Repeat(" ", level)+"|\"", strings.TrimSpace(node.Data), "\"{")
	default:
		fmt.Print(strings.Repeat(" ", level)+"|(type:", node.Type, "){")
	}
	for _, attr := range node.Attr {
		fmt.Print(attr.Key, ":\"", attr.Val, "\",")
	}
	fmt.Print("}[")
	if hasIf, _ := hasAttribute(node, "go-if"); hasIf {
		fmt.Print("if,")
	}
	fmt.Print("]")
	fmt.Print("\n")

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		printNodeLevel(c, level+1)
	}
}

func hasAttribute(node *html.Node, attributeKey string) (bool, int) {
	for attrIdx, attr := range node.Attr {
		if attr.Key == attributeKey {
			return true, attrIdx
		}
	}
	return false, -1
}

type modifiedHTMLNode struct {
	html.Node
}

type conditionalClassAttribute struct {
	SourceNode *html.Node

	Index int
	Class string
}

func hasAttributeWithPrefix(node *html.Node, attributeKeyPrefix string) []conditionalClassAttribute {
	var attrs []conditionalClassAttribute
	for attrIdx, attr := range node.Attr {
		if strings.HasPrefix(attr.Key, attributeKeyPrefix) {
			attrs = append(attrs, conditionalClassAttribute{
				Index: attrIdx,
				Class: strings.TrimPrefix(attr.Key, attributeKeyPrefix),
			})
		}
	}
	return attrs
}
