package docx

import (
	"baliance.com/gooxml/document"
	"goldmark-docx/utils"
)

type runTime struct {
	doc        *document.Document        // 文档操作句柄
	workdir    string                    // 工作目录 供图片加载使用
	paragraphs []*document.Paragraph     // 当前段落栈
	runs       []*document.Run           // 当前排版栈
	table      *document.Table           // 当前处理的表格
	row        *document.Row             // 当前处理的表格行
	styleCache map[string]document.Style // 当前文稿的所有Style定义，key是Style Name名
}

func newRuntime() *runTime {
	runtime := &runTime{
		styleCache: make(map[string]document.Style),
	}
	runtime.newDocumentWithDefaultTemplate()
	runtime.makeStyleCache()
	//runtime.newDocument()
	return runtime
}

func (r *runTime) pushPara(para *document.Paragraph) {
	r.paragraphs = append(r.paragraphs, para)
}

func (r *runTime) pushRun(run *document.Run) {
	r.runs = append(r.runs, run)
}

func (r *runTime) popPara() *document.Paragraph {
	if len(r.paragraphs) == 0 {
		return nil
	}
	ret := r.paragraphs[len(r.paragraphs)-1]
	r.paragraphs = r.paragraphs[:len(r.paragraphs)-1]
	return ret
}

func (r *runTime) popRun() *document.Run {
	if len(r.runs) == 0 {
		return nil
	}
	ret := r.runs[len(r.runs)-1]
	r.runs = r.runs[:len(r.runs)-1]
	return ret
}

func (r *runTime) peekRun() *document.Run {
	if len(r.runs) == 0 {
		return nil
	}
	return r.runs[len(r.runs)-1]
}

func (r *runTime) reRun() {
	if r.peekRun() != nil {
		// 如果链接之前有输出的话需要先结束掉，然后重新开一个
		r.popRun()
		para := r.peekPara()
		run := para.AddRun()
		r.pushRun(&run)
	}
}

func (r *runTime) peekPara() *document.Paragraph {
	if len(r.paragraphs) == 0 {
		return nil
	}
	return r.paragraphs[len(r.paragraphs)-1]
}

func (r *runTime) writeString(content string) {
	if len(r.runs) == 0 || len(content) == 0 {
		return
	}
	r.peekRun().AddText(content)
}

func (r *runTime) writeBytes(content []byte) {
	r.writeString(string(content))
}

func (r *runTime) newDocumentWithDefaultTemplate() error {
	temp, err := utils.CreateTemp()
	if err != nil {
		return err
	}
	defer temp.Close()
	bs, err := utils.Base64Decode(styleTemplate)
	if err != nil {
		return err
	}
	_, err = temp.Write(bs)
	if err != nil {
		return err
	}
	r.doc, err = document.OpenTemplate(temp.Name())
	return err
}

func (r *runTime) StyleNameToId(name string) string {
	if v, ok := r.styleCache[name]; ok {
		return v.StyleID()
	}
	return ""
}

func (r *runTime) makeStyleCache() {
	for _, s := range r.doc.Styles.Styles() {
		r.styleCache[s.Name()] = s
	}
}

func (r *runTime) newDocument() {
	r.doc = document.New()
}
