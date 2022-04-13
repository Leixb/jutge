package commands

import (
	"fmt"
	"os"
	"path/filepath"

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
func (j *jutge) UploadFiles(files []string, code, compiler, annotation string, concurrency uint) (Set, error) {
	extractCode := code == ""

	codes := make(Set)

	codeChan := make(chan string)

	go func() {
		for code := range codeChan {
			codes[code] = true
		}
	}()

	RunParallelFuncs(files, func(file string) error {
		var code string
		var err error

		if extractCode {
			code, err = j.getCode(file)
			if err != nil {
				fmt.Println(" ! Can't get code for file:", file)
				return err
			}

		}

		return j.UploadFile(file, code, compiler, annotation)

	}, concurrency)

	close(codeChan)

	return codes, nil
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
func (j *jutge) UploadFile(fileName, code, compiler, annotation string) error {
	token := j.GetToken()

	if compiler == "AUTO" || compiler == "" {
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

	rq := j.GetReq()
	if err != nil {
		return err
	}

	_, err = rq.Post("https://jutge.org/problems/"+code+"/submissions", param, file)
	return err
}
