package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"strings"
)

// signatures reads the ui package source and returns the exact Go signature of
// every exported component, keyed by name.
//
// Read from source rather than typed into each Section: a signature copied by
// hand drifts the first time a parameter changes, and documentation that lies
// is worse than none.
func signatures(dirs ...string) map[string]string {
	out := map[string]string{}
	cfg := &printer.Config{Mode: printer.RawFormat}
	for _, dir := range dirs {
		collect(out, cfg, dir)
	}
	return out
}

func collect(out map[string]string, cfg *printer.Config, dir string) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		log.Println("signatures:", err)
		return
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Recv != nil || !fn.Name.IsExported() {
					continue
				}
				// Print the whole FuncType: go/printer refuses a bare
				// *ast.FieldList, so parameters cannot be printed on their own.
				var b strings.Builder
				if err := cfg.Fprint(&b, fset, fn.Type); err != nil {
					log.Println("signatures:", err)
					continue
				}
				out[fn.Name.Name] = "func " + fn.Name.Name + strings.TrimPrefix(b.String(), "func")
			}
		}
	}
}

// signatureFor picks the signature for a section, matching the first component
// name in its title ("Field · SelectField" resolves to Field).
func signatureFor(sigs map[string]string, title string) string {
	for _, word := range strings.FieldsFunc(title, func(r rune) bool {
		return r == ' ' || r == '·'
	}) {
		if sig, ok := sigs[strings.TrimSpace(word)]; ok {
			return sig
		}
	}
	return ""
}
