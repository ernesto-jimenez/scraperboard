package scraperboard

// FIXME: Refactor this, it's quite messy

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"

	"code.google.com/p/go.net/html"
	"github.com/PuerkitoBio/goquery"
)

// MarkdownifyReader takes a io.Reader with HTML and returns the text in Markdown
func MarkdownifyReader(r io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	selection := doc.Selection
	return strings.TrimSpace(markdownify(selection)), nil
}

func markdownify(s *goquery.Selection) string {
	var buf bytes.Buffer

	// Slightly optimized vs calling Each: no single selection object created
	for _, n := range s.Nodes {
		buf.WriteString(getNodeText(n))
	}
	return strings.TrimSpace(buf.String())
}

// Get the specified node's text content.
// BUG: It doesn't respect <pre> tags
func getNodeText(node *html.Node) string {
	var buf bytes.Buffer
	// Clear redundant whitespace from text
	if node.Type == html.TextNode {
		text := normalizeWhitespace(node.Data)
		if node.NextSibling == nil || isBlock(node.NextSibling) {
			text = strings.TrimRightFunc(text, unicode.IsSpace)
		}
		if isBlock(node.NextSibling) {
			text = text + "\n\n"
		}
		if isBlock(node.PrevSibling) {
			text = strings.TrimLeftFunc(text, unicode.IsSpace)
		}
		return text
	}
	// change BRs to spaces unless it has two in which case we add extra
	if node.Data == "br" {
		if node.NextSibling.Data == "br" {
			return "\n\n"
		}
		if node.PrevSibling.Data == "br" {
			return ""
		}
		return " "
	}
	if node.FirstChild == nil {
		return ""
	}
	if node.Data == "a" {
		href, exists := getAttributeValue("href", node)
		text := getNodeText(node.FirstChild)
		if !exists {
			return text
		}
		if strings.TrimSpace(text) == "" {
			return " "
		}
		return fmt.Sprintf("[%s](%s)", text, href)
	}
	//buf.WriteString("=> " + node.Data + "|")
	if isHeader(node) {
		buf.WriteString("# ")
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		buf.WriteString(getNodeText(c))
	}
	if isBlock(node) {
		buf.WriteString("\n\n")
	}
	return buf.String()
}

func isBlock(node *html.Node) bool {
	return node != nil && (isParagraph(node) || isHeader(node))
}

func isParagraph(node *html.Node) bool {
	return node != nil && node.Data == "p"
}

func isHeader(node *html.Node) bool {
	return node != nil && len(node.Data) == 2 && node.Data[0] == 'h' && node.Data[1] != 'r'
}

// Private function to get the specified attribute's value from a node.
func getAttributeValue(attrName string, n *html.Node) (val string, exists bool) {
	if n == nil {
		return
	}

	for _, a := range n.Attr {
		if a.Key == attrName {
			val = a.Val
			exists = true
			return
		}
	}
	return
}

func normalizeWhitespace(str string) string {
	exp := regexp.MustCompile("[[:space:]]+")
	str = exp.ReplaceAllString(str, " ")
	return str
}
