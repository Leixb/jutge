package main

import (
	"fmt"
	"strings"

	"github.com/Leixb/jutge/commands"
	"github.com/alecthomas/kong"
)

type Globals struct {
	WorkDir     string         `help:"Directory to save jutge files"`
	Concurrency uint           `help:"Maximum concurrent routines"`
	Regex       string		   `help:"Regular expression used to validate and find problem codes in filenames"`
	Username    string         `help:"Username"`
	Password    string         `help:"Password"`
	Version     VersionFlag    `help:"Print version and exit"`
}

type CLI struct {
	Globals

	Download DownloadCmd `cmd:"" help:"Download problems from jutge.org"`
	Test     TestCmd     `cmd:"" help:"Test your solutions"`
	Upload   UploadCmd   `cmd:"" help:"Upload your solutions"`
	Check    CheckCmd    `cmd:"" help:"Check the status of your solutions"`
	Database DatabaseCmd `cmd:"" help:"Manage your database"`
	New      NewCmd      `cmd:"" help:"Create a new file"`
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}

func main() {
	cli := CLI{
		Globals: Globals{
			Version: VersionFlag("0.3.1"),
		},
	}

	ctx := kong.Parse(&cli,
		kong.Name("jutge"),
		kong.Description("Jutge is a command line tool to download and test problems from jutge.org"),
		kong.UsageOnError(),
		// kong.ConfigureHelp(kong.HelpOptions{
		// 	Compact: true,
		// }),
		kong.Vars{
			"version": "0.3.1",
			"compilers": strings.Join(commands.GetCompilers(), ","),
		})
	err := ctx.Run(&cli.Globals)
	ctx.FatalIfErrorf(err)
}
