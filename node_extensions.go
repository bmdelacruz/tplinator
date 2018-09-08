package tplinator

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type Extension interface {
	Apply(node *Node, dependencies ExtensionDependencies, params EvaluatorParams) (*Node, []*Node, error)
}

type conditionalExtensionCondition struct {
	node                  *Node
	conditionalExpression string
}

type ConditionalExtension struct {
	conditions []conditionalExtensionCondition
	elseNode   *Node
}

func (ce *ConditionalExtension) Apply(node *Node, dependencies ExtensionDependencies, params EvaluatorParams) (*Node, []*Node, error) {
	for _, condition := range ce.conditions {
		evaluator := dependencies.Get(EvaluatorExtDepKey).(Evaluator)
		hasResult, result, err := TryEvaluateBoolUsingContext(
			node, evaluator, condition.conditionalExpression,
		)
		if !hasResult {
			result, err = evaluator.EvaluateBool(condition.conditionalExpression, params)
		}
		if err != nil {
			return nil, nil, err
		} else if result {
			return condition.node, nil, nil
		}
	}
	return ce.elseNode, nil, nil
}

func (ce *ConditionalExtension) addCondition(condition string, node *Node) {
	ce.conditions = append(ce.conditions, conditionalExtensionCondition{
		node:                  node,
		conditionalExpression: condition,
	})
}

func ConditionalExtensionNodeProcessor(node *Node) {
	if hasAttribute, _, ifCondition := node.HasAttribute("go-if"); hasAttribute {
		condBranchSiblings := make([]*Node, 0)
		conditionalExtension := &ConditionalExtension{}

		node.RemoveAttribute("go-if")
		conditionalExtension.addCondition(ifCondition, node)

		node.NextSiblings(func(sibling *Node) bool {
			hasElifAttr, _, elifCondition := sibling.HasAttribute("go-elif")
			hasElseIfAttr, _, elseIfCondition := sibling.HasAttribute("go-else-if")

			if hasElifAttr {
				sibling.RemoveAttribute("go-elif")
				conditionalExtension.addCondition(elifCondition, sibling)
			} else if hasElseIfAttr {
				sibling.RemoveAttribute("go-else-if")
				conditionalExtension.addCondition(elseIfCondition, sibling)
			} else if hasElseAttr, _, _ := sibling.HasAttribute("go-else"); hasElseAttr {
				sibling.RemoveAttribute("go-else")
				conditionalExtension.elseNode = sibling

				condBranchSiblings = append(condBranchSiblings, sibling)

				return false
			} else {
				return false
			}

			condBranchSiblings = append(condBranchSiblings, sibling)
			return true
		})

		node.AddExtension(conditionalExtension)

		for _, condBranchSibling := range condBranchSiblings {
			condBranchSibling.Parent().RemoveChild(condBranchSibling)
		}
	}
}

type conditionalClassExtensionCondition struct {
	className             string
	conditionalExpression string
}

type ConditionalClassExtension struct {
	originalClasses    []string
	conditionalClasses []conditionalClassExtensionCondition
}

func (ce *ConditionalClassExtension) Apply(node *Node, dependencies ExtensionDependencies, params EvaluatorParams) (*Node, []*Node, error) {
	var appliedClasses []string
	copyNode := CopyNode(node)

	appliedClasses = append(appliedClasses, ce.originalClasses...)
	for _, conditionalClass := range ce.conditionalClasses {
		evaluator := dependencies.Get(EvaluatorExtDepKey).(Evaluator)
		hasResult, result, err := TryEvaluateBoolUsingContext(
			node, evaluator, conditionalClass.conditionalExpression,
		)
		if !hasResult {
			result, err = evaluator.EvaluateBool(conditionalClass.conditionalExpression, params)
		}
		if err != nil {
			return nil, nil, err
		} else if result {
			appliedClasses = append(appliedClasses, conditionalClass.className)
		}
	}
	if len(appliedClasses) > 0 {
		copyNode.AddAttribute("class", strings.Join(appliedClasses, " "))
	}
	return copyNode, nil, nil
}

func ConditionalClassExtensionNodeProcessor(node *Node) {
	ifClassAttrs := node.HasAttributes(func(attr Attribute) bool {
		return strings.HasPrefix(attr.Key, "go-if-class-")
	})
	if len(ifClassAttrs) > 0 {
		conditionalClassExtension := &ConditionalClassExtension{}

		hasClass, _, class := node.HasAttribute("class")
		if hasClass && class != "" {
			conditionalClassExtension.originalClasses = strings.Fields(class)
			node.RemoveAttribute("class")
		}

		for _, ifClassAttr := range ifClassAttrs {
			className := strings.TrimPrefix(ifClassAttr.Key, "go-if-class-")
			className = strings.TrimSpace(className)
			if className != "" {
				conditionalClassExtension.conditionalClasses = append(
					conditionalClassExtension.conditionalClasses,
					conditionalClassExtensionCondition{
						className:             className,
						conditionalExpression: ifClassAttr.Value,
					},
				)
				node.RemoveAttribute(ifClassAttr.Key)
			}
		}

		node.AddExtension(conditionalClassExtension)
	}
}

type RangeExtension struct {
	sourceVarName string

	isApplyingOnNewNodes bool
}

