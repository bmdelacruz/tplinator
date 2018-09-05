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
		testerFunc        func(t *testing.T, rootNodes []*tplinator.Node, err error)
	}{
		{
			name:        "simple h1",
			inputString: `<h1></h1>`,

			parserOptionFuncs: nil,
			testerFunc: func(t *testing.T, rootNodes []*tplinator.Node, err error) {
				if err != nil {
					t.Errorf("failed to parse input. cause: %v", err)
				} else if len(rootNodes) != 1 {
					t.Errorf("failed to parse input correctly. rootNodes: %+v", rootNodes)
				} else if h1Node := rootNodes[0]; h1Node.Data != "h1" {
					t.Errorf("failed to parse input correctly. h1Node: %+v", *h1Node)
				}
			},
		},
		{
			name:        "doctype",
			inputString: `<!doctype html>`,

			parserOptionFuncs: nil,
			testerFunc: func(t *testing.T, rootNodes []*tplinator.Node, err error) {
				if err != nil {
					t.Errorf("failed to parse input. cause: %v", err)
				} else if len(rootNodes) != 1 {
					t.Errorf("failed to parse input correctly. rootNodes: %+v", rootNodes)
				} else if doctypeNode := rootNodes[0]; doctypeNode.Type != html.DoctypeNode || doctypeNode.Data != "html" {
					t.Errorf("failed to parse input correctly. doctypeNode: %+v", *doctypeNode)
				}
			},
		},
		{
			name:        "self-closing element (img)",
			inputString: `<img />`,

			parserOptionFuncs: nil,
			testerFunc: func(t *testing.T, rootNodes []*tplinator.Node, err error) {
				if err != nil {
					t.Errorf("failed to parse input. cause: %v", err)
				} else if len(rootNodes) != 1 {
					t.Errorf("failed to parse input correctly. rootNodes: %+v", rootNodes)
				} else if imgNode := rootNodes[0]; imgNode.Type != html.ElementNode || imgNode.Data != "img" {
					t.Errorf("failed to parse input correctly. imgNode: %+v", *imgNode)
				}
			},
		},
		{
			name:        "comment",
			inputString: `<!--some comment-->`,

			parserOptionFuncs: nil,
			testerFunc: func(t *testing.T, rootNodes []*tplinator.Node, err error) {
				if err != nil {
					t.Errorf("failed to parse input. cause: %v", err)
				} else if len(rootNodes) != 0 {
					t.Errorf("failed to parse input correctly. rootNodes: %+v", rootNodes)
				}
			},
		},
		{
			name:        "node processor",
			inputString: `<h1 go-if="isMorning">hello</h1>`,

			parserOptionFuncs: []tplinator.ParserOptionFunc{
				tplinator.NodeProcessorsParserOption(
					tplinator.ConditionalExtensionNodeProcessor,
				),
			},
			testerFunc: func(t *testing.T, rootNodes []*tplinator.Node, err error) {
				if err != nil {
					t.Errorf("failed to parse input. cause: %v", err)
				} else if len(rootNodes) != 1 {
					t.Errorf("failed to parse input correctly. rootNodes: %+v", rootNodes)
				}
			},
		},
		{
			name:        "reach eof",
			inputString: `<h1>hell`,

			parserOptionFuncs: nil,
			testerFunc: func(t *testing.T, rootNodes []*tplinator.Node, err error) {
				if err == nil {
					t.Error("expecting an error because input string unexpectedly ended")
				}
			},
		},
		{
			name:        "incorrect doctype",
			inputString: `<img/><!doctype html>`,

			parserOptionFuncs: nil,
			testerFunc: func(t *testing.T, rootNodes []*tplinator.Node, err error) {
				if err == nil {
					t.Error("expecting an error because doctype is not the first element")
				}
			},
		},
		{
			name:        "tag mismatch",
			inputString: `<h1></h2>`,

			parserOptionFuncs: nil,
			testerFunc: func(t *testing.T, rootNodes []*tplinator.Node, err error) {
				if err == nil {
					t.Error("expecting an error because the start tag and end tag does not match")
				}
			},
		},
		{
			name:        "missing start tag",
			inputString: `</h2>`,

			parserOptionFuncs: nil,
			testerFunc: func(t *testing.T, rootNodes []*tplinator.Node, err error) {
				if err == nil {
					t.Error("expecting an error because an end tag was found but there was no pending start tag")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rootNodes, err := tplinator.ParseNodes(
				strings.NewReader(tc.inputString),
				tc.parserOptionFuncs...,
			)
			tc.testerFunc(t, rootNodes, err)
		})
	}
}
