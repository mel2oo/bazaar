package cli

import (
	"bazaar/api"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func New() *cli.App {
	app := cli.NewApp()
	app.Name = "bazaar-cli"
	app.Usage = "恶意软件仓储命令行工具"
	app.Commands = []*cli.Command{
		{
			Name:   "upload",
			Usage:  "恶意软件上传",
			Action: Upload,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "host", Value: "127.0.0.1:18528"},
				&cli.StringFlag{Name: "file", Aliases: []string{"f"}},
				&cli.StringFlag{Name: "dir", Aliases: []string{"d"}},
				&cli.StringFlag{Name: "ext"},
				&cli.StringSliceFlag{Name: "tag"},
			},
		},
	}
	return app
}

func Upload(ctx *cli.Context) error {
	host := ctx.String("host")
	file := ctx.String("file")
	dir := ctx.String("dir")
	ext := ctx.String("ext")
	tags := ctx.StringSlice("tag")

	if len(file) > 0 {
		if err := api.MalwareUpload(host, file, ext, tags); err != nil {
			fmt.Println("upload file error", file, err)
		} else {
			fmt.Println("upload file success", file)
		}
	}

	if len(dir) > 0 {
		filepath.Walk(dir, func(path string, info fs.FileInfo, _ error) error {
			if info.IsDir() {
				return nil
			}

			if err := api.MalwareUpload(host, path, ext, tags); err != nil {
				fmt.Println("upload file error", path, err)
			} else {
				fmt.Println("upload file success", path)
			}

			return nil
		})
	}

	return nil
}
