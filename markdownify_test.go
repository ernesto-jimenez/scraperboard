package scraperboard

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestMarkdonwConvert(t *testing.T) {
	file, err := os.Open("testdata/markdown.html")
	fatalIfError(t, err)
	doc, err := goquery.NewDocumentFromReader(file)
	fatalIfError(t, err)

	selection := doc.Find("#content")
	actual := strings.TrimSpace(markdownify(selection))

	md, err := ioutil.ReadFile("testdata/markdown.md")
	fatalIfError(t, err)
	expected := strings.TrimSpace(string(md))

	if actual != expected {
		t.Fatalf("Expected: %#v\n\tGot: %#v", expected, actual)
	}
}
