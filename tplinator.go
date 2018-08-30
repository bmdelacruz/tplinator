package tplinator

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/Knetic/govaluate"
	"golang.org/x/net/html"
)

// Template related constants
const (
	attributePrefix    = "go-"
	ifAttribute        = attributePrefix + "if"
	elseIfAttribute    = attributePrefix + "else-if"
	elseAttribute      = attributePrefix + "else"
	rangeAttribute     = attributePrefix + "range"
	interpolationStart = "{{go:"
	interpolationEnd   = "}}"
)

var interpolationRegexPattern = regexp.MustCompile(interpolationStart + "[\\d\\w]+" + interpolationEnd)

type Template struct {
	rootNode *html.Node
}

func CreateTemplate(reader io.Reader) (*Template, error) {
	rootNode, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}
	return &Template{
		rootNode: rootNode,
	}, nil
}

func (t *Template) Execute(data map[string]interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	if _, err := execute(t.rootNode, data); err != nil {
		return nil, err
	}
	if err := html.Render(&buffer, t.rootNode); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func execute(node *html.Node, data map[string]interface{}) (bool, error) {
	if node.Type == html.ElementNode {
		if didGetIfAttr, attrIndex, attrExpr := tryGetIfAttribute(node); didGetIfAttr {
			// remove the go-if attribute from the node
			node.Attr = append(node.Attr[:attrIndex], node.Attr[attrIndex+1:]...)

			boolResult, err := applyIfAttribute(node, attrExpr, data)
			if err != nil {
				return false, err
			} else if !boolResult {
				// TODO: check if next sibling node has go-else-if/go-else attributes

				// stop processing this node and its children since it is
				// going to be removed anyway
				return true, nil
			}
		}

		// TODO: process go-range elements

		// evaluate string interpolation on attributes
		attrsToReplace := make(map[int]html.Attribute)
		for attrIdx, attr := range node.Attr {
			// try to check if there are interpolation variables within the
			// attribute value. if there is, replace the attribute value with
			// the evaluated version of it.
			hasInterpolated, interpolated, err := tryToInterpolate(attr.Val, data)
			if err != nil {
				return false, err
			} else if hasInterpolated {
				attrsToReplace[attrIdx] = html.Attribute{
					Key:       attr.Key,
					Namespace: attr.Namespace,
					Val:       interpolated,
				}
			}
		}
		// idk if it's a bug, but `html.Attribute` struct seems to be immutable
		// so i recreate them instead
		for k, v := range attrsToReplace {
			node.Attr[k] = v
		}
	} else if node.Type == html.TextNode {
		// evaluate string interpolation on `html.TextNode`s
		hasInterpolated, interpolated, err := tryToInterpolate(node.Data, data)
		if err != nil {
			return false, err
		} else if hasInterpolated {
			node.Data = interpolated
		}
	}

	var childrenToBeRemoved []*html.Node

	for n := node.FirstChild; n != nil; n = n.NextSibling {
		shouldRemoveNode, err := execute(n, data)
		if err != nil {
			return false, err
		} else if shouldRemoveNode {
			childrenToBeRemoved = append(childrenToBeRemoved, n)
		}
	}
	for _, n := range childrenToBeRemoved {
		node.RemoveChild(n)
	}

	return false, nil
}

func tryGetIfAttribute(node *html.Node) (bool, int, string) {
	if node.Type == html.ElementNode {
		for attrIndex, attr := range node.Attr {
			if attr.Key == ifAttribute {
				return true, attrIndex, attr.Val
			}
		}
	}
	return false, 0, ""
}

func applyIfAttribute(node *html.Node, ifExprStr string, data map[string]interface{}) (bool, error) {
	expr, err := govaluate.NewEvaluableExpression(ifExprStr)
	if err != nil {
		return false, err
	}

	result, err := expr.Evaluate(data)
	if err != nil {
		return false, err
	}
	boolResult, isBool := result.(bool)
	if !isBool {
		return false, fmt.Errorf(`"%v" is not a boolean expression`, ifExprStr)
	}

	return boolResult, nil
}

func tryToInterpolate(inputStr string, data map[string]interface{}) (bool, string, error) {
	outputString := inputStr
	matches := interpolationRegexPattern.FindAllString(inputStr, -1)

	for _, match := range matches {
		origMatch := match

		match = strings.TrimLeft(match, interpolationStart)
		match = strings.TrimRight(match, interpolationEnd)
		exprStr := strings.TrimSpace(match)

		expr, err := govaluate.NewEvaluableExpression(exprStr)
		if err != nil {
			return false, "", err
		}
		result, err := expr.Evaluate(data)
		if err != nil {
			return false, "", err
		}

		outputString = strings.Replace(outputString, origMatch, fmt.Sprintf("%v", result), 1)
	}

	return len(matches) > 0, outputString, nil
}
