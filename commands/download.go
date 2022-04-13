package commands

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/imroc/req"
)

// DownloadProblems downloads all problems from `codes []string` concurrently
func (j *jutge) DownloadProblems(codes []string, overwrite bool, concurrency uint) error {
	return RunParallelFuncs(codes, func(code string) error {
		return j.DownloadProblem(getCodeOrSame(code, j.regex), j.folder, overwrite)
		}, concurrency)
}

// downloadProblem downloads problem data and stores it in folder
func (j *jutge) DownloadProblem(code, folder string, overwrite bool) error {
	rq := req.New()

	var err error

	if code[0] == byte('X') {
		rq = j.GetReq()
		if err != nil {
			return err
		}
	}

	r, err := rq.Get("https://jutge.org/problems/" + code + "/zip")
	if err != nil {
		return err
	}

	file, err := ioutil.TempFile("", "jutge_problem_*.zip")
	if err != nil {
		return err
	}
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

		fpath := filepath.Join(folder, f.Name)

		if _, err = os.Stat(fpath); err == nil {
			if !overwrite {
				fmt.Println(" + Skipping:", fpath)
				continue
			}
		}

		fmt.Println(" - Extracting:", fpath)

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
