package tplinator

type Extension interface {
	Apply(node *Node, dependencies ExtensionDependencies, params EvaluatorParams) (*Node, error)
}

type conditionalExtensionCondition struct {
	node                  *Node
	conditionalExpression string
}

type ConditionalExtension struct {
	conditions []conditionalExtensionCondition
	elseNode   *Node
}

func (ce *ConditionalExtension) Apply(node *Node, dependencies ExtensionDependencies, params EvaluatorParams) (*Node, error) {
	for _, condition := range ce.conditions {
		evaluator := dependencies.Get(evaluatorExtDepKey).(Evaluator)
		result, err := evaluator.EvaluateBool(condition.conditionalExpression, params)
		if err != nil {
			return nil, err
		} else if result {
			return condition.node, nil
		}
	}
	return ce.elseNode, nil
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

func ConditionalClassExtensionNodeProcessor(node *Node) {
	// TODO
}

func RangeExtensionNodeProcessor(node *Node) {
	// TODO
}
