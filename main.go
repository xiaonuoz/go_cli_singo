package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/xiaonuoz/go_cli_singo_generate_code/generate"
)

func main() {
	app := &cli.App{
		Name:  "gogen",
		Usage: "生成model和api等通用代码",
		Commands: []*cli.Command{
			// {
			// 	Name:    "param",
			// 	Aliases: []string{"p"},
			// 	Usage:   "生成param",
			// 	Action: func(c *cli.Context) error {
			// 		generate.ProjectDir = c.Args().Get(1)
			// 		err := generate.GenerateParamCode(generate.GetStructInfoArr(c.Args().Get(0)))
			// 		if err != nil {
			// 			fmt.Println(err)
			// 		}
			// 		return nil
			// 	},
			// },
			{
				Name:    "model",
				Aliases: []string{"m"},
				Usage:   "gogen model filePath - 增加model crud方法",
				Action: func(c *cli.Context) error {
					generate.ProjectDir = c.Args().Get(1)
					err := generate.GenerateModelCode(generate.GetStructInfoArr(c.Args().Get(0)))
					if err != nil {
						fmt.Println(err)
					}
					return nil
				},
			},
			{
				Name:    "api",
				Aliases: []string{"a"},
				Usage:   "gogen api filePath - 生成api",
				Action: func(c *cli.Context) error {
					generate.ProjectDir = c.Args().Get(1)
					err := generate.GenerateApiCode(generate.GetStructInfoArr(c.Args().Get(0)))
					if err != nil {
						fmt.Println(err)
					}
					return nil
				},
			},
			{
				Name:    "handler",
				Aliases: []string{"h"},
				Usage:   "gogen handler filePath - 生成handler",
				Action: func(c *cli.Context) error {
					generate.ProjectDir = c.Args().Get(1)
					err := generate.GenerateHandlerCode(generate.GetStructInfoArr(c.Args().Get(0)))
					if err != nil {
						fmt.Println(err)
					}
					return nil
				},
			},
			// {
			// 	Name:  "sdk",
			// 	Usage: "生成sdk",
			// 	Action: func(c *cli.Context) error {
			// 		generate.ProjectDir = c.Args().Get(1)
			// 		err := generate.GenerateSDKCode(generate.GetStructInfoArr(c.Args().Get(0)))
			// 		if err != nil {
			// 			fmt.Println(err)
			// 		}
			// 		return nil
			// 	},
			// },
		},
	}

	// 启动命令行应用程序
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}
}
