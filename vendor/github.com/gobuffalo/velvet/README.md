# Velvet [![GoDoc](https://godoc.org/github.com/gobuffalo/velvet?status.svg)](https://godoc.org/github.com/gobuffalo/velvet) [![Build Status](https://travis-ci.org/gobuffalo/velvet.svg?branch=master)](https://travis-ci.org/gobuffalo/velvet) [![Code Climate](https://codeclimate.com/github/gobuffalo/velvet/badges/gpa.svg)](https://codeclimate.com/github/gobuffalo/velvet)

Velvet is a templating package for Go. It bears a striking resemblance to "handlebars" based templates, there are a few small changes/tweaks, that make it slightly different.

## General Usage

If you know handlebars, you basically know how to use Velvet.

Let's assume you have a template (a string of some kind):

```handlebars
<!-- some input -->
<h1>{{ name }}</h1>
<ul>
  {{#each names}}
    <li>{{ @value }}</li>
  {{/each}}
</ul>
```

Given that string, you can render the template like such:

```go
ctx := velvet.NewContext()
ctx.Set("name", "Mark")
ctx.Set("names", []string{"John", "Paul", "George", "Ringo"})
s, err := velvet.Render(input, ctx)
if err != nil {
  // handle errors
}
```
Which would result in the following output:

```html
<h1>Mark</h1>
<ul>
  <li>John</li>
  <li>Paul</li>
  <li>George</li>
  <li>Ringo</li>
</ul>
```

## Helpers

### If Statements

What to do? Should you render the content, or not? Using Velvet's built in `if`, `else`, and `unless` helpers, let you figure it out for yourself.

```handlebars
{{#if true }}
  render this
{{/if}}
```

#### Else Statements

```handlebars
{{#if false }}
  won't render this
{{ else }}
  render this
{{/if}}
```

#### Unless Statements

```handlebars
{{#unless true }}
  won't render this
{{/unless}}
```

### Each Statements

Into everyone's life a little looping must happen. We can't avoid the need to write loops in applications, so Velvet helps you out by coming loaded with an `each` helper to iterate through `arrays`, `slices`, and `maps`.

#### Arrays

When looping through `arrays` or `slices`, the block being looped through will be access to the "global" context, as well as have four new variables available within that block:

* `@first` [`bool`] - is this the first pass through the iteration?
* `@last` [`bool`] - is this the last pass through the iteration?
* `@index` [`int`] - the counter of where in the loop you are, starting with `0`.
* `@value` - the current element in the array or slice that is being iterated over.

```handlebars
<ul>
  {{#each names}}
    <li>{{ @index }} - {{ @value }}</li>
  {{/each}}
</ul>
```

By using "block parameters" you can change the "key" of the element being accessed from `@value` to a key of your choosing.

```handlebars
<ul>
  {{#each names as |name|}}
    <li>{{ name }}</li>
  {{/each}}
</ul>
```

To change both the key and the index name you can pass two "block parameters"; the first being the new name for the index and the second being the name for the element.

```handlebars
<ul>
  {{#each names as |index, name|}}
    <li>{{ index }} - {{ name }}</li>
  {{/each}}
</ul>
```

#### Maps

Looping through `maps` using the `each` helper is also supported, and follows very similar guidelines to looping through `arrays`.

* `@first` [`bool`] - is this the first pass through the iteration?
* `@last` [`bool`] - is this the last pass through the iteration?
* `@key` - the key of the pair being accessed.
* `@value` - the value of the pair being accessed.

```handlebars
<ul>
  {{#each users}}
    <li>{{ @key }} - {{ @value }}</li>
  {{/each}}
</ul>
```

By using "block parameters" you can change the "key" of the element being accessed from `@value` to a key of your choosing.

```handlebars
<ul>
  {{#each users as |user|}}
    <li>{{ @key }} - {{ user }}</li>
  {{/each}}
</ul>
```

To change both the key and the value name you can pass two "block parameters"; the first being the new name for the key and the second being the name for the value.

