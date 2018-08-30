package tplinator_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/bmdelacruz/tplinator"
	"github.com/yosssi/gohtml"
)

func ExampleTemplate_Execute() {
	sampleHtml := `
	<!doctype html>
	<html>
	<head>
		<title></title>
	</head>
	<body>
		<form action="/test/{{go:uid}}" method="POST">
			<input type="hidden" name="secret" value="{{go:secret}}"/>
			<input type="hidden" name="confcode" value="{{go:confcode}}"/>
		</form>
		<div class="container">
			<h1 go-if="!shouldBeRendered">INVERTED; This should appear conditionally.</h1>
			<h2>This should appear conditionally.</h2>
			<div class="col-sm-4" go-if="shouldBeRendered">
				<h3>This should appear conditionally.</h3>
				<h4>Username: {{go:username}}</h4>
			</div>
		</div>
	</body>
	</html>
	`

	tpl, err := tplinator.CreateTemplate(strings.NewReader(sampleHtml))
	if err != nil {
		log.Fatal(err)
	}

	docBytes, err := tpl.Execute(map[string]interface{}{
		"shouldBeRendered": true,

		"username": "bryanmdlx",
		"uid":      "28473664853",
		"secret":   "hImydjwKwixFa9yv08wRDw",
		"confcode": "Y3g1tT3sH-aHQ_rMJOSB7A",
	})
	if err != nil {
		log.Fatal(err)
	}

	gohtml.Condense = true
	fmt.Printf("%s\n", gohtml.FormatBytes(docBytes))

	// Output:
	// <!DOCTYPE html>
	// <html>
	//   <head>
	//     <title></title>
	//   </head>
	//   <body>
	//     <form action="/test/28473664853" method="POST">
	//       <input type="hidden" name="secret" value="hImydjwKwixFa9yv08wRDw"/>
	//       <input type="hidden" name="confcode" value="Y3g1tT3sH-aHQ_rMJOSB7A"/>
	//     </form>
	//     <div class="container">
	//       <h2>This should appear conditionally.</h2>
	//       <div class="col-sm-4">
	//         <h3>This should appear conditionally.</h3>
	//         <h4>Username: bryanmdlx</h4>
	//       </div>
	//     </div>
	//   </body>
	// </html>
}
