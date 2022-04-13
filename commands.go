package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Leixb/jutge/commands"
	"github.com/Leixb/jutge/database"

	"github.com/alecthomas/kong"
)

func jutgeConf(globals *Globals) (*commands.JutgeConfig, error) {
	regex, err := regexp.Compile(globals.Regex)
	if err != nil {
		return nil, err
	}

	return &commands.JutgeConfig{
		URL:    globals.URL,
		Folder: globals.WorkDir,
		Regex:  regex,
	}, nil
}

type DownloadCmd struct {
	Codes []string `arg:"" required:"" help:"the codes of problems to download"`

	Overwrite bool `help:"overwrite the existing files" default:"false"`
}

func (d *DownloadCmd) Run(ctx *kong.Context, globals *Globals) error {
	conf, err := jutgeConf(globals)
	if err != nil {
		return err
	}

	return commands.Jutge(conf).DownloadProblems(
		d.Codes, d.Overwrite, globals.Concurrency)
}

type TestCmd struct {
	Programs []string `arg:"" required:"" type:"path" help:"the programs to test"`

	Code            string `help:"the code of problem to test"`
	DownloadMissing bool   `help:"download the missing programs" default:"false"`
	Overwrite       bool   `help:"overwrite the existing files" default:"false"`
}

func (t *TestCmd) Run(ctx *kong.Context, globals *Globals) error {
	conf, err := jutgeConf(globals)
	if err != nil {
		return err
	}

	passedTotal, countTotal, err := commands.Jutge(conf).TestPrograms(
		t.Code, t.Programs, t.DownloadMissing, t.Overwrite, globals.Concurrency)

	println("Passed:", passedTotal, "Total:", countTotal)

	return err
}

type UploadCmd struct {
	Files []string `arg:"" required:"" type:"path" help:"the files to upload"`

	Code       string `help:"the code of problem to upload"`
	Compiler   string `help:"the compiler of problem to upload" default:"AUTO" enum:"${compilers}"`
	Annotation string `help:"the annotation of problem to upload" default:"Uploaded with jutge-go"`
	Check      bool   `help:"check the uploaded files" default:"false"`
}

func (u *UploadCmd) Run(ctx *kong.Context, globals *Globals) error {
	conf, err := jutgeConf(globals)
	if err != nil {
		return err
	}

	jutge := commands.Jutge(conf)
	codes, err := jutge.UploadFiles(u.Files, u.Code, u.Compiler, u.Annotation, globals.Concurrency)
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
		err = jutge.CheckProblems(codeList, globals.Concurrency)
	}

	return err
}

type CheckCmd struct {
	Codes []string `arg:"" required:"" help:"the codes of problems to check"`

	Submission int `optional:"" short:"s" help:"the submission number to check (negative to check from back, 0 to check problem status as a whole)" default:"0"`
}

func (c *CheckCmd) Run(ctx *kong.Context, globals *Globals) error {
	conf, err := jutgeConf(globals)
	if err != nil {
		return err
	}
	jutge := commands.Jutge(conf)

	if c.Submission == 0 {
		return jutge.CheckProblems(c.Codes, globals.Concurrency)
	}

	return commands.RunParallelFuncs(c.Codes, func(code string) error {
		subn := c.Submission
		if c.Submission < 0 {
			n, err := jutge.GetNumSubmissions(code)
			if err != nil {
				fmt.Println("Error checking submissions of", code)
				return err
			}
			subn = n + 1 + c.Submission
		}
		veredict, err := jutge.CheckSubmission(code, subn)
		if err != nil {
			fmt.Println("Error checking submission", subn, "of", code)
			return err
		}
		fmt.Println(code, subn, veredict)
		return nil
	}, globals.Concurrency)
}

type DatabaseCmd struct {
	Dump struct{} `cmd:"" help:"dump the database contents"`
	Add  struct {
		Code  string `arg:"" required:"" help:"the code of problem to add"`
		Title string `arg:"" required:"" help:"the title of problem to add"`
	} `cmd:"" help:"add a problem to the database"`
	Query struct {
		Code string `arg:"" required:"" help:"the code of problem to query"`
	} `cmd:"" help:"query the database"`
	Import struct {
		ZipFile string `arg:"" required:"" type:"path" help:"the zip file to import"`
	} `cmd:"" help:"import zip file into the database"`
	Download struct{} `cmd:"" help:"download database from remote"`

	Database string `type:"path" help:"the database file" env:"JUTGE_DATABASE"`
}

func (d *DatabaseCmd) Run(ctx *kong.Context, globals *Globals) error {
	if d.Database == "" {
		d.Database = filepath.Join(globals.WorkDir, "jutge.db")
	}

	db := database.NewJutgeDB(d.Database)

	command := strings.SplitN(ctx.Command(), " ", 3)[1]

	switch command {
	case "dump":
		return db.Print()
	case "add":
		return db.Add(d.Add.Code, d.Add.Title)
	case "query":
		if title, err := db.Query(d.Query.Code); err != nil {
			fmt.Println("Code not found in database")
			return err
		} else {
			fmt.Println(title)
		}
	case "import":
		return db.ImportZip(d.Import.ZipFile)
	case "download":
		return db.Download()
	default:
		return fmt.Errorf("unknown command: %s", command)
	}

	return nil
}

type NewCmd struct {
	Code string `arg:"" required:"" help:"the code of problem to create"`

	Extension string `help:"the extension of problem to create" default:"cc"`
}

func (n *NewCmd) Run(ctx *kong.Context, globals *Globals) error {
	conf, err := jutgeConf(globals)
	if err != nil {
		return err
	}

	filename, err := commands.Jutge(conf).GetFilename(n.Code, n.Extension)
	if err != nil {
		return err
	}

	fmt.Println(filename)

	_, err = os.Create(filename)

	return err
}
