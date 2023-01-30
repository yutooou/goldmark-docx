package docx

import (
	"github.com/yuin/goldmark/ast"
)

type NodeRendererFunc func(r *runTime, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error)

func renderDocument(r *runTime, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func renderParagraph(r *runTime, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		para := r.doc.AddParagraph()
		r.pushPara(&para)
		run := para.AddRun()
		r.pushRun(&run)
	} else {
		r.popRun()
		r.popPara()
	}
	return ast.WalkContinue, nil
}

func renderText(r *runTime, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	segment := node.(*ast.Text).Segment
	r.writeBytes(segment.Value(source))
	return ast.WalkContinue, nil
}

func renderHeading(r *runTime, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	if entering {
		para := r.doc.AddParagraph()
		switch n.Level {
		case 1:
			para.SetStyle(heading1)
		case 2:
			para.SetStyle(heading2)
		case 3:
			para.SetStyle(heading3)
		case 4:
			para.SetStyle(heading4)
		case 5:
			para.SetStyle(heading5)
		case 6:
			para.SetStyle(heading6)
		default:
			para.SetStyle(heading1)
		}
		r.pushPara(&para)
		run := para.AddRun()
		r.pushRun(&run)
	} else {
		r.popPara()
		r.popRun()
	}
	return ast.WalkContinue, nil
}
