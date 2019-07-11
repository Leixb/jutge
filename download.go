package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/imroc/req"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Download object that wraps its settings
type Download struct {
	codes       []string
	overwrite   bool
	concurrency int
}

// NewDownload return new Download object
func NewDownload() *Download {
	return &Download{}
}

// ConfigCommand configure kingpin options
func (d *Download) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("download", "Download problem files from jutge.org").Action(d.Run)

	// Arguments
	cmd.Arg("code", "Codes of problems to download").Required().StringsVar(&d.codes)

	// Flags
	cmd.Flag("overwrite", "Overwrite existing files").BoolVar(&d.overwrite)
	cmd.Flag("concurrency", "Number of simultaneous uploads").Default("3").IntVar(&d.concurrency)
}

// Run the command
func (d *Download) Run(c *kingpin.ParseContext) error {
	var wg sync.WaitGroup
	sem := make(chan bool, d.concurrency)

	for _, code := range d.codes {
		sem <- true
		wg.Add(1)

		// Try to regex match code or fallback to passed value
		codeReg, err := getCode(Conf.Regex, code)
		if err != nil {
			codeReg = code
		}

		go func(c string) {
			d.downloadProblem(c)
			<-sem
			wg.Done()
		}(codeReg)
	}
	wg.Wait()
	return nil
}

func (d *Download) downloadProblem(code string) error {

	r, err := req.Get("https://jutge.org/problems/" + code + "/zip")
	if err != nil {
		return err
	}

	file, err := ioutil.TempFile("", "jutge_problem_*.zip")
	defer os.Remove(file.Name())

	err = r.ToFile(file.Name())
	if err != nil {
		return err
	}

	z, err := zip.OpenReader(file.Name())
	if err != nil {
		return err
	}

	for _, f := range z.File {

		fpath := filepath.Join(Conf.WorkDir, f.Name)

		if !Conf.Quiet {
			fmt.Println("Extracting:", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()

		rc.Close()

		if err != nil {
			return err
		}

	}

	os.Remove(file.Name())

	return nil
}