func (re *RangeExtension) Apply(node *Node, dependencies ExtensionDependencies, params EvaluatorParams) (*Node, []*Node, error) {
	if re.isApplyingOnNewNodes {
		return node, nil, nil
	}

	evaluator := dependencies.Get(EvaluatorExtDepKey).(Evaluator)
	hasResult, result, err := TryEvaluateUsingContext(node, evaluator, re.sourceVarName)
	if !hasResult {
		result, err = evaluator.Evaluate(re.sourceVarName, params)
	}
	if err != nil {
		return nil, nil, err
	}
	rangeEvalParams, isRangeEvalParam := result.(RangeEvaluatorParams)
	if !isRangeEvalParam {
		return nil, nil, fmt.Errorf("the type of `%s` is not RangeEvaluatorParams", re.sourceVarName)
	}

	re.isApplyingOnNewNodes = true

	var newNodes []*Node
	for _, rangeEvalParam := range rangeEvalParams {
		nodeCopy := CopyNode(node)
		nodeCopy.SetContextParams(rangeEvalParam)
		nodeCopy.SetParentEvaluatorContextSource(node.Parent())

		// ignore new siblings produced by this Node#ApplyExtensions func call
		newNodeCopy, _, err := nodeCopy.ApplyExtensions(dependencies, params)
		if err != nil {
			return nil, nil, err
		} else if newNodeCopy != nil {
			newNodes = append(newNodes, newNodeCopy)
		}
	}

	re.isApplyingOnNewNodes = false

	return nil, newNodes, nil
}

type RangeEvaluatorParams []EvaluatorParams

func RangeParams(params ...EvaluatorParams) RangeEvaluatorParams {
	return params
}

func RangeExtensionNodeProcessor(node *Node) {
	if hasRange, _, rangeDeclaration := node.HasAttribute("go-range"); hasRange {
		rangeExtension := &RangeExtension{
			sourceVarName: rangeDeclaration,
		}
		node.AddExtension(rangeExtension)
		node.RemoveAttribute("go-range")
	}
}

var stringInterpolationMarkerRegex = regexp.MustCompile("{{go:[a-zA-Z]+[a-zA-Z\\d\\.]*[a-zA-Z\\d]*}}")

type strInterpMarker struct {
	marker string
	key    string
}

type attrStrInterpMarkers struct {
	attributeKey string
	markers      []strInterpMarker
}

type AttrStringInterpExtension struct {
	markers []attrStrInterpMarkers
}

func (asie AttrStringInterpExtension) Apply(node *Node, dependencies ExtensionDependencies, params EvaluatorParams) (*Node, []*Node, error) {
	evaluator := dependencies.Get(EvaluatorExtDepKey).(Evaluator)

	for _, marker := range asie.markers {
		hasAttr, _, attrVal := node.HasAttribute(marker.attributeKey)
		if !hasAttr {
			return nil, nil, fmt.Errorf("attr string interp ext: assertion error. cannot find attr `%v`", marker.attributeKey)
		}
		for _, marker := range marker.markers {
			hasResult, result, err := TryEvaluateStringUsingContext(node, evaluator, marker.key)
			if !hasResult {
				result, err = evaluator.EvaluateString(marker.key, params)
			}
			if err != nil {
				return nil, nil, fmt.Errorf("attr string interp ext: %v", err)
			}
			attrVal = strings.Replace(attrVal, marker.marker, result, 1)
		}
		node.ReplaceAttribute(marker.attributeKey, attrVal)
	}

	return node, nil, nil
}

type TextStringInterpExtension struct {
	markers []strInterpMarker
}

func (tsie TextStringInterpExtension) Apply(node *Node, dependencies ExtensionDependencies, params EvaluatorParams) (*Node, []*Node, error) {
	evaluator := dependencies.Get(EvaluatorExtDepKey).(Evaluator)

	for _, marker := range tsie.markers {
		hasResult, result, err := TryEvaluateStringUsingContext(node, evaluator, marker.key)
		if !hasResult {
			result, err = evaluator.EvaluateString(marker.key, params)
		}
		if err != nil {
			return nil, nil, fmt.Errorf("text string interp ext: %v", err)
		}
		node.Data = strings.Replace(node.Data, marker.marker, result, 1)
	}

	return node, nil, nil
}

func StringInterpolationNodeProcessor(node *Node) {
	switch node.Type {
	case html.ElementNode:
		var attrMarkers []attrStrInterpMarkers
		for _, attr := range node.Attributes() {
			matches := stringInterpolationMarkerRegex.FindAllString(attr.Value, -1)
			if len(matches) > 0 {
				attrMarker := attrStrInterpMarkers{
					attributeKey: attr.Key,
				}
				for _, marker := range matches {
					key := strings.TrimPrefix(marker, "{{go:")
					key = strings.TrimSuffix(key, "}}")
					attrMarker.markers = append(attrMarker.markers, strInterpMarker{
						marker: marker,
						key:    key,
					})
				}
				attrMarkers = append(attrMarkers, attrMarker)
			}
		}
		if len(attrMarkers) > 0 {
			node.AddExtension(&AttrStringInterpExtension{
				markers: attrMarkers,
			})
		}
	case html.TextNode:
		matches := stringInterpolationMarkerRegex.FindAllString(node.Data, -1)
		if len(matches) > 0 {
			tsiExt := &TextStringInterpExtension{}
			for _, marker := range matches {
				key := strings.TrimPrefix(marker, "{{go:")
				key = strings.TrimSuffix(key, "}}")
				tsiExt.markers = append(tsiExt.markers, strInterpMarker{
					marker: marker,
					key:    key,
				})
			}
			node.AddExtension(tsiExt)
		}
	}
}
