package docx

import (
	"fmt"
	"github.com/yuin/goldmark/ast"
	goldrender "github.com/yuin/goldmark/renderer"
	"io"
)

func New() goldrender.Renderer {
	runtime := newRuntime()
	r := &renderer{
		nodeRendererFuncs: make(map[ast.NodeKind]NodeRendererFunc),
		runTime:           runtime,
	}
	r.RegisterNodeRendererFuncs()
	return r
}

type renderer struct {
	nodeRendererFuncs map[ast.NodeKind]NodeRendererFunc // 文档渲染函数
	runTime           *runTime                          // 运行时
	config            *Config                           // 配置
}

func (r *renderer) Render(w io.Writer, source []byte, n ast.Node) (err error) {
	err = ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkContinue
		var err error
		fmt.Println("kind: ", n.Kind())
		f := r.nodeRendererFuncs[n.Kind()]
		if f != nil {
			s, err = f(r.runTime, source, n, entering)
		}
		return s, err
	})
	if err != nil {
		return err
	}
	return r.runTime.doc.Save(w)
}

func (r *renderer) AddOptions(_ ...goldrender.Option) {
	//TODO implement me
	panic("implement me")
}

func (r *renderer) RegisterNodeRendererFuncs() {
	r.nodeRendererFuncs[ast.KindDocument] = renderDocument
	r.nodeRendererFuncs[ast.KindHeading] = renderHeading
	//r.nodeRendererFuncs[ast.KindBlockquote] = renderBlockquote
	r.nodeRendererFuncs[ast.KindFencedCodeBlock] = renderCodeBlock
	r.nodeRendererFuncs[ast.KindCodeBlock] = renderCodeBlock
	r.nodeRendererFuncs[ast.KindParagraph] = renderParagraph
	r.nodeRendererFuncs[ast.KindText] = renderText
	r.nodeRendererFuncs[ast.KindEmphasis] = renderEmphasis
	//r.nodeRendererFuncs[ast.KindLink] = renderLink
	//r.nodeRendererFuncs[ast.KindImage] = renderImage
	r.nodeRendererFuncs[ast.KindCodeSpan] = renderCodeSpan
	//r.nodeRendererFuncs[ast.KindHTMLBlock] = renderHTMLBlock
	//r.nodeRendererFuncs[ast.KindHTMLSpan] = renderHTMLSpan
	//r.nodeRendererFuncs[ast.KindList] = renderList
	//r.nodeRendererFuncs[ast.KindListItem] = renderListItem
	//r.nodeRendererFuncs[ast.KindThematicBreak] = renderThematicBreak
	//r.nodeRendererFuncs[ast.KindHardBreak] = renderHardBreak
	//r.nodeRendererFuncs[ast.KindSoftBreak] = renderSoftBreak
	//r.nodeRendererFuncs[ast.KindStrikethrough] = renderStrikethrough
	//r.nodeRendererFuncs[ast.KindTable] = renderTable
	//r.nodeRendererFuncs[ast.KindTableHead] = renderTableHead
	//r.nodeRendererFuncs[ast.KindTableRow] = renderTableRow
	//r.nodeRendererFuncs[ast.KindTableCell] = renderTableCell
	//r.nodeRendererFuncs[ast.KindFootnote] = renderFootnote
	//r.nodeRendererFuncs[ast.KindFootnoteBlock] = renderFootnoteBlock
	//r.nodeRendererFuncs[ast.KindFootnoteDefinition] = renderFootnoteDefinition
	//r.nodeRendererFuncs[ast.KindFootnoteList] = renderFootnoteList
	//r.nodeRendererFuncs[ast.KindFootnoteListItem] = renderFootnoteListItem
	//r.nodeRendererFuncs[ast.KindRawHTML] = renderRawHTML
}
