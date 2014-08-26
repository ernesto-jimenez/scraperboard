package scraperboard

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestMarkdownifyReader(t *testing.T) {
	file, err := os.Open("testdata/markdown.html")
	fatalIfError(t, err)

	actual, err := MarkdownifyReader(file)
	fatalIfError(t, err)

	md, err := ioutil.ReadFile("testdata/markdown_reader.md")
	fatalIfError(t, err)
	expected := strings.TrimSpace(string(md))

	if actual != expected {
		t.Fatalf("Expected: %#v\n\tGot: %#v", expected, actual)
	}
}

func TestMarkdownifyEmptyString(t *testing.T) {
	actual, err := MarkdownifyReader(strings.NewReader(""))
	fatalIfError(t, err)

	expected := ""

	if actual != expected {
		t.Fatalf("Expected: %#v\n\tGot: %#v", expected, actual)
	}
}

func TestMarkdownifyBrAsLastChild(t *testing.T) {
	str := "<span>content <br /></span>"

	actual, err := MarkdownifyReader(strings.NewReader(str))
	fatalIfError(t, err)

	expected := "content"

	if actual != expected {
		t.Fatalf("Expected: %#v\n\tGot: %#v", expected, actual)
	}
}

func TestMarkdownConvert(t *testing.T) {
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
