package tplinator

// import (
// 	"io"

// 	"golang.org/x/net/html"
// )

// // Tplinate parses the provided HTML document using the `html.Parse`
// // function, duplicates the HTML document's node structure into a more
// // manipulatable node structure, and then returns a pointer to a
// // `PrecompiledTemplate` struct instance if successful.
// func Tplinate(documentReader io.Reader) (*PrecompiledTemplate, error) {
// 	documentNode, err := html.Parse(documentReader)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &PrecompiledTemplate{
// 		documentNode: precompileToNode(documentNode),
// 	}, nil
// }
