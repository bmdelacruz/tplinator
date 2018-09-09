# tplinator

An HTML template renderer heavily inspired by the syntax used for Angular and VueJS applications.

[![Stage](https://img.shields.io/badge/experimental-red.svg)](https://img.shields.io/badge/experimental-red.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/bmdelacruz/tplinator)](https://goreportcard.com/report/github.com/bmdelacruz/tplinator)
[![Coverage Status](https://coveralls.io/repos/github/bmdelacruz/tplinator/badge.svg?branch=master)](https://coveralls.io/github/bmdelacruz/tplinator?branch=master)
[![CircleCI](https://circleci.com/gh/bmdelacruz/tplinator/tree/master.svg?style=svg)](https://circleci.com/gh/bmdelacruz/tplinator/tree/master)


## Features

### String Interpolation

Uses the `{{go:[VARIABLE_NAME]}}` syntax which can be placed on attribute values and inside HTML elements. It only accepts variables which are of type string.

#### Example

*Golang code snippet*

```golang
bytes, err := template.RenderBytes(map[string]interface{}{
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

Uses the `go-if`, `go-else-if` (or `go-elif`), and `go-else` to define that the target element/s will be rendered conditionally. The value of the conditional attribute must be a boolean expression.

#### Example

*Golang code snippet*

```golang
bytes, err := template.RenderBytes(map[string]interface{}{
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

### Conditional Classes

Uses `go-if-class-*` to define that the target element's class will be added conditionally. The value of the conditional attribute must be a boolean expression.

#### Example

*Golang code snippet*

```golang
bytes, err := template.RenderBytes(map[string]interface{}{
    "isFood": true,
    "isDrink": false,
    "name": "Blueberry Cheesecake",
    "description": "A delicious treat from the realm of gods and goddesses.",
})
```

*HTML template snippet*

```html
<div class="menu-entry" go-if-class-food="isFood" go-if-class-drink="isDrink">
    <h1>{{go:name}}</h1>
    <p>{{go:description}}</p>
</div>
```

*Output HTML*

```html
<div class="menu-entry food">
    <h1>Blueberry Cheesecake</h1>
    <p>A delicious treat from the realm of gods and goddesses.</p>
</div>
```

### List Rendering

Uses the `go-range` attribute to define that the target element must be rendered `n` times, where `n` is the length of the specified range variable, under its original parent element.

Each entry on the `RangeParams` function will be accessible only to the target element and its children. If it does not contain a variable needed for one of the string interpolations inside the target element, the variable will be looked for on the `EvaluatorParams` passed to the `Template#Render*` function.

#### Example

*Golang code snippet*

```golang
bytes, err := template.RenderBytes(tplinator.EvaluatorParams{
    "userId": "412897373847523",
    "menuEntries": tplinator.RangeParams(
        tplinator.EvaluatorParams{
            "id": "399585827",
            "isFood": true,
            "isDrink": false,
            "favoriteWhat": "food",
            "name": "Blueberry Cheesecake",
            "description": "A delicious treat from the realm of gods and goddesses.",
        },
        tplinator.EvaluatorParams{
            "id": "518273743",
            "isFood": false,
            "isDrink": true,
            "favoriteWhat": "drink",
            "name": "Iced Coffee",
            "description": "The tears of the trees that lives on the summit of the Alps.",
        },
    ),
})
```

*HTML template snippet*

```html
<div class="menu">
    <div class="menu-entry" go-range="menuEntries" go-if-class-food="isFood" go-if-class-drink="isDrink">
        <h1>{{go:name}}</h1>
        <p class="description">{{go:description}}</p>
        <form action="/user/{{go:userId}}/favorite?what={{go:favoriteWhat}}&id={{go:id}}" method="POST">
            <button type="submit">Add {{go:name}} to my favorites</button>
        </form>
    </div>
</div>
```

*Output HTML*

```html
<div class="menu">
    <div class="menu-entry food">
        <h1>Blueberry Cheesecake</h1>
        <p class="description">A delicious treat from the realm of gods and goddesses.</p>
        <form action="/user/412897373847523/favorite?what=food&id=399585827" method="POST">
            <button type="submit">Add Blueberry Cheesecake to my favorites</button>
        </form>
    </div>
    <div class="menu-entry drink">
        <h1>Iced Coffee</h1>
        <p class="description">The tears of the trees that lives on the summit of the Alps.</p>
        <form action="/user/412897373847523/favorite?what=drink&id=518273743" method="POST">
            <button type="submit">Add Iced Coffee to my favorites</button>
        </form>
    </div>
</div>
```

## Contributions

This is Go package is currently highly experimental. Contributions from y'all would be much appreciated.