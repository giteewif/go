package parser

import (
	"go/ast"
	"go/token"
)

func AddImport(file *ast.File, packageName string, filename string) {
	if file.Name.Name != "main" {
		return
	}
	print("AddImportaaa ", filename, "\n")
	noImport := true
	toInsert := &ast.ImportSpec{
		Name: &ast.Ident{
			Name: "_",
		},
		Path: &ast.BasicLit{
			ValuePos: 0,
			Kind:     token.STRING,
			Value:    packageName,
		},
		EndPos: 0,
	}
	var toInsertIfStmt *ast.IfStmt
	var mainBody *ast.BlockStmt
	print("decl len", len(file.Decls))
	for _, decl := range file.Decls {
		fd, ok := decl.(*ast.GenDecl)
		print("gen", ok)
		if ok && fd.Tok == token.IMPORT {
			imports := make([]ast.Spec, 0, len(fd.Specs)+1)
			imports = append(imports, toInsert)
			imports = append(imports, fd.Specs...)
			fd.Specs = imports
			noImport = false
		}
		fc, okFunc := decl.(*ast.FuncDecl)
		print("func", okFunc)
		if okFunc {
			print("decl func")
			print(fc.Name.Name)
		}
		if okFunc && fc.Name.Name == "useless" {
			bdstmt := fc.Body
			for _, stmt := range bdstmt.List {
				ifstmt, okIf := stmt.(*ast.IfStmt)
				if okIf {
					toInsertIfStmt = ifstmt
				}
			}
		}
		if okFunc && fc.Name.Name == "main" {
			mainBody = fc.Body
		}
	}
	if mainBody != nil && toInsertIfStmt != nil {
		toInsertIfStmt.If = mainBody.Lbrace + 1
		bodyList := make([]ast.Stmt, 0, len(mainBody.List)+1)
		bodyList = append(bodyList, toInsertIfStmt)
		bodyList = append(bodyList, mainBody.List...)
		mainBody.List = bodyList
	}
	if noImport {
		decls := make([]ast.Decl, 0, len(file.Decls)+1)
		imports := make([]ast.Spec, 0, 1)
		imports = append(imports, toInsert)
		decl := &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: imports,
		}
		decls = append(decls, decl)
		decls = append(decls, file.Decls...)
		file.Decls = decls
	}
	fset := token.NewFileSet()
	ast.Print(fset, file)

	//var cfg printer.Config
	//var buf bytes.Buffer
	//
	//cfg.Fprint(&buf, fset, file)
	//fmt.Printf(buf.String())
}
