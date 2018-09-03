package tplinator

import (
	"io"

	"golang.org/x/net/html"
)

func Tplinate(documentReader io.Reader) (*PrecompiledTemplate, error) {
	documentNode, err := html.Parse(documentReader)
	if err != nil {
		return nil, err
	}
	return &PrecompiledTemplate{
		documentNode: precompileToNode(documentNode),
	}, nil
}
