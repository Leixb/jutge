package commands

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type check struct{}

// NewCheck returns check object
func NewCheck() *check {
	return &check{}
}

// CheckProblems concurrently runs CheckProblem() for all problems in `codes []string`
func (c *check) CheckProblems(codes []string) error {
	var wg sync.WaitGroup
	sem := make(chan bool, conf.concurrency)

	for _, code := range codes {
		sem <- true
		wg.Add(1)
		go func(pCode string) {
			defer func() { <-sem; wg.Done() }()

			pCode = getCodeOrSame(pCode)

			veredict, err := c.CheckProblem(pCode)
			if err != nil {
				fmt.Println(" ! Error", pCode, err)
			} else {
				fmt.Printf(" - %s: %s\n", pCode, veredict)
			}
		}(code)

	}

	wg.Wait()

	return nil
}

// CheckProblem gets veredict for the given problem code
func (c *check) CheckProblem(code string) (string, error) {
	rq, err := getReq()
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

// CheckSubmission gets submission veredict for the given code and submission number
func (c *check) CheckSubmission(code string, submission int) (string, error) {
	rq, err := getReq()
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

// GetNumSubmissions gets the number of submissions for the given code
func (c *check) GetNumSubmissions(code string) (int, error) {
	rq, err := getReq()
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

// CheckLast gets veredict of last submission for the given code
func (c *check) CheckLast(code string) (string, error) {
	n, err := c.GetNumSubmissions(code)
	if err != nil {
		return "Error", err
	}

	return c.CheckSubmission(code, n)
}
