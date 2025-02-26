package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// TODO: Read all go source files seeking strings for translation.
// These strings may not be split ("xxx ..."+ "yyy ...") and must start with
// "<" and end with ">".
// These are collected into a translation file, including source file info.
// These files can be read into a map for performing the translation.

var tregexp *regexp.Regexp
var nlregexp *regexp.Regexp

func init() {
	tregexp = regexp.MustCompile("(?s)^[\"`]<([a-zA-Z]*)>(.+)>[\"`]$")
	//tregexp = regexp.MustCompile("(?s)^`<([a-zA-Z]*)>(.+)>`$")
	nlregexp = regexp.MustCompile(`[\t \f\r]*[\n][\t\n \f\r]*>?`)
}

type trItem struct {
	line int
	tag  string
	text string
}

type trData struct {
	path        string
	packageName string
	items       []trItem
}

func main() {
	// Get base directory
	//TODO: Another possibility might be to search upwards for go.mod?
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		if len(args) == 0 {
			log.Fatalln("ERROR* No input directory")
		}
		log.Fatalf("*ERROR* Too many command-line arguments:\n  %+v\n", args)
	}
	dirpath := args[0]
	fileInfo, err := os.Stat(dirpath)
	if err != nil {
		log.Fatalln("ERROR* " + err.Error())
	}
	abspath, err := filepath.Abs(dirpath)
	if err != nil {
		log.Fatalf("*ERROR* " + err.Error())
	}
	if !fileInfo.IsDir() {
		log.Fatalf("*ERROR* %s is not a directory", abspath)
	}
	trdir := filepath.Join(abspath, "translations")
	fileInfo, err = os.Stat(trdir)
	if err != nil {
		log.Fatalln("ERROR* " + err.Error())
	}
	if !fileInfo.IsDir() {
		log.Fatalf("*ERROR* %s is not a directory", trdir)
	}
	trfile := filepath.Join(trdir, "messages")
	ftr, err := os.Create(trfile)
	if err != nil {
		log.Fatalf("Couldn't open output file: %s\n", trfile)
	}
	defer ftr.Close()

	t := time.Now().UTC().Format("2006-01-02 15:04:05")
	writeLine(ftr, "#UTC:"+t)

	// Get all go files in this directory (and subdirectories ...)
	gofiles := listFiles(abspath)

	mmap := map[string]string{} // used in checking for duplicates
	for _, f := range gofiles {
		fmt.Printf("Reading %s\n", f)

		data := getTrStrings(f)
		writeLine(ftr,
			fmt.Sprintf("\n::%s@%s",
				data.packageName, filepath.Base(data.path)))
		for _, tr := range data.items {
			// Check this is not a duplicate
			if where, ok := mmap[tr.text]; ok {
				log.Printf("ERROR: Message is a duplicate of %s:\n"+
					"  %04d: %s\n", where, tr.line, tr.text)
				continue
			}
			mmap[tr.text] = fmt.Sprintf("%s:%04d", f, tr.line)
			writeLine(ftr,
				fmt.Sprintf("<<%04d:%s>%s", tr.line, tr.tag, tr.text))
		}
	}
}

func writeLine(f *os.File, line string) {
	_, err := f.WriteString(line + "\n")
	if err != nil {
		log.Fatalf("Couldn't write line to: %s\n  -- %s\n", f.Name(), line)
	}
}

func getTrStrings(f string) trData {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, f, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	data := trData{
		path:        f,
		packageName: node.Name.Name,
		items:       []trItem{},
	}

	ast.Inspect(node, func(n ast.Node) bool {
		ret, ok := n.(*ast.BasicLit)
		if ok {

			if ret.Kind == token.STRING {

				rm := tregexp.FindStringSubmatch(ret.Value)
				if rm != nil {
					var rmt string
					if rm[0][0] == '`' {
						// Strip and replace newlines
						rmt = nlregexp.ReplaceAllString(rm[2], "")
						// Remove remaining control characters
						rmt = strings.Map(func(r rune) rune {
							if unicode.IsPrint(r) {
								return r
							}
							return -1
						}, rmt)
						data.items = append(data.items, trItem{
							fset.Position(ret.Pos()).Line,
							rm[1],
							rmt,
						})

					} else {
						// This warns when a "suitable" string is found, but
						// enclosed in '"'.
						fmt.Printf("TODO: %d: <%s>%s>\n",
							fset.Position(ret.Pos()).Line,
							rm[1], rm[2])
					}
				}
			}
			return true
		}
		return true
	})
	return data
}

func listFiles(dir string) []string {
	var files []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && filepath.Ext(path) == ".go" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return files
}
