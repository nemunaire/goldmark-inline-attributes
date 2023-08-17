package attributes

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

// Kind is the kind of hashtag AST nodes.
var Kind = ast.NewNodeKind("InlineAttributes")

// Node is a parsed attributes node.
type Node struct {
	ast.BaseInline
}

func (*Node) Kind() ast.NodeKind { return Kind }

func (n *Node) Dump(src []byte, level int) {
	attrs := n.Attributes()
	list := make(map[string]string, len(attrs))
	for _, attr := range attrs {
		name := util.BytesToReadOnlyString(attr.Name)
		value := util.BytesToReadOnlyString(util.EscapeHTML(attr.Value.([]byte)))
		list[name] = value
	}

	ast.DumpHelper(n, src, level, list, nil)
}
