package main

import (
	"fmt"
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
	return &Check{}
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
				fmt.Println(veredict)
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

	return strings.TrimSpace(veredict), nil
}
