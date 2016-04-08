[![Build Status](https://travis-ci.org/valyala/quicktemplate.svg)](https://travis-ci.org/valyala/quicktemplate)
[![GoDoc](https://godoc.org/github.com/valyala/quicktemplate?status.svg)](http://godoc.org/github.com/valyala/quicktemplate)
[![Coverage](http://gocover.io/_badge/github.com/valyala/quicktemplate)](http://gocover.io/github.com/valyala/quicktemplate)
[![Go Report](http://goreportcard.com/badge/valyala/quicktemplate)](http://goreportcard.com/report/valyala/quicktemplate)

# quicktemplate

Fast, powerful, yet easy to use template engine for Go.
Inspired by [mako templates](http://www.makotemplates.org/) philosophy.


# Features

  * [Extremely fast](#performance-comparison-with-htmltemplate).
    Templates are converted into Go code and then compiled.
  * Almost all the bugs are caught during template compilation, so production
    suffers less from template-related bugs.
  * Easy to use. See [quickstart](#quick-start) and [examples](https://github.com/valyala/quicktemplate/tree/master/examples)
    for details.
  * Powerful. Arbitrary Go code may be embedded into and mixed with templates.
    Be careful with this power - do not query db and/or external resources from
    templates unless you miss php way in Go :) This power is mostly for
    arbitrary data transformations.
  * Easy to understand template inheritance powered by [Go interfaces](https://golang.org/doc/effective_go.html#interfaces).
    See [this example](https://github.com/valyala/quicktemplate/tree/master/examples/basicserver) for details.
  * Templates are compiled into a single binary, so there is no need in copying
    template files to the server.

# Drawbacks

  * Templates cannot be updated on the fly on the server, since they
    are compiled into a single binary.
    Take a look at [fasttemplate](https://github.com/valyala/fasttemplate)
    if you need fast template engine for simple dynamically updated templates.

# Performance comparison with html/template

Quicktemplate is more than 20x faster than [html/template](https://golang.org/pkg/html/template/).
The following simple template is used in the benchmark:

  * [html/template version](https://github.com/valyala/quicktemplate/blob/master/testdata/templates/bench.tpl)
  * [quicktemplate version](https://github.com/valyala/quicktemplate/blob/master/testdata/templates/bench.qtpl)

Benchmark results:

```
$ go test -bench=Template -benchmem
BenchmarkQuickTemplate1-4  	10000000	       158 ns/op	       0 B/op	       0 allocs/op
BenchmarkQuickTemplate10-4 	 2000000	       604 ns/op	       0 B/op	       0 allocs/op
BenchmarkQuickTemplate100-4	  300000	      5498 ns/op	       0 B/op	       0 allocs/op
BenchmarkHTMLTemplate1-4   	  500000	      2807 ns/op	     752 B/op	      23 allocs/op
BenchmarkHTMLTemplate10-4  	  100000	     13527 ns/op	    3521 B/op	     117 allocs/op
BenchmarkHTMLTemplate100-4 	   10000	    133503 ns/op	   34499 B/op	    1152 allocs/op
```

# Security

By default all the template placeholders are html-escaped.

# Examples

See [examples](https://github.com/valyala/quicktemplate/tree/master/examples).

# Quick start

Let's start with a minimal template example:

```qtpl
All the text outside function templates is treated as comments,
i.e. it is just ignored by quicktemplate compiler (qtc). It is for humans.

Hello is a simple template function.
{% func Hello(name string) %}
	Hello, {%s name %}!
{% endfunc %}
```

Save this file into `templates` folder under the name `hello.qtpl`
and run [qtc](https://github.com/valyala/quicktemplate/tree/master/qtc)
inside this folder. `qtc` may be installed by issuing:

```
go get -u github.com/valyala/quicktemplate/qtc
```

If all went ok, `hello.qtpl.go` file must appear in the `templates` folder.
This file contains Go code for `hello.qtpl`. Let's use it!

Create a file main.go outside `templates` folder and put the following
code there:

```go
package main

import (
	"fmt"

	"./templates"
)

func main() {
	fmt.Printf("%s\n", templates.Hello("Foo"))
	fmt.Printf("%s\n", templates.Hello("Bar"))
}
```

Then run `go run`. If all went ok, you'll see something like this:

```

	Hello, Foo!


	Hello, Bar!

```

Let's create more complex template, which calls other template functions,
contains loops, conditions, breaks, continue and returns.
Put the following template into `templates/greetings.qtpl`:

```qtpl

Greetings greets up to 42 names.
It also greets John differently comparing to others.
{% func Greetings(names []string) %}
	{% if len(names) == 0 %}
		Nobody to greet :(
		{% return %}
	{% endif %}

	{% for i, name := range names %}
		{% if i == 42 %}
			I'm tired to greet so many people...
			{% break %}
		{% elseif name == "John" %}
			{%= sayHi("Mr. " + name) %}
			{% continue %}
		{% else%}
			{%= Hello(name) %}
		{% endif %}
	{% endfor %}
{% endfunc %}

sayHi is unexported, since it starts with lowercase letter.
{% func sayHi(name string) %}
	Hi, {%s name %}
{% endfunc %}

Note that every template file may contain arbitrary number
of template functions. For instance, this file contains Greetings and sayHi
functions.
```

Run `qtc` inside `templates` folder. Now the folder should contain
two files with Go code: `hello.qtpl.go` and `greetings.qtpl.go`. These files
form a single `templates` Go package. Template functions and other template
stuff is shared between template files located in the same folder.
So `Hello` template function may be used inside `greetings.qtpl` while
it is defined in `hello.qtpl`.
Moreover, the folder may contain ordinary Go files and its' contents may
be used inside templates.

Now put the following code into `main.go`:

```go
package main

import (
	"bytes"
	"fmt"

	"./templates"
)

func main() {
	names := []string{"Kate", "Go", "John", "Brad"}

	// qtc creates Write* function for each template function.
	// Such functions accept io.Writer as first parameter:
	var buf bytes.Buffer
	templates.WriteGreetings(&buf, names)

	fmt.Printf("buf=\n%s", buf.Bytes())
}
```

Careful readers may notice different output tags were used in these
templates: `{%s name %}` and `{%= Hello(name) %}`. What's the difference?
The `{%s x %}` is used for printing html-safe strings, while `{%= F() %}`
is used for embedding template function calls. Quicktemplate supports also
other output tags:

  * `{%d num %}` for integers
  * `{%f float %}` for float64
  * `{%z bytes %}` for byte slices
  * `{%q str %}` for json-compatible quoted strings.
  * `{%j str %}` for embedding str into json string. Unlike `{%q str %}`
    it doesn't quote the string.
  * `{%u str %}` for [URL encoding](https://en.wikipedia.org/wiki/Percent-encoding) the given str.
  * `{%v anything %}` is equivalent to `%v` in [printf-like functions](https://golang.org/pkg/fmt/).

All these output tags produce html-safe output, i.e. they escape `<` to `&lt;`,
`>` to `&gt;`, etc. If you don't want html-safe output, then just put `=` after
the tag. For example: `{%s= "<h1>This h1 won't be escaped</h1>" %}`.

As you may notice `{%= F() %}` and `{%s= F() %}` produce the same output for `{% func F() %}`.
But the first one is optimized for speed - it avoids memory allocations and copy.
So stick to it when embedding template function calls.

All the ouptut tags except of `{%= F() %}` may contain arbitrary valid
Go expression instead of just identifier. For example:

```qtpl
Import fmt for fmt.Sprintf()
{% import "fmt" %}

FmtFunc uses fmt.Sprintf() inside output tag
{% func FmtFunc(s string) %}
	{%s fmt.Sprintf("FmtFunc accepted %q string", s) %}
{% endfunc %}
```

There are other useful tags supported by quicktemplate:

  * `{% comment %}`

    ```qtpl
    {% comment %}
        This is a comment. It won't trap into the output.
        It may contain {% arbitrary tags %}. They are just ignored.
    {% endcomment %}
    ```

  * `{% plain %}`

  ```qtpl
      {% plain %}
          Tags will {% trap into %} the output {% unmodified %}.
          Plain block may contain invalid and {% incomplete tags.
      {% endplain %}
    ```

  * `{% collapsespace %}`

    ```qtpl
    {% collapsespace %}
        <div>
            <div>space between lines</div>
               and {%s " tags" %}
             <div>is collapsed into a single space
             unless{% newline %}or{% space %}is used</div>
        </div>
    {% endcollapsespace %}
    ```

    Is converted into

    ```
    <div> <div>space between lines</div> and tags <div>is collapsed into a single space unless
    or is used</div> </div>
    ```

  * `{% stripspace %}`

    ```qtpl
    {% stripspace %}
         <div>
             <div>space between lines</div>
                and {%s " tags" %}
             <div>is removed unless{% newline %}or{% space %}is used</div>
         </div>
    {% endstripspace %}
    ```

    Is converted into

    ```
    <div><div>space between lines</div>and tags<div>is removed unless
    or is used</div></div>
    ```

  * `{% code %}`:

    ```qtpl
    {% code
    // arbitrary Go code may be embedded here!
    type FooArg struct {
        Name string
        Age int
    }
    %}
    ```

  * `{% import %}`:

    ```qtpl
    Import external packages.
    {% import "foo/bar" %}
    {% import (
        "foo"
        bar "baz/baa"
    ) %}
    ```

  * `{% interface %}`:

    ```qtpl
    Interfaces allow powerful templates' inheritance
    {%
    interface Page {
        Title()
        Body(s string, n int)
        Footer()
    }
    %}

    PrintPage prints Page
    {% func PrintPage(p Page) %}
        <html>
            <head><title>{%= p.Title() %}</title></head>
            <body>
                <div>{%= p.Body("foo", 42) %}</div>
                <div>{%= p.Footer() %}</div>
            </body>
        </html>
    {% endfunc %}

    Base page implementation
    {% code
    type BasePage struct {
        TitleStr string
        FooterStr string
    }
    %}
    {% func (bp *BasePage) Title() %}{%s bp.TitleStr %}{% endfunc %}
    {% func (bp *BasePage) Body(s string, n int) %}
        <b>s={%q s %}, n={%d n %}</b>
    {% endfunc %}
    {% func (bp *BasePage) Footer() %}{%s bp.FooterStr %}{% endfunc %}

    Main page implementation
    {% code
    type MainPage struct {
        // inherit from BasePage
        BasePage

        // real body for main page
        BodyStr string
    }
    %}

    Override only Body
    Title and Footer are used from BasePage.
    {% func (mp *MainPage) Body(s string, n int) %}
        <div>
            main body: {%s mp.BodyStr %}
        </div>
        <div>
            base body: {%= mp.BasePage.Body(s, n) %}
        </div>
    {% endfunc %}
    ```

    See [basicserver example](https://github.com/valyala/quicktemplate/tree/master/examples/basicserver)
    for more details.


# FAQ

  * *Why quicktemplate syntax is incompatible with [html/template](https://golang.org/pkg/html/template/)?*

    Because `html/template` syntax isn't expressive enough for `quicktemplate`.

  * *What's the difference between quicktemplate and [ego](https://github.com/benbjohnson/ego)?*

    `Ego` is similar to `quicktemplate` in the sense it converts templates into Go code.
    But it misses the following stuff, which makes `quicktemplate` so powerful
    and easy to use:

      * Defining multiple function templates in a single template file.
      * Embedding function templates inside other function templates.
      * Template interfaces, inheritance and overriding.
        See [this example](https://github.com/valyala/quicktemplate/tree/master/examples/basicserver)
        for details.
      * Top-level comments outside function templates.
      * Template packages.
      * Combining arbitrary Go files with template files in template packages.
      * Performance optimizations.

  * *What's the difference between quicktemplate and [gorazor](https://github.com/sipin/gorazor)?*

    `Gorazor` is similar to `quicktemplate` in the sense it converts templates into Go code.
    But it misses the following useful features:

      * Clear syntax insead of hard-to-understand magic stuff related
        to template arguments, template inheritance and embedding function
        templates into other templates.
      * Performance optimizations.

* *I didn't find an answer for my question here*

  Try exploring [these questions](https://github.com/valyala/quicktemplate/issues?q=label%3Aquestion).
