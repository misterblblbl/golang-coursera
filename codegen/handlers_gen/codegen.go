package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := os.Create(os.Args[2])
	fmt.Fprintln(out, "package "+node.Name.Name+"\n")

	for _, f := range node.Decls {
		genDecl, ok := f.(*ast.GenDecl)
		if !ok {
			fmt.Printf("SKIP %T is not *ast.GenDecl\n", f)
			continue
		}

		for _, spec := range genDecl.Specs {
			imports, ok := spec.(*ast.ImportSpec)
			if ok {
				fmt.Fprintf(out, "import %s\n", imports.Path.Value)
			}
		}
	}
}
