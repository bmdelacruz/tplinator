# tplinator

An HTML template renderer heavily inspired by the syntax used for Angular and VueJS applications.

[![Stage](https://img.shields.io/badge/experimental-red.svg)](https://img.shields.io/badge/experimental-red.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/bmdelacruz/tplinator)](https://goreportcard.com/report/github.com/bmdelacruz/tplinator)
[![Coverage Status](https://coveralls.io/repos/github/bmdelacruz/tplinator/badge.svg?branch=master)](https://coveralls.io/github/bmdelacruz/tplinator?branch=master)
[![CircleCI](https://circleci.com/gh/bmdelacruz/tplinator/tree/master.svg?style=svg)](https://circleci.com/gh/bmdelacruz/tplinator/tree/master)


## Features

### String Interpolation

Uses the `{{go:[EXPRESSION]}}` syntax which can be placed on attribute values and inside HTML elements. It is recommended to use simple expressions that will produce a string. 

#### Example

*Golang code snippet*

```golang
bytes, err := precompiledTemplate.Execute(map[string]interface{}{
    "username": "bryanmdlx",
    "currentDisplayName": "Bryan Dela Cruz",
})
```

*HTML template snippet*

```html
<h1>Username: {{go:username}}</h1>
<form action="/user/{{go:username}}/details" method="POST">
    <input type="text" name="displayName" value="{{go:currentDisplayName}}"/>
    <button type="submit">Update</button>
</form>
```

*Output HTML*

```html
<h1>Username: bryanmdlx</h1>
<form action="/user/bryanmdlx/details" method="POST">
    <input type="text" name="displayName" value="Bryan Dela Cruz"/>
    <button type="submit">Update</button>
</form>
```

### Conditional Rendering

Uses the `go-if`, `go-else-if` (or `go-elif`), and `go-else` to define that the target element/s will be rendered conditionally. The value of the conditional attribute must be a boolean expression; `PrecompiledTemplate#ExecuteStrict` will return an error if the result of the specified expression is not of boolean type while `PrecompiledTemplate#Execute` will simply ignore it and use the default value: `false`.

#### Example

*Golang code snippet*

```golang
bytes, err := precompiledTemplate.Execute(map[string]interface{}{
    "hasOne": false,
    "hasTwo": true,
    "hasThree": true,
})
```

*HTML template snippet*

```html
<div class="messages">
    <p go-if="hasOne">Has one</p>
    <p go-else-if="hasTwo && !hasThree">Has two but does not have three</p>
    <p go-else-if="hasTwo && hasThree">Has two and three</p>
    <p go-else>Does not have any</p>
</div>
```

*Output HTML*

```html
<div class="messages">
    <p>Has two and three</p>
</div>
```

### List Rendering

[![Not yet implemented](https://img.shields.io/badge/-Not%20yet%20implemented-red.svg)](https://img.shields.io/badge/-Not%20yet%20implemented-red.svg)

Uses the `go-range` attribute to define that the target element must be rendered `n` times, where `n` is the length of the specified range variable, under its original parent element.

#### Example

*Golang code snippet*

```golang
bytes, err := precompiledTemplate.Execute(map[string]interface{}{
    "foodMenu": tplinator.Range(
            foodMenuEntry{
            Name: "Blueberry Cheesecake",
            Description: "A delicious treat from the realm of gods and goddesses.",
        },
        foodMenuEntry{
            Name: "Iced Coffee",
            Description: "The tears of the trees that lives on the summit of the Alps.",
        },
    ),
})
```

*HTML template snippet*

```html
<div class="menu food">
    <div class="menu-entry" go-range="food in foodMenu">
        <h1>{{go:food.Name}}</h1>
        <p class="description">{{go:food.Description}}</p>
    </div>
</div>
```

*Output HTML*

```html
<div class="menu food">
    <div class="menu-entry">
        <h1>Blueberry Cheesecake</h1>
        <p class="description">A delicious treat from the realm of gods and goddesses.</p>
    </div>
    <div class="menu-entry">
        <h1>Iced Coffee</h1>
        <p class="description">The tears of the trees that lives on the summit of the Alps.</p>
    </div>
</div>
```

## Contributions

This is Go package is currently highly experimental. Contributions from y'all would be much appreciated.