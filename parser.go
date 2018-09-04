package tplinator

import (
	"errors"
	"io"
	"strings"

	"github.com/golang-collections/collections/stack"
	"golang.org/x/net/html"
)

type ParserOptionFunc func(*Parser)

func StrictnessParserOption(isStrict bool) ParserOptionFunc {
	return func(p *Parser) {
		p.isStrict = isStrict
	}
}

type parserNodeProcessorFunc func(*node)

func nodeProcessorsParserOption(npf []parserNodeProcessorFunc) ParserOptionFunc {
	return func(p *Parser) {
		p.nodeProcessorFuncs = npf
	}
}

type Parser struct {
	tokenizer *html.Tokenizer

	nodeProcessorFuncs []parserNodeProcessorFunc

	isStrict bool // TODO: use strictness value
}

// func ParseSample() { // TODO: REMOVE
// 	sample := `
// 		<h2>Hello, world!</h2>
// 		<h3>What is love?</h3>
// 	`

// 	nodes, err := parseTemplate(
// 		strings.NewReader(sample),
// 		nodeProcessorsParserOption([]parserNodeProcessorFunc{}),
// 	)
// 	if err != nil {
// 		fmt.Println("error: ", err)
// 	}
// 	spew.Dump(nodes)
// }

func parseTemplate(rdr io.Reader, opts ...ParserOptionFunc) ([]*node, error) {
	parser := Parser{
		tokenizer: html.NewTokenizer(rdr),
	}
	for _, parserOption := range opts {
		parserOption(&parser)
	}
	return parser.parse()
}

func (p Parser) parse() ([]*node, error) {
	templateNodes := make([]*node, 0)
	parserStack := stack.New()

	postProcessThenAddNode := func(newNode *node) {
		for _, processNode := range p.nodeProcessorFuncs {
			processNode(newNode)
		}
		if top := parserStack.Peek(); top != nil {
			top.(*node).AppendChild(newNode)
		} else {
			templateNodes = append(templateNodes, newNode)
		}
	}
	createAttributes := func(origAttrs []html.Attribute) []attribute {
		attributes := make([]attribute, len(origAttrs))
		for attrIdx, attr := range origAttrs {
			attributes[attrIdx] = attribute{
				key:   attr.Key,
				value: attr.Val,
			}
		}
		return attributes
	}

	for {
		tokenType := p.tokenizer.Next()
		err := p.tokenizer.Err()
		if err != nil {
			if err != io.EOF {
				return templateNodes, err
			} else if parserStack.Len() > 0 {
				return templateNodes, errors.New(
					"parser: reached the end of the file unexpectedly",
				)
			}
			return templateNodes, nil
		}

		switch tokenType {
		case html.ErrorToken, html.CommentToken:
			// simply ignore comment and error tokens
		case html.TextToken:
			// if the text token's data becomes an empty string
			// after trimming, do not add it to the current node
			// as its child or to the template nodes.
			token := p.tokenizer.Token()
			trimmedText := strings.TrimSpace(token.Data)
			if len(trimmedText) > 0 {
				postProcessThenAddNode(&node{
					nodeType: html.TextNode,
					data:     trimmedText,
				})
			}
		case html.DoctypeToken:
			// the doctype token should be the first token to be found
			// if the template is a complete HTML document
			if parserStack.Peek() != nil && len(templateNodes) > 0 {
				return templateNodes, errors.New(
					"parser: unexpectedly found a doctype",
				)
			}
			token := p.tokenizer.Token()
			templateNodes = append(templateNodes, &node{
				nodeType: html.DoctypeNode,
				data:     token.Data,
			})
		case html.SelfClosingTagToken:
			token := p.tokenizer.Token()
			postProcessThenAddNode(&node{
				nodeType:   html.ElementNode,
				data:       token.Data,
				attributes: createAttributes(token.Attr),

				isSelfClosingFlag: true,
			})
		case html.StartTagToken:
			token := p.tokenizer.Token()
			newNode := &node{
				nodeType:   html.ElementNode,
				data:       token.Data,
				attributes: createAttributes(token.Attr),
			}
			postProcessThenAddNode(newNode)
			parserStack.Push(newNode)
		case html.EndTagToken:
			if top := parserStack.Peek(); top != nil {
				currentNode := top.(*node)
				token := p.tokenizer.Token()

				// if the tag of the start tag and the tag
				// of the current end tag token does not match,
				// return an error
				if token.Data != currentNode.data {
					return templateNodes, errors.New(
						"parser: found an end tag that does not " +
							"match with the start tag",
					)
				}

				parserStack.Pop()
			} else {
				return templateNodes, errors.New(
					"parser: found an end tag but there's no " +
						"start tag available",
				)
			}
		default:
			return templateNodes, errors.New("parser: unknown token type")
		}
	}
}
