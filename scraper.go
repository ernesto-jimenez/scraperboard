package scraperboard

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type DebugLogger interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})
}

type debugger struct {
	logger DebugLogger
	debug  bool
}

func (d *debugger) Printf(str string, v ...interface{}) {
	if d.debug {
		d.logger.Printf(str, v...)
	}
}

func (d *debugger) Print(v ...interface{}) {
	if d.debug {
		d.logger.Print(v...)
	}
}

var debuglog = &debugger{
	logger: log.New(os.Stderr, "SCRAPER DEBUG - ", log.LstdFlags),
}

// Set boolean flag to enable logging
func Debug(debug bool) {
	debuglog.debug = debug
}

// Sets the debug logger. By default it logs to STDERR
func DefaultDebugLogger(logger DebugLogger) {
	debuglog.logger = logger
}

// NewScraperFromString constructs a Scraper based on the XML passed as a string
func NewScraperFromString(str string) (Scraper, error) {
	return NewScraper(strings.NewReader(str))
}

// NewScraperFromFile constructs a Scraper reading the XML from the file provided
func NewScraperFromFile(name string) (Scraper, error) {
	file, err := os.Open(name)
	if err != nil {
		return Scraper{}, err
	}
	return NewScraper(file)
}

// NewScraper constructs a Scraper reading the XML from the provided io.Reader
func NewScraper(r io.Reader) (scraper Scraper, err error) {
	// TODO: Validate XML: tags have required attributes, filter chain works
	err = xml.NewDecoder(r).Decode(&scraper)
	return
}

// ScrapeFromURL scrapes the provided URL and returns a map[string]interface{} that can be encoded into JSON or go structs
func (s *Scraper) ScrapeFromURL(url string) (result map[string]interface{}, err error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return
	}
	return s.scrape(doc)
}

// ScrapeFromResponse scrapes the HTML in the provided http.Response Body and returns a map[string]interface{} that can be encoded into JSON or go structs
func (s *Scraper) ScrapeFromResponse(res *http.Response) (result map[string]interface{}, err error) {
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return
	}
	return s.scrape(doc)
}

// ScrapeFromReader scrapes the HTML from the provided io.Reader and returns a map[string]interface{} that can be encoded into JSON or go structs
func (s *Scraper) ScrapeFromReader(reader io.Reader) (result map[string]interface{}, err error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return
	}
	return s.scrape(doc)
}

func (s *Scraper) scrape(doc *goquery.Document) (result map[string]interface{}, err error) {
	var sel *goquery.Selection

	if s.Selector != "" {
		sel = doc.Filter(s.Selector)
	} else {
		sel = doc.Selection
	}

	result = make(map[string]interface{})
	var k string
	var v interface{}

	for _, each := range s.EachList {
		k, v, err = each.scrape(sel)
		if err != nil {
			return
		}
		result[k] = v
	}

	for _, property := range s.PropertyList {
		k, v, err = property.scrape(sel)
		if err != nil {
			return
		}
		result[k] = v
	}

	for _, property := range s.ArrayPropertyList {
		k, v, err = property.scrape(sel)
		if err != nil {
			return
		}
		result[k] = v
	}

	if s.Name != "" {
		result = map[string]interface{}{s.Name: result}
	}
	return
}

func (s *Each) scrape(sel *goquery.Selection) (key string, value []map[string]interface{}, err error) {
	find := sel.Find(s.Selector)
	key = s.Name
	value = make([]map[string]interface{}, find.Size())

	find.Each(func(i int, sel *goquery.Selection) {
		debuglog.Printf("Processing %s/%d", s.Name, i)
		value[i] = make(map[string]interface{})

		for _, property := range s.PropertyList {
			k, v, err := property.scrape(sel)
			if err != nil {
				log.Print(err)
			} else {
				value[i][k] = v
			}
		}

		for _, property := range s.ArrayPropertyList {
			k, v, err := property.scrape(sel)
			if err != nil {
				log.Print(err)
			} else {
				value[i][k] = v
			}
		}
	})
	return
}

