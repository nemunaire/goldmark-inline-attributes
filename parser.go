package attributes

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type attributesParser struct {
}

var defaultAttributesParser = &attributesParser{}

// NewAttributesParser return a new InlineParser that parses inline attributes
// like '[txt]{.underline}' .
func NewAttributesParser() parser.InlineParser {
	return defaultAttributesParser
}

func (s *attributesParser) Trigger() []byte {
	return []byte{'['}
}

func (s *attributesParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	savedLine, savedPosition := block.Position()

	line, seg := block.PeekLine()

	endText := bytes.Index(line, []byte{']'})
	if endText < 0 {
		return nil // must close on the same line
	}

	if len(line) <= endText || line[endText+1] != '{' {
		return nil
	}
	block.Advance(endText + 1)

	attrs, ok := parser.ParseAttributes(block)
	if !ok {
		block.SetPosition(savedLine, savedPosition)
		return nil
	}

	n := &Node{}
	for _, attr := range attrs {
		n.SetAttribute(attr.Name, attr.Value)
	}

	n.AppendChild(n, ast.NewTextSegment(text.NewSegment(seg.Start+1, seg.Start+endText)))
	return n
}
