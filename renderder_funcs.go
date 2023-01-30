package docx

import (
	"bytes"
	"github.com/yuin/goldmark/ast"
	textm "github.com/yuin/goldmark/text"
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

func renderCodeBlock(r *runTime, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(interface {
			Lines() *textm.Segments
		})

		var content []byte
		l := n.Lines().Len()
		for i := 0; i < l; i++ {
			line := n.Lines().At(i)
			content = append(content, line.Value(source)...)
		}
		// 分行渲染代码段 否则无法展示\n
		lines := bytes.Split(bytes.TrimSpace(content), []byte{'\n'})
		for _, line := range lines {
			renderCodeBlockLike(r, line)
		}
	}
	return ast.WalkContinue, nil
}

func renderCodeSpan(r *runTime, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		para := r.peekPara()
		run := para.AddRun()
		run.Properties().SetStyle(r.StyleNameToId("Code"))
		r.pushRun(&run)
	} else {
		r.popRun()
		r.reRun()
	}
	return ast.WalkContinue, nil
}

func renderCodeBlockLike(r *runTime, content []byte) {
	para := r.doc.AddParagraph()
	r.pushPara(&para)
	run := para.AddRun()
	run.Properties().SetStyle(r.StyleNameToId("CodeBlock"))
	r.pushRun(&run)
	r.writeBytes(content)
	//run.AddBreak()
	r.popRun()
	r.reRun()
	r.popPara()
}

func renderEmphasis(r *runTime, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Emphasis)
	if entering {
		run := r.peekPara().AddRun()
		r.pushRun(&run)
		switch n.Level {
		case 2:
			run.Properties().SetBold(true)
		default:
			run.Properties().SetItalic(true)
		}
	} else {
		r.popRun()
		r.reRun()
	}
	return ast.WalkContinue, nil
}

func renderHeading(r *runTime, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	if entering {
		para := r.doc.AddParagraph()
		switch n.Level {
		case 1:
			para.SetStyle(r.StyleNameToId("heading 1"))
		case 2:
			para.SetStyle(r.StyleNameToId("heading 2"))
		case 3:
			para.SetStyle(r.StyleNameToId("heading 3"))
		case 4:
			para.SetStyle(r.StyleNameToId("heading 4"))
		case 5:
			para.SetStyle(r.StyleNameToId("heading 5"))
		case 6:
			para.SetStyle(r.StyleNameToId("heading 6"))
		default:
			para.SetStyle(r.StyleNameToId("heading 1"))
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
