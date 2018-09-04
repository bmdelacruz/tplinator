package tplinator

import (
	"errors"
	"io"
	"strings"

	"github.com/alediaferia/stackgo"
	"golang.org/x/net/html"
)

type ParserOptionFunc func(*Parser)

func StrictnessParserOption(isStrict bool) ParserOptionFunc {
	return func(p *Parser) {
		p.isStrict = isStrict
	}
}

type NodeProcessorFunc func(*Node)

func NodeProcessorsParserOption(npfs ...NodeProcessorFunc) ParserOptionFunc {
	return func(p *Parser) {
		p.nodeProcessors = append(p.nodeProcessors, npfs...)
	}
}

type Parser struct {
	tokenizer *html.Tokenizer

	nodeProcessors []NodeProcessorFunc

	isStrict bool // TODO: use strictness value
}

func ParseNodes(rdr io.Reader, opts ...ParserOptionFunc) ([]*Node, error) {
	parser := Parser{
		tokenizer: html.NewTokenizer(rdr),
	}
	for _, parserOption := range opts {
		parserOption(&parser)
	}
	return parser.parse()
}

func (p Parser) parse() ([]*Node, error) {
	parserStack := stackgo.NewStack()
	templateNodes := make([]*Node, 0)

	postProcessThenAddNode := func(newNode *Node) {
		for _, processNode := range p.nodeProcessors {
			processNode(newNode)
		}
		if top := parserStack.Top(); top != nil {
			top.(*Node).AppendChild(newNode)
		} else {
			templateNodes = append(templateNodes, newNode)
		}
	}

	for {
		tokenType := p.tokenizer.Next()

		err := p.tokenizer.Err()
		if err != nil {
			if err != io.EOF {
				return templateNodes, err
			} else if parserStack.Size() > 0 {
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
				postProcessThenAddNode(
					CreateNode(
						html.TextNode, trimmedText, nil, false,
					),
				)
			}
		case html.DoctypeToken:
			// the doctype token should be the first token to be found
			// if the template is a complete HTML document
			if parserStack.Top() != nil && len(templateNodes) > 0 {
				return templateNodes, errors.New(
					"parser: unexpectedly found a doctype",
				)
			}
			token := p.tokenizer.Token()
			templateNodes = append(templateNodes,
				CreateNode(
					html.DoctypeNode, token.Data, nil, false,
				),
			)
		case html.SelfClosingTagToken:
			token := p.tokenizer.Token()
			postProcessThenAddNode(
				CreateNode(
					html.ElementNode, token.Data, token.Attr, true,
				),
			)
		case html.StartTagToken:
			token := p.tokenizer.Token()
			newNode := CreateNode(
				html.ElementNode, token.Data, token.Attr, false,
			)
			postProcessThenAddNode(newNode)
			parserStack.Push(newNode)
		case html.EndTagToken:
			if top := parserStack.Top(); top != nil {
				currentNode := top.(*Node)
				token := p.tokenizer.Token()

				// if the tag of the start tag and the tag
				// of the current end tag token does not match,
				// return an error
				if token.Data != currentNode.Data {
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
