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
		<title>Hello, world!</title>
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
	//     <title>Hello, world!</title>
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

func ExampleConditionalRendering() {
	sampleHtml := `
	<div class="messages">
		<p go-if="hasOne">Has one</p>
		<p go-else-if="hasTwo && !hasThree">Has two but does not have three</p>
		<p go-else-if="hasTwo && hasThree">Has two and three</p>
		<p go-else>Does not have any</p>
	</div>
	`

	template, err := tplinator.Tplinate(strings.NewReader(sampleHtml))
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := template.RenderBytes(map[string]interface{}{
		"hasOne":   false,
		"hasTwo":   true,
		"hasThree": true,
	})
	if err != nil {
		log.Fatal(err)
	}

	gohtml.Condense = true
	fmt.Printf("%s\n", gohtml.FormatBytes(bytes))

	// Output:
	// <div class="messages">
	//   <p>Has two and three</p>
	// </div>
}

func ExampleConditionalClasses() {
	sampleHtml := `
	<div class="menu-entry" go-if-class-food="isFood" go-if-class-drink="isDrink">
		<h1>{{go:name}}</h1>
		<p>{{go:description}}</p>
	</div>
	`

	template, err := tplinator.Tplinate(strings.NewReader(sampleHtml))
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := template.RenderBytes(map[string]interface{}{
		"isFood":      true,
		"isDrink":     false,
		"name":        "Blueberry Cheesecake",
		"description": "A delicious treat from the realm of gods and goddesses.",
	})
	if err != nil {
		log.Fatal(err)
	}

	gohtml.Condense = true
	fmt.Printf("%s\n", gohtml.FormatBytes(bytes))

	// Output:
	// <div class="menu-entry food">
	//   <h1>Blueberry Cheesecake</h1>
	//   <p>A delicious treat from the realm of gods and goddesses.</p>
	// </div>
}

func ExampleListRendering() {
	// FIXME
	// since range elements are not attached to the dom tree,
	// that element and its children will not be able to access
	// context params of the elements outside their scope.

	sampleHtml := `
	<div class="menu">
		<div class="menu-entry" go-range="menuEntries" go-if-class-food="isFood" go-if-class-drink="isDrink">
			<h1>{{go:name}}</h1>
			<p class="description">{{go:description}}</p>
			<form action="/user/{{go:userId}}/favorite?what={{go:favoriteWhat}}&id={{go:id}}" method="POST">
				<button type="submit">Add {{go:name}} to my favorites</button>
			</form>
		</div>
	</div>
	`

	template, err := tplinator.Tplinate(strings.NewReader(sampleHtml))
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := template.RenderBytes(tplinator.EvaluatorParams{
		"userId": "412897373847523",
		"menuEntries": tplinator.RangeParams(
			tplinator.EvaluatorParams{
				"isFood":       true,
				"isDrink":      false,
				"id":           "399585827",
				"favoriteWhat": "food",
				"name":         "Blueberry Cheesecake",
				"description":  "A delicious treat from the realm of gods and goddesses.",
			},
			tplinator.EvaluatorParams{
				"isFood":       false,
				"isDrink":      true,
				"id":           "518273743",
				"favoriteWhat": "drink",
				"name":         "Iced Coffee",
				"description":  "The tears of the trees that lives on the summit of the Alps.",
			},
		),
	})
	if err != nil {
		log.Fatal(err)
	}

	gohtml.Condense = true
	fmt.Printf("%s\n", gohtml.FormatBytes(bytes))

	// Output:
	// <div class="menu">
	//   <div class="menu-entry food">
	//     <h1>Blueberry Cheesecake</h1>
	//     <p class="description">A delicious treat from the realm of gods and goddesses.</p>
	//     <form action="/user/412897373847523/favorite?what=food&id=399585827" method="POST">
	//       <button type="submit">Add Blueberry Cheesecake to my favorites</button>
	//     </form>
	//   </div>
	//   <div class="menu-entry drink">
	//     <h1>Iced Coffee</h1>
	//     <p class="description">The tears of the trees that lives on the summit of the Alps.</p>
	//     <form action="/user/412897373847523/favorite?what=drink&id=518273743" method="POST">
	//       <button type="submit">Add Iced Coffee to my favorites</button>
	//     </form>
	//   </div>
	// </div>
}
