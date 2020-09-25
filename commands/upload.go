package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/imroc/req"
)

type upload struct {
	Code       string
	Compiler   string
	Annotation string
	Check      bool

	codes map[string]bool
}

// NewUpload returns upload object
func NewUpload() *upload {
	return &upload{codes: make(map[string]bool), Compiler: "G++11"}
}

// GetCompilers returns a list with all valid compilers for upload
func GetCompilers() []string {
	return compilers
}

var compilers = []string{
	"AUTO", "BEEF", "Chicken", "Clang", "Clang++17", "CLISP", "Crystal",
	"Erlang", "F2C", "FBC", "FPC", "G++", "G++11", "G++17", "GCC", "GCJ",
	"GDC", "GFortran", "GHC", "GNAT", "Go", "GObjC", "GPC", "Guile", "IVL08",
	"JDK", "Lua", "MakePRO2", "MonoCS", "Nim", "nodejs", "P1++", "P2C", "Perl",
	"PHP", "PRO2", "Python", "Python3", "Quiz", "R", "Ruby", "RunHaskell",
	"RunPython", "Rust", "Stalin", "Verilog", "WS",
}

var associatedCompilers = map[string]string{
	".ada":  "GNAT",
	".bas":  "FBC",
	".bf":   "BEEF",
	".c":    "GCC",
	".cc":   "P1++",
	".cpp":  "G++11",
	".cr":   "Crystal",
	".cs":   "MonoCS",
	".d":    "GDC",
	".erl":  "Erlang",
	".f":    "GFortran",
	".go":   "Go",
	".hs":   "GHC",
	".java": "JDK",
	".js":   "nodejs",
	".lisp": "CLISP",
	".lua":  "Lua",
	".m":    "GObjC",
	".nim":  "Nim",
	".pas":  "FPC",
	".php":  "PHP",
	".pl":   "Perl",
	".py":   "Python3",
	".py2":  "Python",
	".r":    "R",
	".rb":   "Ruby",
	".scm":  "Chicken",
	".v":    "Verilog",
	".ws":   "WS",
}

// UploadFiles concurrently uploads all files in `files []string`
func (u *upload) UploadFiles(files []string) error {
	var err error

	extractCode := u.Code == ""

	var wg sync.WaitGroup
	sem := make(chan bool, conf.concurrency)

	for _, file := range files {
		sem <- true
		wg.Add(1)
		go func(f string) {
			defer func() { <-sem; wg.Done() }()

			fmt.Println(" - Uploading:", f)

			code := u.Code
			if extractCode {
				code, err = getCode(f)
				if err != nil {
					fmt.Println(" ! Can't get code for file:", f)
					return
				}

			}

			err = u.UploadFile(f, code)
			if err != nil {
				fmt.Println(" ! Upload failed", f, err)
				return
			}
			// Add code to set so it can be checked later
			u.codes[code] = true

		}(file)
	}

	wg.Wait()

	return err
}

func guessCompiler(fileName string) string {
	if compiler, ok := associatedCompilers[filepath.Ext(fileName)]; ok {
		fmt.Println(" - compiler: " + compiler)
		return compiler
	}
	fmt.Printf(" ! Warning: could not guess compiler from extension for file %s, defaulting to G++11\n", fileName)
	fmt.Println(filepath.Ext(fileName))
	panic("!!!")
	//return "G++11"
}

// UploadFile submits file to jutge.org for the problem given by code
func (u *upload) UploadFile(fileName, code string) error {
	token, err := getToken()
	if err != nil {
		return err
	}

	compiler := u.Compiler
	if compiler == "AUTO" {
		compiler = guessCompiler(fileName)
	}

	param := req.Param{
		"annotation":  u.Annotation,
		"compiler_id": compiler,
		"submit":      "submit",
		"token_uid":   token,
	}

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	file := req.FileUpload{
		File:      f,
		FieldName: "file",
		FileName:  "program",
	}

	rq, err := getReq()
	if err != nil {
		return err
	}

	_, err = rq.Post("https://jutge.org/problems/"+code+"/submissions", param, file)
	return err
}

// CheckUploaded checks the veredict of uploaded problems
//
// Bugs(): It checks the last submission for every problem code,
// meaning that if you submit more than 1 solution for the same problem,
// it will only check one of the veredicts.
func (u *upload) CheckUploaded() error {
	var wg sync.WaitGroup
	sem := make(chan bool, conf.concurrency)

	checker := NewCheck()

	for code := range u.codes {
		sem <- true
		wg.Add(1)

		go func(c string) {
			defer func() { <-sem; wg.Done() }()
			for i := 0; i < 6; i++ {
				time.Sleep(time.Second * 5)
				veredict, err := checker.CheckLast(c)
				if err != nil {
					fmt.Println(" ! Error checking", c, err)
					return
				}
				if veredict != "Not found" {
					fmt.Println(" -", c, veredict)
					return
				}
			}
			fmt.Println(" !", c, "Timed out")
		}(code)
	}

	wg.Wait()
	return nil
}
