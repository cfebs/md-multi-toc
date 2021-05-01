package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	mdmulti "github.com/cfebs/md-multi-toc"
	"github.com/yuin/goldmark"
)

var (
	md             goldmark.Markdown
	flagUpdateFile string
	flagMaxLevel   int
)

const TocStartTag = "<!-- ts -->"
const TocEndTag = "<!-- te -->"

func init() {
	flag.StringVar(&flagUpdateFile, "update", "", "file to update")
	flag.IntVar(&flagMaxLevel, "maxlevel", 3, "max level of heading to find")
}

func main() {
	flag.Parse()
	fileNames := flag.Args()

	var toc strings.Builder

	for _, f := range fileNames {
		log.Printf("Parsing file: %s", f)
		doc, src, err := mdmulti.ParseFile(f)
		if err != nil {
			log.Printf("Error parsing file: %s", f)
			continue
		}

		headers, err := mdmulti.ExtractHeaders(doc, src, 3)
		if err != nil {
			log.Printf("Error extracting headers: %s", err)
			continue
		}

		toc.WriteString(fmt.Sprintf("## %s\n", f))
		for _, h := range headers {
			log.Printf("Found header in %s - level: %d text: %s\n", f, h.Level, h.Text)
			toc.WriteString(h.ToListItem(f))
			toc.WriteString("\n")
		}
	}

	if flagUpdateFile != "" {
		file, err := os.Open(flagUpdateFile)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		var newFile strings.Builder
		shouldAddLine := true
		sawStartTag := false
		sawEndTag := false
		tocWritten := false

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineTxt := scanner.Text()

			if tocWritten {
				newFile.WriteString(lineTxt + "\n")
				continue
			}

			if lineTxt == TocStartTag {
				sawStartTag = true
				newFile.WriteString(lineTxt + "\n")
				shouldAddLine = false
				continue
			}

			if lineTxt == TocEndTag {
				sawEndTag = true
				shouldAddLine = true
				if sawStartTag {
					newFile.WriteString(toc.String())
					tocWritten = true
				}
			}

			if shouldAddLine {
				newFile.WriteString(lineTxt + "\n")
			}
		}

		if !sawStartTag {
			log.Fatalf("Did not see a start tag %s", TocStartTag)
		}

		if !sawEndTag {
			log.Fatalf("Did not see a end tag %s", TocEndTag)
		}

		// get existing file's stat for mode
		fInfo, _ := file.Stat()
		file.Close()

		// write file back
		err = ioutil.WriteFile(flagUpdateFile, []byte(newFile.String()), fInfo.Mode())
		if err != nil {
			log.Fatalf("Err writing file %s", err)
		}

		log.Printf("Wrote TOC to %s", flagUpdateFile)
	}
}
