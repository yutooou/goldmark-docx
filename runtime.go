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
	runtime := &runTime{}
	//runtime.newDocumentWithDefaultTemplate()
	runtime.newDocument()
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

func (r *runTime) newDocument() {
	r.doc = document.New()
}
