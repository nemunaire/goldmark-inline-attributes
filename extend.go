package attributes

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type Extender struct{}

var (
	defaultParser   = new(attributesParser)
	defaultRenderer = new(attributesRenderer)
)

func (e *Extender) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(defaultParser, 100),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(defaultRenderer, 100),
		),
	)
}

// Extension is a goldmark.Extender with markdown inline attributes support.
var Extension goldmark.Extender = new(Extender)

// Enable is a goldmark.Option with inline attributes support.
var Enable = goldmark.WithExtensions(Extension)
