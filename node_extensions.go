package tplinator

import (
	"strings"
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
		evaluator := dependencies.Get(evaluatorExtDepKey).(Evaluator)
		result, err := evaluator.EvaluateBool(condition.conditionalExpression, params)
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
	copyNode := *node

	appliedClasses = append(appliedClasses, ce.originalClasses...)
	for _, conditionalClass := range ce.conditionalClasses {
		evaluator := dependencies.Get(evaluatorExtDepKey).(Evaluator)
		result, err := evaluator.EvaluateBool(conditionalClass.conditionalExpression, params)
		if err != nil {
			return nil, nil, err
		} else if result {
			appliedClasses = append(appliedClasses, conditionalClass.className)
		}
	}
	if len(appliedClasses) > 0 {
		copyNode.AddAttribute("class", strings.Join(appliedClasses, " "))
	}
	return &copyNode, nil, nil
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

func RangeExtensionNodeProcessor(node *Node) {
	// TODO
}
