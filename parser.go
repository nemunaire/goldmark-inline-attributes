package attributes

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var spanContentStateKey = parser.NewContextKey()

type attributesParser struct {
}

var defaultAttributesParser = &attributesParser{}

// NewAttributesParser return a new InlineParser that parses inline attributes
// like '[txt]{.underline}' .
func NewAttributesParser() parser.InlineParser {
	return defaultAttributesParser
}

func (s *attributesParser) Trigger() []byte {
	return []byte{'[', ']'}
}

var attributeBottom = parser.NewContextKey()

func (s *attributesParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	sl, ss := block.Position()
	line, seg := block.PeekLine()

	if line[0] == '[' {
		openIdx := bytes.IndexByte(line[1:], '[')
		closeIdx := bytes.IndexByte(line, ']')
		if (openIdx < 0 || openIdx > closeIdx) && closeIdx > 0 && closeIdx+1 < len(line) && line[closeIdx+1] == '{' {
			pc.Set(attributeBottom, pc.LastDelimiter())
			return processSpanContentOpen(block, seg.Start, pc)
		}
	}

	// line[0] == ']'
	tlist := pc.Get(spanContentStateKey)
	if tlist == nil {
		block.SetPosition(sl, ss)
		return nil
	}

	last := tlist.(*spanContentState).Last
	if last == nil {
		return nil
	}
	block.Advance(1)

	removeSpanContentState(pc, last)

	c := block.Peek()
	if c != '{' {
		block.SetPosition(sl, ss)
		return nil
	}

	attrs, ok := parser.ParseAttributes(block)
	if !ok {
		block.SetPosition(sl, ss)
		return nil
	}

	span := &Node{}
	for _, attr := range attrs {
		span.SetAttribute(attr.Name, attr.Value)
	}

	var bottom ast.Node
	if v := pc.Get(attributeBottom); v != nil {
		bottom = v.(ast.Node)
	}
	pc.Set(attributeBottom, nil)
	parser.ProcessDelimiters(bottom, pc)
	for c := last.NextSibling(); c != nil; {
		next := c.NextSibling()
		parent.RemoveChild(parent, c)
		span.AppendChild(span, c)
		c = next
	}

	last.Parent().RemoveChild(last.Parent(), last)
	return span
}

type spanContentState struct {
	ast.BaseInline

	Segment text.Segment

	Prev *spanContentState

	Next *spanContentState

	First *spanContentState

	Last *spanContentState
}

func (s *spanContentState) Text(source []byte) []byte {
	return s.Segment.Value(source)
}

func (s *spanContentState) Dump(source []byte, level int) {
	fmt.Printf("%sspanContentState: \"%s\"\n", strings.Repeat("    ", level), s.Text(source))
}

var kindSpanContentState = ast.NewNodeKind("SpanContentState")

func (s *spanContentState) Kind() ast.NodeKind {
	return kindSpanContentState
}

func spanContentStateLength(v *spanContentState) int {
	if v == nil || v.Last == nil || v.First == nil {
		return 0
	}
	return v.Last.Segment.Stop - v.First.Segment.Start
}

func pushSpanContentState(pc parser.Context, v *spanContentState) {
	tlist := pc.Get(spanContentStateKey)
	var list *spanContentState
	if tlist == nil {
		list = v
		v.First = v
		v.Last = v
		pc.Set(spanContentStateKey, list)
	} else {
		list = tlist.(*spanContentState)
		l := list.Last
		list.Last = v
		l.Next = v
		v.Prev = l
	}
}

func removeSpanContentState(pc parser.Context, d *spanContentState) {
	tlist := pc.Get(spanContentStateKey)
	var list *spanContentState
	if tlist == nil {
		return
	}
	list = tlist.(*spanContentState)

	if d.Prev == nil {
		list = d.Next
		if list != nil {
			list.First = d
			list.Last = d.Last
			list.Prev = nil
			pc.Set(spanContentStateKey, list)
		} else {
			pc.Set(spanContentStateKey, nil)
		}
	} else {
		d.Prev.Next = d.Next
		if d.Next != nil {
			d.Next.Prev = d.Prev
		}
	}
	if list != nil && d.Next == nil {
		list.Last = d.Prev
	}
	d.Next = nil
	d.Prev = nil
	d.First = nil
	d.Last = nil
}

func processSpanContentOpen(block text.Reader, pos int, pc parser.Context) *spanContentState {
	state := &spanContentState{
		Segment: text.NewSegment(pos, pos+1),
	}
	pushSpanContentState(pc, state)
	block.Advance(1)
	return state
}
