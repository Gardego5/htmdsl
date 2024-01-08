//go:generate go run .
package main

import (
	"errors"
	"strings"
	"sync"
)

func main() {
	errs := []error{}
	errch := make(chan error)
	wg := sync.WaitGroup{}

	// aggregate errors
	go func() {
		for {
			errs = append(errs, <-errch)
			wg.Done()
		}
	}()

	// run each codegen func in a goroutine
	for _, fn := range []CodegenFunc{
		codegenerator("tags", tags),
	} {
		wg.Add(1)
		go func(fn func() error) {
			errch <- fn()
		}(fn)
	}

	// wait for all codegen funcs to finish
	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		panic(err)
	}
}

var tags = [...]struct {
	Name string
	Void bool
	Desc string
}{
	{Name: "a", Desc: strings.Join([]string{
		"The `a` HTML element (or *anchor* element), with its `href`",
		"attribute, creates a hyperlink to webpages, files, email",
		"addresses, locations in the same page, or anything else a URL",
		"can address.",
	}, "\n// ")},
	{Name: "abbr"},
	{Name: "address"},
	{Name: "area", Void: true},
	{Name: "article"},
	{Name: "aside"},
	{Name: "audio"},
	{Name: "b"},
	{Name: "base", Void: true},
	{Name: "bdi"},
	{Name: "bdo"},
	{Name: "blockquote"},
	{Name: "body"},
	{Name: "br", Void: true},
	{Name: "button"},
	{Name: "canvas"},
	{Name: "caption"},
	{Name: "cite"},
	{Name: "code"},
	{Name: "col", Void: true},
	{Name: "colgroup"},
	{Name: "data"},
	{Name: "datalist"},
	{Name: "dd"},
	{Name: "del"},
	{Name: "details"},
	{Name: "dfn"},
	{Name: "dialog"},
	{Name: "div"},
	{Name: "dl"},
	{Name: "dt"},
	{Name: "em"},
	{Name: "embed", Void: true},
	{Name: "fieldset"},
	{Name: "figcaption"},
	{Name: "figure"},
	{Name: "footer"},
	{Name: "form"},
	{Name: "h1"},
	{Name: "h2"},
	{Name: "h3"},
	{Name: "h4"},
	{Name: "h5"},
	{Name: "h6"},
	{Name: "head"},
	{Name: "header"},
	{Name: "hgroup"},
	{Name: "hr", Void: true},
	{Name: "html"},
	{Name: "i"},
	{Name: "iframe"},
	{Name: "img", Void: true},
	{Name: "input", Void: true},
	{Name: "ins"},
	{Name: "kbd"},
	{Name: "label"},
	{Name: "legend"},
	{Name: "li"},
	{Name: "link", Void: true},
	{Name: "main"},
	{Name: "map"},
	{Name: "mark"},
	{Name: "math"},
	{Name: "menu"},
	{Name: "menuitem", Void: true},
	{Name: "meta", Void: true},
	{Name: "meter"},
	{Name: "nav"},
	{Name: "noscript"},
	{Name: "object"},
	{Name: "ol"},
	{Name: "optgroup"},
	{Name: "option"},
	{Name: "output"},
	{Name: "p"},
	{Name: "param", Void: true},
	{Name: "picture"},
	{Name: "pre"},
	{Name: "progress"},
	{Name: "q"},
	{Name: "rb"},
	{Name: "rp"},
	{Name: "rt"},
	{Name: "rtc"},
	{Name: "ruby"},
	{Name: "s"},
	{Name: "samp"},
	{Name: "script"},
	{Name: "search"},
	{Name: "section"},
	{Name: "select"},
	{Name: "slot"},
	{Name: "small"},
	{Name: "source", Void: true},
	{Name: "span"},
	{Name: "strong"},
	{Name: "style"},
	{Name: "sub"},
	{Name: "summary"},
	{Name: "sup"},
	{Name: "svg"},
	{Name: "table"},
	{Name: "tbody"},
	{Name: "td"},
	{Name: "template"},
	{Name: "textarea"},
	{Name: "tfoot"},
	{Name: "th"},
	{Name: "thead"},
	{Name: "time"},
	{Name: "title"},
	{Name: "tr"},
	{Name: "track", Void: true},
	{Name: "u"},
	{Name: "ul"},
	{Name: "var"},
	{Name: "video"},
	{Name: "wbr", Void: true},
}
