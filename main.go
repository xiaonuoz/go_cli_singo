package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/XiaoNuoZ/go_cli_singo_generate_code/generate"
	"github.com/fatih/camelcase"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Code Generator"
	app.Usage = "Generate code based on Golang struct definitions"
	app.Version = "1.0.0"

	// 获取工作路径
	// workDir := "github.com/XiaoNuoZ/go_cli_singo_generate_code"
	// workDir = filepath.Join(os.Getenv("GOPATH"), "src", workDir)
	sourceFilePath := "./test.go"
	// 读取源代码文件
	sourceFile, err := os.ReadFile(sourceFilePath)
	if err != nil {
		fmt.Println("无法读取源代码文件:", err)
		return
	}

	// 创建 Go 语法分析器
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, sourceFilePath, sourceFile, parser.ParseComments)
	if err != nil {
		fmt.Println("语法解析错误:", err)
		return
	}
	// 查找文件中的所有结构体
	var StructInfoArr []generate.StructInfo
	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						var structInfo generate.StructInfo
						structInfo.Name = ts.Name.String()
						if st, ok := ts.Type.(*ast.StructType); ok {

							for index, field := range st.Fields.List {
								fieldComment := node.Comments[index].Text()
								for _, fieldName := range field.Names {
									// 生成tag，符合驼峰命名
									entries := camelcase.Split(fieldName.Name)
									var tagNameArr []string
									for _, v := range entries {
										tagNameArr = append(tagNameArr, strings.ToLower(v))
									}
									tagName := strings.Join(tagNameArr, "_")

									structInfo.Field = append(structInfo.Field, fmt.Sprint(fieldName.Name))
									structInfo.FieldType = append(structInfo.FieldType, fmt.Sprint(field.Type))
									structInfo.Tsgs = append(structInfo.Tsgs, tagName)
									structInfo.Comments = append(structInfo.Comments, strings.ReplaceAll(fieldComment, "\n", ""))
								}
							}
						}
						StructInfoArr = append(StructInfoArr, structInfo)
					}
				}
			}
		}
	}

	if len(StructInfoArr) == 0 {
		panic("文件中不存在结构体！")
	}

	fmt.Printf("%+v", StructInfoArr)

	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"gen"},
			Usage:   "init 生成代码，不带-和--",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "param",
					Usage: "生成param，--param",
				},
				// Add more flags for customization if needed
			},
			Action: func(c *cli.Context) error {
				// structName := c.String("param")
				generate.GenerateParamCode(StructInfoArr)

				return nil
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func generateCode(structName string) string {
	// Implement the code generation logic based on the structName
	// This can include creating methods, boilerplate code, etc.

	// For demonstration purposes, let's generate a simple method.
	return fmt.Sprintf("package sdk\nfunc (s *%s) SomeMethod() {\n\t// TODO: Implement this method\n}\n", structName)
}
