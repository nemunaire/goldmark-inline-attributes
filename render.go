package attributes

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type attributesRenderer struct{}

func (r *attributesRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(Kind, r.Render)
}

func (r *attributesRenderer) Render(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n, ok := node.(*Node)
	if !ok {
		return ast.WalkStop, fmt.Errorf("unexpected node %T, expected *Node", node)
	}

	if entering {
		if err := r.enter(w, n); err != nil {
			return ast.WalkStop, err
		}
	} else {
		r.exit(w, n)
	}

	return ast.WalkContinue, nil
}

func (r *attributesRenderer) enter(w util.BufWriter, n *Node) error {
	w.WriteString(`<span`)
	html.RenderAttributes(w, n, nil)
	w.WriteString(`>`)
	return nil
}

func (r *attributesRenderer) exit(w util.BufWriter, n *Node) {
	w.WriteString("</span>")
}
