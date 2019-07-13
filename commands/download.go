package commands

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/imroc/req"
)

type download struct {
	Overwrite bool
}

// NewDownload return download object
func NewDownload() *download {
	return &download{Overwrite: false}
}

// DownloadProblems download all problems from d.codes
func (d *download) DownloadProblems(codes []string) error {
	var wg sync.WaitGroup
	sem := make(chan bool, conf.concurrency)

	for _, code := range codes {
		sem <- true
		wg.Add(1)

		code = getCodeOrSame(code)
		go func(c string) {
			defer func() { <-sem; wg.Done() }()

			err := d.DownloadProblem(c)
			if err != nil {
				fmt.Println("failed", c, err)
			}
		}(code)
	}
	wg.Wait()
	return nil
}

// downloadProblem download problem data to Conf.WorkDir
func (d *download) DownloadProblem(code string) error {
	rq := req.New()

	var err error

	if code[0] == byte('X') {
		rq, err = getReq()
		if err != nil {
			return err
		}
	}

	r, err := rq.Get("https://jutge.org/problems/" + code + "/zip")
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

		fpath := filepath.Join(conf.workDir, f.Name)

		if _, err = os.Stat(fpath); err == nil {
			if !d.Overwrite {
				fmt.Println("Skipping:", fpath)
				continue
			}
		}

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

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}

	}

	os.Remove(file.Name())

	return nil
}
