package zksyncwrapper

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generation/generate/genwrapper"
)

func WrapZksyncDeploy(zksyncBytecodePath, className, pkgName, outDirSuffixInput string) {
	fmt.Printf("Generating zk bytecode binding for %s\n", pkgName)
	outDir := genwrapper.GetOutPath(pkgName, outDirSuffixInput)
	outPath := filepath.Join(outDir, pkgName+"_zksync.go")

	fileNode := &ast.File{
		Name: ast.NewIdent(pkgName),
		Decls: []ast.Decl{
			declareImports(),
			declareDeployFunction(className),
			declareBytecodeVar(zksyncBytecodePath)}}

	println(outPath)
	writeFile(fileNode, outPath)
}

const comment = `// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.
`

func declareImports() ast.Decl {
	values := []string{
		`"context"`,
		``,
		`"github.com/ethereum/go-ethereum/accounts/abi/bind"`,
		`"github.com/ethereum/go-ethereum/common"`,
		`"github.com/ethereum/go-ethereum/ethclient"`,
		`"github.com/zksync-sdk/zksync2-go/accounts"`,
		`"github.com/zksync-sdk/zksync2-go/clients"`,
		`"github.com/zksync-sdk/zksync2-go/types"`,
	}
	specs := make([]ast.Spec, len(values))
	for i, value := range values {
		specs[i] = &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: value}}
	}

	return &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: specs}
}

func declareDeployFunction(contractName string) ast.Decl {
	data, err := os.ReadFile("../generation/zksync/zk_deploy_template.go")
	if err != nil {
		panic(err)
	}

	template := string(data)

	// remove imports, function name, first indent and closing bracket
	var (
		count = 0
		index = 0
	)
	for count < 15 { // lines to skip
		if template[index] == '\n' {
			count++
		}
		index++
	}
	template = template[index+1 : len(template)-3]

	template = strings.Replace(template, "PlaceholderContractName", contractName, 2)

	return &ast.FuncDecl{
		Name: ast.NewIdent("Deploy" + contractName + "Zk"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{{
					Names: []*ast.Ident{ast.NewIdent("auth")},
					Type:  &ast.Ident{Name: "*bind.TransactOpts"}}, {
					Names: []*ast.Ident{ast.NewIdent("ethClient")},
					Type:  &ast.Ident{Name: "*ethclient.Client"}}, {
					Names: []*ast.Ident{ast.NewIdent("wallet")},
					Type:  &ast.Ident{Name: "accounts.Wallet"}}, {
					Names: []*ast.Ident{ast.NewIdent("args")},
					Type:  &ast.Ellipsis{Elt: &ast.Ident{Name: "interface{}"}}}}},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.Ident{Name: "common.Address"}},
					{Type: &ast.Ident{Name: "*types.Receipt"}},
					{Type: &ast.StarExpr{X: &ast.Ident{Name: contractName}}},
					{Type: &ast.Ident{Name: "error"}}}}},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.BasicLit{
						Kind:  token.STRING,
						Value: template}}}}}
}

func declareBytecodeVar(srcFile string) ast.Decl {
	bytecode, err := os.ReadFile(srcFile)
	if err != nil {
		panic(err)
	}

	return &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{ast.NewIdent("ZkBytecode")},
				Values: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("common"),
							Sel: ast.NewIdent("Hex2Bytes")},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: fmt.Sprintf(`"%s"`, bytecode[2:])}}}}}}}
}

func writeFile(fileNode *ast.File, dstFile string) {
	var buf bytes.Buffer
	fset := token.NewFileSet()
	if err := format.Node(&buf, fset, fileNode); err != nil {
		panic(err)
	}

	bs := buf.Bytes()
	bs = append([]byte(comment), bs...)

	if err := os.WriteFile(dstFile, bs, 0600); err != nil {
		panic(err)
	}
}