```handlebars
<ul>
  {{#each users as |key, user|}}
    <li>{{ key }} - {{ user }}</li>
  {{/each}}
</ul>
```

### Other Builtin Helpers

* `json` - returns a JSON marshaled string of the value passed to it.
* `js_escape` - safely escapes a string to be used in a JavaScript bit of code.
* `html_escape` - safely escapes a string to be used in an HTML bit of code.
* `upcase` - upper cases the entire string passed to it.
* `downcase` - lower cases the entire string passed to it.
* `markdown` - converts markdown to HTML.
* `len` - returns the length of an array or slice

Velvet also imports all of the helpers found [https://github.com/markbates/inflect/blob/master/helpers.go](https://github.com/markbates/inflect/blob/master/helpers.go)

## Custom Helpers

No templating package would be complete without allowing for you to build your own, custom, helper functions.

### Return Values

The first thing to understand about building custom helper functions is their are a few "valid" return values:

#### `string`

Return just a `string`. The `string` will be HTML escaped, and deemed "not"-safe.

```go
func() string {
  return ""
}
```

#### `string, error`

Return a `string` and an error. The `string` will be HTML escaped, and deemed "not"-safe.

```go
func() (string, error) {
  return "", nil
}
```

#### `template.HTML`

[https://golang.org/pkg/html/template/#HTML](https://golang.org/pkg/html/https://golang.org/pkg/html/template/#HTMLlate/#HTML)

Return a `template.HTML` string. The `template.HTML` will **not** be HTML escaped, and will be deemed safe.

```go
func() template.HTML {
  return template.HTML("")
}
```


#### `template.HTML, error`

Return a `template.HTML` string and an error. The `template.HTML` will **not** be HTML escaped, and will be deemed safe.

```go
func() ( template.HTML, error ) {
  return template.HTML(""), error
}
```

### Input Values

Custom helper functions can take any type, and any number of arguments. There is an option last argument, [`velvet.HelperContext`](https://godoc.org/github.com/gobuffalo/velvet#HelperContext), that can be received. It's quite useful, and I would recommend taking it, as it provides you access to things like the context of the call, the block associated with the helper, etc...

### Registering Helpers

Custom helpers can be registered in one of two different places; globally and per template.

#### Global Helpers

```go
err := velvet.Helpers.Add("greet", func(name string) string {
  return fmt.Sprintf("Hi %s!", name)
})
if err != nil {
  // handle errors
}
```

The `greet` function is now available to all templates that use Velvet.

```go
s, err := velvet.Render(`<h1>{{greet "mark"}}</h1>`, velvet.NewContext())
if err != nil {
  // handle errors
}
fmt.Print(s) // <h1>Hi mark!</h1>
```

#### Per Template Helpers

```go
t, err := velvet.Parse(`<h1>{{greet "mark"}}</h1>`)
if err != nil {
  // handle errors
}
t.Helpers.Add("greet", func(name string) string {
  return fmt.Sprintf("Hi %s!", name)
})
if err != nil {
  // handle errors
}
```

The `greet` function is now only available to the template it was added to.

```go
s, err := t.Exec(velvet.NewContext())
if err != nil {
  // handle errors
}
fmt.Print(s) // <h1>Hi mark!</h1>
```

### Block Helpers

Like the `if` and `each` helpers, block helpers take a "block" of text that can be evaluated and potentially rendered, manipulated, or whatever you would like. To write a block helper, you have to take the `velvet.HelperContext` as the last argument to your helper function. This will give you access to the block associated with that call.

#### Example

```go
velvet.Helpers.Add("upblock", func(help velvet.HelperContext) (template.HTML, error) {
  s, err := help.Block()
  if err != nil {
    return "", err
  }
  return strings.ToUpper(s), nil
})

s, err := velvet.Render(`{{#upblock}}hi{{/upblock}}`, velvet.NewContext())
if err != nil {
  // handle errors
}
fmt.Print(s) // HI
```

