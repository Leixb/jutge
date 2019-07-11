package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/imroc/req"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Download object that wraps its settings
type Download struct {
	code      string
	overwrite bool
}

// NewDownload return new Download object
func NewDownload() *Download {
	return &Download{}
}

// ConfigCommand configure kingpin options
func (d *Download) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("download", "Download problem files from jutge.org").Action(d.Run)

	// Arguments
	cmd.Arg("code", "Code of problem to download").Required().StringVar(&d.code)

	// Flags
	cmd.Flag("overwrite", "Overwrite existing files").BoolVar(&d.overwrite)
}

// Run the command
func (d *Download) Run(c *kingpin.ParseContext) error {

	r, err := req.Get("https://jutge.org/problems/" + d.code + "/zip")
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

		fmt.Println("Extracting:", fpath)

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

		// Close the file without defer to close before next iteration of loop

		outFile.Close()

		rc.Close()

		if err != nil {
			return err
		}

	}

	return nil
}
