package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/Leixb/jutge/commands"
)

type DownloadCmd struct {
	Codes []string `arg:"" required:"" help:"the codes of problems to download"`

	Overwrite bool `help:"overwrite the existing files" default:"false"`
}

func (d *DownloadCmd) Run(globals *Globals) error {
	return commands.DownloadProblems(
		d.Codes, globals.WorkDir, globals.Concurrency, d.Overwrite, regexp.MustCompile(globals.Regex))
}

type TestCmd struct {
	Programs []string `arg:"" required:"" help:"the programs to test"`

	Code            string `help:"the code of problem to test"`
	DownloadMissing bool   `help:"download the missing programs" default:"false"`
	Overwrite       bool   `help:"overwrite the existing files" default:"false"`
}

func (t *TestCmd) Run(globals *Globals) error {
	passedTotal, countTotal, err := commands.TestPrograms(
		t.Code, t.Programs, globals.WorkDir, t.DownloadMissing, t.Overwrite, globals.Concurrency, regexp.MustCompile(globals.Regex))

	println("Passed:", passedTotal, "Total:", countTotal)

	return err
}

type UploadCmd struct {
	Files []string `arg:"" required:"" help:"the files to upload"`

	Code       string `help:"the code of problem to upload"`
	Compiler   string `help:"the compiler of problem to upload" default:"AUTO" enum:"${compilers}"`
	Annotation string `help:"the annotation of problem to upload" default:"Uploaded with jutge-go"`
	Check      bool   `help:"check the uploaded files" default:"false"`
}

func (u *UploadCmd) Run(globals *Globals) error {
	codes, err := commands.UploadFiles(u.Files, u.Code, u.Compiler, u.Annotation, globals.Concurrency, regexp.MustCompile(globals.Regex))
	if err != nil {
		return err
	}

	if u.Check {
		codeList := make([]string, len(codes))
		i := 0
		for code := range codes {
			codeList[i] = code
			i++
		}
		time.Sleep(time.Second * 10)
		err = commands.CheckProblems(codeList, globals.Concurrency, regexp.MustCompile(globals.Regex))
	}

	return err
}

type CheckCmd struct {
	Codes []string `arg:"" required:"" help:"the codes of problems to check"`
}

func (c *CheckCmd) Run(globals *Globals) error {
	return commands.CheckProblems(c.Codes, globals.Concurrency, regexp.MustCompile(globals.Regex))
}

type DatabaseCmd struct{}

func (d *DatabaseCmd) Run(globals *Globals) error {
	return nil
}

type NewCmd struct {
	Code string `arg:"" required:"" help:"the code of problem to create"`

	Extension string `help:"the extension of problem to create" default:"cc"`
}

func (n *NewCmd) Run(globals *Globals) error {
	filename, err := commands.GetFilename(globals.WorkDir, n.Code, n.Extension)
	if err != nil {
		return err
	}

	fmt.Println(filename)

	_, err = os.Create(filename)

	return err
}
