package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Leixb/jutge/database"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req"

	"regexp"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type newfile struct {
	Code      string
	Extension string
}

func NewNewfile() *newfile {
	return &newfile{Extension: "cpp"}
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

var reg = regexp.MustCompile(`[^a-zA-Z0-9_\-.]`)

func normalizeString(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	res, _, _ := transform.String(t, s)
	res = reg.ReplaceAllString(res, "_")
	return res
}

func GetName(code string) (string, error) {
	rq := req.New()

	var err error

	if code[0] == byte('X') {
		rq, err = getReq()
		if err != nil {
			return "", err
		}
	}

	r, err := rq.Get("https://jutge.org/problems/" + code)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromResponse(r.Response())
	if err != nil {
		return "", err
	}
	name := doc.Find(".my-trim").First().Clone().Children().Remove().End().Text()
	name = strings.TrimSpace(name)

	return name, nil
}

func (n *newfile) GetFilename() (string, error) {

	dbFile := filepath.Join(*WorkDir(), "jutge.db")
	db := database.NewJutgeDB(dbFile)
	defer db.Close()

	problemName, err := db.Query(n.Code)
	if err != nil {
		return "", err
	}

	if problemName == "" {
		problemName, err = GetName(n.Code)
		if err != nil {
			return "", err
		}
		if problemName == "" {
			return fmt.Sprintf("%s.%s", n.Code, n.Extension), nil
		}
		db.Add(n.Code, problemName)
	}

	problemName = normalizeString(problemName)

	filename := fmt.Sprintf("%s_%s.%s", n.Code, problemName, n.Extension)

	return filename, nil
}
