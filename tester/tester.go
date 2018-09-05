package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alediaferia/stackgo"
	"github.com/bmdelacruz/tplinator"
	"github.com/yosssi/gohtml"
)

func main() {
	templateFile, err := os.Open("tester.html")
	if err != nil {
		log.Fatalln(err)
	}
	rootNodes, err := tplinator.ParseTemplate(templateFile)
	if err != nil {
		log.Fatalln(err)
	}

	for _, rootNode := range rootNodes {
		printNode(rootNode)

	}
}

type startTag struct {
	node *tplinator.Node
	tag  string
}

func (t startTag) getTag() string {
	return t.tag
}

type endTag struct {
	tag string
}

func (t endTag) getTag() string {
	return t.tag
}

func printNode(node *tplinator.Node) {
	var sb strings.Builder

	nodeStack := stackgo.NewStack()
	pushNode := func(node *tplinator.Node) {
		st, et := node.Tags()
		if et != "" {
			nodeStack.Push(endTag{
				tag: et,
			})
		}
		nodeStack.Push(startTag{
			node: node,
			tag:  st,
		})
	}

	pushNode(node)

	for nodeStack.Size() > 0 {
		switch tag := nodeStack.Pop().(type) {
		case startTag:
			sb.WriteString(tag.tag)

			var children []*tplinator.Node
			tag.node.Children(func(_ int, child *tplinator.Node) bool {
				children = append(children, child)
				return true
			})
			for i := len(children) - 1; i >= 0; i-- {
				pushNode(children[i])
			}
		case endTag:
			sb.WriteString(tag.tag)
		default:
			fmt.Printf("assertion error. %+v\n", tag)
		}
	}

	gohtml.Condense = true
	fmt.Println(gohtml.Format(sb.String()))
}
