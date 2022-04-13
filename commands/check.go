package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// CheckProblems concurrently runs CheckProblem() for all problems in `codes []string`
func (j *jutge) CheckProblems(codes []string, concurrency uint) error {
	return RunParallelFuncs(codes, func(code string) error {
		code = getCodeOrSame(code, j.regex)

		veredict, err := j.CheckProblem(code)
		if err != nil {
			return err
		}

		fmt.Println(" ->", code, ":", veredict)
		return nil
	}, concurrency)
}

// CheckProblem gets veredict for the given problem code
func (j *jutge) CheckProblem(code string) (string, error) {
	rq := j.GetReq()

	r, err := rq.Get("https://jutge.org/problems/" + code)
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromResponse(r.Response())
	if err != nil {
		return "", err
	}

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
func (j *jutge) CheckSubmission(code string, submission int) (string, error) {
	rq := j.GetReq()

	r, err := rq.Get(fmt.Sprintf("https://jutge.org/problems/%s/submissions/S%03d", code, submission))
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromResponse(r.Response())
	if err != nil {
		return "", err
	}

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
func (j *jutge) GetNumSubmissions(code string) (int, error) {
	rq := j.GetReq()

	r, err := rq.Get("https://jutge.org/problems/" + code)
	if err != nil {
		return 0, err
	}
	doc, err := goquery.NewDocumentFromResponse(r.Response())
	if err != nil {
		return 0, err
	}

	submissions := doc.Find(".equal > div:nth-child(1) > div:nth-child(1) > div:nth-child(2) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(3) > dl:nth-child(1) > dd:nth-child(2)").Text()

	if submissions == "" {
		submissions = "0"
	}

	return strconv.Atoi(submissions)
}

// CheckLast gets veredict of last submission for the given code
func (j *jutge) CheckLast(code string) (string, error) {
	n, err := j.GetNumSubmissions(code)
	if err != nil {
		return "Error", err
	}

	return j.CheckSubmission(code, n)
}
