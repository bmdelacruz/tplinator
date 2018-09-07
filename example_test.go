package tplinator_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/bmdelacruz/tplinator"
	"github.com/yosssi/gohtml"
)

func ExampleTplinate() {
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
			<h1 go-if="!shouldBeRendered">shouldBeRendered is equal to false</h1>
			<h2>shouldBeRendered is equal to true</h2>
			<div class="col-sm-4" go-if="shouldBeRendered">
				<h3>shouldBeRendered is equal to true</h3>
				<h4>Username: {{go:username}}</h4>
			</div>
			<ul>
				<li go-range="indexes">{{go:index}}</li>
			</ul>
		</div>
	</body>
	</html>
	`

	tpl, err := tplinator.Tplinate(strings.NewReader(sampleHtml))
	if err != nil {
		log.Fatal(err)
	}

	docBytes, err := tpl.RenderBytes(map[string]interface{}{
		"shouldBeRendered": true,

		"username": "bryanmdlx",
		"uid":      "28473664853",
		"secret":   "hImydjwKwixFa9yv08wRDw",
		"confcode": "Y3g1tT3sH-aHQ_rMJOSB7A",

		"indexes": tplinator.RangeParams(
			tplinator.EvaluatorParams{"index": "1"},
			tplinator.EvaluatorParams{"index": "2"},
			tplinator.EvaluatorParams{"index": "3"},
			tplinator.EvaluatorParams{"index": "4"},
			tplinator.EvaluatorParams{"index": "5"},
			tplinator.EvaluatorParams{"index": "6"},
		),
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
	//       <h2>shouldBeRendered is equal to true</h2>
	//       <div class="col-sm-4">
	//         <h3>shouldBeRendered is equal to true</h3>
	//         <h4>Username: bryanmdlx</h4>
	//       </div>
	//       <ul>
	//         <li>1</li>
	//         <li>2</li>
	//         <li>3</li>
	//         <li>4</li>
	//         <li>5</li>
	//         <li>6</li>
	//       </ul>
	//     </div>
	//   </body>
	// </html>
}
