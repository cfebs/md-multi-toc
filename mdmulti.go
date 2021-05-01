package mdmulti

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gosimple/slug"
	"github.com/yuin/goldmark"
	gm_ast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var md goldmark.Markdown

func init() {
	md = goldmark.New()

	slug.CustomSub = map[string]string{
		".": "",
	}
}

func ParseFile(f string) (gm_ast.Node, []byte, error) {
	fReader, err := os.Open(f)
	if err != nil {
		return nil, nil, err
	}
	defer fReader.Close()

	all, err := ioutil.ReadAll(fReader)
	if err != nil {
		return nil, nil, err
	}

	reader := text.NewReader(all)
	doc := md.Parser().Parse(reader)
	return doc, all, nil
}

type Header struct {
	Level int
	Text  []byte
}

func (h Header) ToListItem(filePath string) string {
	indent := strings.Repeat("    ", h.Level-1)
	var sb strings.Builder

	sb.WriteString(indent)
	sb.WriteString("* ")

	headerIdSlug := slug.Make(string(h.Text))
	sb.WriteString(fmt.Sprintf("[%s](%s#%s)", h.Text, filePath, headerIdSlug))

	return sb.String()
}

type Headers []*Header

func ExtractHeaders(doc gm_ast.Node, source []byte, maxLevel int) (Headers, error) {
	var headers Headers
	gm_ast.Walk(doc, func(n gm_ast.Node, entering bool) (gm_ast.WalkStatus, error) {
		if entering {
			switch t := n.(type) {
			case *gm_ast.Heading:
				if t.Level > maxLevel {
					log.Printf("Skipping heading %s at level %d", t.Text(source), t.Level)
					break
				}
				h := &Header{Level: t.Level, Text: t.Text(source)}
				headers = append(headers, h)
			}
		}

		return gm_ast.WalkContinue, nil
	})

	return headers, nil
}
