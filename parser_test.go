package tplinator_test

import (
	"strings"
	"testing"

	"github.com/bmdelacruz/tplinator"
	"golang.org/x/net/html"
)

func TestParseTemplate(t *testing.T) {
	testCases := []struct {
		name        string
		inputString string

		parserOptionFuncs []tplinator.ParserOptionFunc
		testerFunc        func(t *testing.T, rootNodes []*tplinator.Node)
	}{
		{
			name:        "simple H1",
			inputString: `<h1>Hello, world!</h1>`,

			parserOptionFuncs: nil,
			testerFunc: func(t *testing.T, rootNodes []*tplinator.Node) {
				if len(rootNodes) != 1 {
					t.Errorf("failed to parse input correctly. rootNodes: %+v", rootNodes)
					return
				}

				h1Node := rootNodes[0]

				h1NodeChildren := make([]*tplinator.Node, 0)
				h1Node.Children(func(_ int, child *tplinator.Node) bool {
					h1NodeChildren = append(h1NodeChildren, child)
					return true
				})

				if h1Node.Data != "h1" || len(h1NodeChildren) != 1 {
					t.Errorf("failed to parse input correctly. h1Node: %+v", *h1Node)
					return
				}

				h1TextNode := h1NodeChildren[0]
				expectedText := "Hello, world!"

				if h1TextNode.Type != html.TextNode || h1TextNode.Data != expectedText {
					t.Errorf("failed to parse input correctly. got %v, wanted `%v`",
						h1TextNode.Data, expectedText)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rootNodes, err := tplinator.ParseTemplate(
				strings.NewReader(tc.inputString),
			)
			if err != nil {
				t.Errorf("failed to parse input. cause: %v", err)
			} else {
				tc.testerFunc(t, rootNodes)
			}
		})
	}
}
