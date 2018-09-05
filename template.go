package tplinator

// import (
// 	"errors"
// 	"fmt"
// 	"log"
// 	"regexp"
// 	"strings"

// 	"github.com/Knetic/govaluate"
// )

// // PrecompiledTemplate is a struct which basically holds the information
// // about the original node structure of the parsed HTML document.
// type PrecompiledTemplate struct {
// 	documentNode *node
// }

// // Execute returns the bytes of the rendered node structure of the underlying
// // HTML document. It evaluates expressions, specified within the HTML document,
// // used by the templating features currently implemented in this package and
// // performs string interpolation on all placeholders found in the HTML document.
// // When the expression evaluation fails or when variable are not present in the
// // `data` parameter, it simply uses the default values, e.g. `false` for boolean
// // value requiring expression, empty string for string value requiring expressions,
// // etc., instead of returning an error and then logs a warning about it.
// func (pt *PrecompiledTemplate) Execute(data map[string]interface{}) ([]byte, error) {
// 	docStr, err := pt.documentNode.Execute(
// 		createInterpolatorFunc(data),
// 		createEvaluatorFunc(data),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return []byte(docStr), nil
// }

// // ExecuteStrict returns the bytes of the rendered node structure of the
// // underlying HTML document. It evaluates expressions, specified within the
// // HTML document, used by the templating features currently implemented in
// // this package and performs string interpolation on all placeholders found
// // in the HTML document. When the expression evaluation fails or when variable
// // are not present in the `data` parameter, it immediately stops the processing
// // of the template and returns an error.
// func (pt *PrecompiledTemplate) ExecuteStrict(data map[string]interface{}) ([]byte, error) {
// 	docStr, err := pt.documentNode.Execute(
// 		createStrictInterpolatorFunc(data),
// 		createStrictEvaluatorFunc(data),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return []byte(docStr), nil
// }

// type evaluator func(string) (interface{}, error)
// type interpolator func(string) (string, error)

// var errEvaluationFailed = errors.New("cannot evaluate provided expression")

// func createInterpolatorFunc(params map[string]interface{}) interpolator {
// 	return func(data string) (string, error) {
// 		return interpolate(data, params, false)
// 	}
// }

// func createStrictInterpolatorFunc(params map[string]interface{}) interpolator {
// 	return func(data string) (string, error) {
// 		return interpolate(data, params, true)
// 	}
// }

// func createEvaluatorFunc(params map[string]interface{}) evaluator {
// 	return func(data string) (interface{}, error) {
// 		return evaluate(data, params, false)
// 	}
// }

// func createStrictEvaluatorFunc(params map[string]interface{}) evaluator {
// 	return func(data string) (interface{}, error) {
// 		return evaluate(data, params, true)
// 	}
// }

// // Interpolation related constants
// const (
// 	interpolationStart = "{{go:"
// 	interpolationEnd   = "}}"
// )

// // Interpolation related global variables
// var (
// 	interpolationRegexPattern = regexp.MustCompile(interpolationStart + "[\\d\\w]+" + interpolationEnd)
// )

// func interpolate(expression string, params map[string]interface{}, strict bool) (string, error) {
// 	var hasWarning bool

// 	outputString := expression
// 	matches := interpolationRegexPattern.FindAllString(expression, -1)

// 	for _, match := range matches {
// 		originalMatch := match

// 		if !strict && hasWarning {
// 			// replace remaining matches with empty string because there's an error
// 			outputString = strings.Replace(outputString, originalMatch, "", 1)
// 			continue
// 		}

// 		match = strings.TrimLeft(match, interpolationStart)
// 		match = strings.TrimRight(match, interpolationEnd)
// 		varStr := strings.TrimSpace(match)

// 		expr, err := govaluate.NewEvaluableExpression(varStr)
// 		if err != nil {
// 			if !strict {
// 				hasWarning = true
// 				logInterpolationWarning(originalMatch, err)

// 				outputString = strings.Replace(outputString, originalMatch, "", 1)
// 				continue
// 			}
// 			return "", err
// 		}
// 		result, err := expr.Evaluate(params)
// 		if err != nil {
// 			if !strict {
// 				hasWarning = true
// 				logInterpolationWarning(originalMatch, err)

// 				outputString = strings.Replace(outputString, originalMatch, "", 1)
// 				continue
// 			}
// 			return "", err
// 		}

// 		outputString = strings.Replace(outputString, originalMatch, fmt.Sprintf("%v", result), 1)
// 	}

// 	return outputString, nil
// }

// func evaluate(expressionStr string, params map[string]interface{}, strict bool) (interface{}, error) {
// 	expr, err := govaluate.NewEvaluableExpression(expressionStr)
// 	if err != nil {
// 		if !strict {
// 			return nil, errEvaluationFailed
// 		}
// 		return nil, err
// 	}
// 	result, err := expr.Evaluate(params)
// 	if err != nil {
// 		if !strict {
// 			return nil, errEvaluationFailed
// 		}
// 		return nil, err
// 	}
// 	return result, nil
// }

// func logInterpolationWarning(expression string, cause error) {
// 	log.Printf("warning: failed to properly interpolate expression `%v`. cause: %v\n", expression, cause)
// }

// func logDefaultValueWarning(expression string, defaultValue interface{}) {
// 	log.Printf("warning: failed to properly evaluate expression `%v`. resolution: using default value (%v)\n", expression, defaultValue)
// }
