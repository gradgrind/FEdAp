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
	"strings"
)

// TODO: Read all go source files seeking strings for translation.
// These strings may not be split ("xxx ..."+ "yyy ...") and must start with
// "<" and end with ">".
// These are collected into a translation file, including source file info.
// These files can be read into a map for performing the translation.

type trItem struct {
	line int
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

	// Get all go files in this directory (and subdirectories ...)
	gofiles := listFiles(abspath)

	for _, f := range gofiles {
		fmt.Printf("\n[%s]\n", f)
		data := getTrStrings(f)
		fmt.Printf("++ %s :: %s\n", data.packageName, filepath.Base(data.path))
		for _, tr := range data.items {
			fmt.Printf("    -- %04d: %s\n", tr.line, tr.text)
		}
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
				s := strings.Trim(ret.Value, "\"`")
				if strings.HasPrefix(s, "<") && strings.HasSuffix(s, ">") {
					data.items = append(data.items, trItem{
						fset.Position(ret.Pos()).Line,
						s[1 : len(s)-1],
					})
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
