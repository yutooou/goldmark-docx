package docx_test

import (
	"github.com/yuin/goldmark"
	docx "goldmark-docx"
	"os"
	"testing"
)

func Test(t *testing.T) {
	source, err := os.ReadFile("test.md")
	if err != nil {
		panic(err)
	}

	md := goldmark.New(
		goldmark.WithRenderer(docx.New()),
	)
	file, err := os.OpenFile("out.docx", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := md.Convert(source, file); err != nil {
		panic(err)
	}
}
