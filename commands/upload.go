package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/imroc/req"
)

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

type Set map[string]bool

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
func UploadFiles(files []string, code, compiler, annotation string, concurrency uint, regex *regexp.Regexp) (Set, error) {
	var err error

	extractCode := code == ""

	var wg sync.WaitGroup
	sem := make(chan bool, concurrency)

	codes := make(Set)

	for _, file := range files {
		sem <- true
		wg.Add(1)
		go func(f string) {
			defer func() { <-sem; wg.Done() }()

			fmt.Println(" - Uploading:", f)

			if extractCode {
				code, err = getCode(f, regex)
				if err != nil {
					fmt.Println(" ! Can't get code for file:", f)
					return
				}

			}

			err = UploadFile(f, code, compiler, annotation)
			if err != nil {
				fmt.Println(" ! Upload failed", f, err)
				return
			}
			// Add code to set so it can be checked later
			codes[code] = true

		}(file)
	}

	wg.Wait()

	return codes, err
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
func UploadFile(fileName, code, compiler, annotation string) error {
	token := getToken()

	if compiler == "AUTO" {
		compiler = guessCompiler(fileName)
	}

	param := req.Param{
		"annotation":  annotation,
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

	rq := getReq()
	if err != nil {
		return err
	}

	_, err = rq.Post("https://jutge.org/problems/"+code+"/submissions", param, file)
	return err
}
