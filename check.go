package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Check settings
type Check struct {
	codes       []string
	concurrency int
}

// NewCheck return new Check object
func NewCheck() *Check {
	return &Check{concurrency: 3}
}

// ConfigCommand configure kingpin options
func (c *Check) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("check", "Check problem files from jutge.org").Action(c.Run)

	// Arguments
	cmd.Arg("code", "Codes of problems to check").Required().StringsVar(&c.codes)

	// Flags
	cmd.Flag("concurrency", "Number of simultaneous uploads").Default("3").IntVar(&c.concurrency)
}

// Run the command
func (c *Check) Run(*kingpin.ParseContext) error {
	return c.CheckProblems()
}

// CheckProblems check all problems in c.codes
func (c *Check) CheckProblems() error {
	var wg sync.WaitGroup
	sem := make(chan bool, c.concurrency)

	for _, code := range c.codes {
		sem <- true
		wg.Add(1)
		go func(pCode string) {
			defer func() { <-sem; wg.Done() }()

			pCode = getCodeOrSame(pCode)

			veredict, err := c.CheckProblem(pCode)
			if err != nil {
				fmt.Println("Error", pCode, err)
			} else {
				fmt.Printf("%s: %s\n", pCode, veredict)
			}
		}(code)

	}

	wg.Wait()

	return nil
}

// CheckProblem get problem veredict
func (c *Check) CheckProblem(code string) (string, error) {
	rq, err := Conf.getReq()
	if err != nil {
		return "", err
	}

	r, err := rq.Get("https://jutge.org/problems/" + code)
	doc, err := goquery.NewDocumentFromResponse(r.Response())

	veredict := doc.Find(".equal > div:nth-child(1) > div:nth-child(1) > div:nth-child(1)").Text()

	// Remove white space and get veredit after ":"
	splited := strings.Split(
		strings.TrimSpace(veredict), ": ")
	veredict = splited[len(splited)-1]

	if veredict == "" {
		veredict = "Not found"
	}

	return veredict, nil
}

// CheckSubmission get submission veredict
func (c *Check) CheckSubmission(code string, submission int) (string, error) {
	rq, err := Conf.getReq()
	if err != nil {
		return "", err
	}

	r, err := rq.Get(fmt.Sprintf("https://jutge.org/problems/%s/submissions/S%03d", code, submission))
	doc, err := goquery.NewDocumentFromResponse(r.Response())

	veredict := doc.Find("div.col-sm-6:nth-child(1) > div:nth-child(1) > div:nth-child(2) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > a:nth-child(1)").Text()

	// Remove white space and get veredit after ":"
	splited := strings.Split(
		strings.TrimSpace(veredict), ": ")
	veredict = splited[len(splited)-1]

	if veredict == "" {
		veredict = "Not found"
	}

	return veredict, nil
}

func (c *Check) GetNumSubmissions(code string) (int, error) {
	rq, err := Conf.getReq()
	if err != nil {
		return 0, err
	}

	r, err := rq.Get("https://jutge.org/problems/" + code)
	doc, err := goquery.NewDocumentFromResponse(r.Response())

	submissions := doc.Find(".equal > div:nth-child(1) > div:nth-child(1) > div:nth-child(2) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(3) > dl:nth-child(1) > dd:nth-child(2)").Text()

	if submissions == "" {
		submissions = "0"
	}

	return strconv.Atoi(submissions)
}

func (c *Check) CheckLast(code string) (string, error) {
	n, err := c.GetNumSubmissions(code)
	if err != nil {
		return "Error", err
	}

	return c.CheckSubmission(code, n)
}