func (s *Property) scrape(sel *goquery.Selection) (key string, value interface{}, err error) {
	key = s.Name
	find := sel.Find(s.Selector)
	value = find
	debuglog.Printf("Property %v from %v matches", s.Name, find.Length())

	if find.Length() == 0 {
		debuglog.Print("No matches for ", s.Selector)
		value = nil
		return
	}

	if len(s.FilterList) == 0 {
		s.FilterList = defaultFilterList()
	}

	defer func() {
		if r := recover(); r != nil {
			log.Panic(r)
		}
	}()

	debuglog.Printf("Passing filters on %v", s.Name)
	for _, filter := range s.FilterList {
		value, err = filter.run(value)
		if err != nil {
			return
		}
	}
	debuglog.Printf("Property %v: %v", s.Name, value)
	return
}

func (s *ArrayProperty) scrape(sel *goquery.Selection) (key string, value interface{}, err error) {
	key = s.Name
	total := sel.Length()
	value = sel.Find(s.Selector).Map(func(i int, selection *goquery.Selection) string {
		debuglog.Printf("Scraping %s[%d/%d]", s.Name, i, total)
		var val interface{}
		val = selection
		for _, filter := range s.FilterList {
			if err != nil {
				return ""
			}
			val, err = filter.run(val)
		}
		return val.(string)
	})
	return
}

// TODO: Refactor filters using reflection to avoid type casting
func (f *Filter) run(val interface{}) (result interface{}, err error) {
	switch f.Type {
	case "first":
		result = val.(*goquery.Selection).First()
	case "last":
		result = val.(*goquery.Selection).Last()
	case "text":
		result = val.(*goquery.Selection).Text()
	case "markdown":
		result = markdownify(val.(*goquery.Selection))
	case "attr":
		result, _ = val.(*goquery.Selection).Attr(f.Argument)
	case "exists":
		count := val.(*goquery.Selection).Length()
		if count > 0 {
			result = "true"
		} else {
			result = "false"
		}
	case "queryParameter":
		var uri *url.URL
		uri, err = url.Parse(val.(string))
		result = uri.Query().Get(f.Argument)
	case "html":
		result, _ = val.(*goquery.Selection).Html()
	case "regex":
		exp := regexp.MustCompile(f.Argument)
		matches := exp.FindAllStringSubmatch(val.(string), 1)
		if matches != nil && len(matches) > 0 {
			if len(matches[0]) > 1 {
				result = matches[0][1]
			}
		}
	case "stringf":
		result = fmt.Sprintf(f.Argument, val.(string))
	case "parseDate":
		result, err = time.Parse(f.Argument, val.(string))
	default:
		err = errors.New("Unknown filter " + f.Type)
	}
	debuglog.Printf("FILTER \"%s\" (%s): %#v", f.Type, f.Argument, result)
	return
}

func defaultFilterList() []Filter {
	return []Filter{
		Filter{Type: "first"},
		Filter{Type: "text"},
	}
}

// Scraper defines a scraper template to extract structured data from HTML documents
type Scraper struct {
	Property
	EachList          []Each          `xml:"Each"`
	PropertyList      []Property      `xml:"Property"`
	ArrayPropertyList []ArrayProperty `xml:"ArrayProperty"`
}

// Each tags allow you to extract arrays of structured data (e.g: lists of reviews)
type Each struct {
	Property
	sortBy            string          `xml:"sortBy,attr"`
	PropertyList      []Property      `xml:"Property"`
	ArrayPropertyList []ArrayProperty `xml:"ArrayProperty"`
}

// Property defines a property to be extracted
type Property struct {
	Name       string   `xml:"name,attr"`
	Selector   string   `xml:"selector,attr"`
	FilterList []Filter `xml:"Filter"`
}

// ArrayProperty is used to extract array properties
type ArrayProperty struct {
	Property
	FilterList []Filter `xml:"Filter"`
}

// Filter allows you to shape the values for a property
type Filter struct {
	Type     string `xml:"type,attr"`
	Argument string `xml:"argument,attr"`
}
