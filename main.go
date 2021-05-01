package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/yuin/goldmark"
	gm_ast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var updateFile string

var md goldmark.Markdown

func init() {
	flag.StringVar(&updateFile, "update", "", "file to update")
	md = goldmark.New()
}

func parseFile(f string) (gm_ast.Node, []byte, error) {
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

	// ast.Inspect(doc, func(n ast.Node) bool {
	// 	return true
	// })

	// fmt.Printf("%#v \n", doc)
	// return nil
}

type Headers []string

func extractHeaders(doc gm_ast.Node, source []byte, maxLevel int) (Headers, error) {
	var headers Headers
	gm_ast.Walk(doc, func(n gm_ast.Node, entering bool) (gm_ast.WalkStatus, error) {
		if entering {
			switch t := n.(type) {
			case *gm_ast.Heading:
				// fmt.Printf("%#v \n", source)
				fmt.Printf("%s \n", t.Text(source))
				// lines := t.Lines()
				// fmt.Printf("LINES %#v \n", lines)
				//for _, s := range lines.Sliced(0, lines.Len()-1) {
				//}
				// fmt.Printf("%#v \n", t.Lines().Len())
				// fmt.Printf("Level: %#v \n", t.Level)
				// fmt.Printf("type: %#v, kind: %#v, kindstr: %#v \n", n.Type(), n.Kind(), n.Kind().String())
			}
		}

		return gm_ast.WalkContinue, nil
	})

	return headers, nil
}

func main() {
	flag.Parse()
	fileNames := flag.Args()

	for _, f := range fileNames {
		log.Printf("Parsing file: %s", f)
		doc, src, err := parseFile(f)
		if err != nil {
			log.Printf("Error parsing file: %s", f)
			continue
		}

		extractHeaders(doc, src, 3)

	}
}
