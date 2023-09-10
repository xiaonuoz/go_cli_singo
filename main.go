package main

import (
	"fmt"
	"os"

	"github.com/XiaoNuoZ/go_cli_singo_generate_code/generate"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "singo_make_api",
		Usage: "make singo project api",
		Commands: []*cli.Command{
			{
				Name:    "param",
				Aliases: []string{"p"},
				Usage:   "生成param",
				Action: func(c *cli.Context) error {
					generate.ProjectDir = c.Args().Get(1)
					generate.GenerateParamCode(generate.GetStructInfoArr(c.Args().Get(0)))
					return nil
				},
			},
			{
				Name:    "model",
				Aliases: []string{"m"},
				Usage:   "增加model crud方法",
				Action: func(c *cli.Context) error {
					generate.ProjectDir = c.Args().Get(1)
					generate.GenerateModelCode(generate.GetStructInfoArr(c.Args().Get(0)))
					return nil
				},
			},
			{
				Name:    "service",
				Aliases: []string{"s"},
				Usage:   "生成service",
				Action: func(c *cli.Context) error {
					generate.ProjectDir = c.Args().Get(1)
					generate.GenerateServiceCode(generate.GetStructInfoArr(c.Args().Get(0)))
					return nil
				},
			},
			{
				Name:    "handler",
				Aliases: []string{"h"},
				Usage:   "生成handler",
				Action: func(c *cli.Context) error {
					generate.ProjectDir = c.Args().Get(1)
					generate.GenerateHandlerCode(generate.GetStructInfoArr(c.Args().Get(0)))
					return nil
				},
			},
			{
				Name:  "sdk",
				Usage: "生成sdk",
				Action: func(c *cli.Context) error {
					generate.ProjectDir = c.Args().Get(1)
					generate.GenerateSDKCode(generate.GetStructInfoArr(c.Args().Get(0)))
					return nil
				},
			},
		},
	}

	// 启动命令行应用程序
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}
}
