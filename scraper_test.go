package scraperboard

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNewScraperFromFile(t *testing.T) {
	_, err := NewScraperFromFile("testdata/title.xml")
	if err != nil {
		t.Fatal(err)
	}
}

func TestErrorNewScraperFromFile(t *testing.T) {
	_, err := NewScraperFromFile("nonexistent")
	if err == nil {
		t.Fatal("Should pass on File.Open errors")
	}
}

func TestNewScraperFromString(t *testing.T) {
	_, err := NewScraperFromString("<Scraper/>")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewScraper(t *testing.T) {
	reader := strings.NewReader("<Scraper/>")
	_, err := NewScraper(reader)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewScraperFailsWithInvalidXML(t *testing.T) {
	reader := strings.NewReader("<Scraper>")
	_, err := NewScraper(reader)
	if err == nil {
		t.Fatal("Invalid XML must return an error")
	}
}

func TestScrapeProperty(t *testing.T) {
	actual, err := scraperResult("title.xml", "title_and_list_links.html")
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]interface{}{
		"title": "Page title",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected: %v\nGot: %v", expected, actual)
	}
}

func TestScrapeEach(t *testing.T) {
	actual, err := scraperResult("list_link_names.xml", "title_and_list_links.html")
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]interface{}{
		"links": []map[string]interface{}{
			map[string]interface{}{"text": "Item 1"},
			map[string]interface{}{"text": "Item 2"},
			map[string]interface{}{"text": "Item 3"},
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected: %v\nGot: %v", expected, actual)
	}
}

func TestScrapeEachAndProperties(t *testing.T) {
	actual, err := scraperResult("title_and_list_links.xml", "title_and_list_links.html")
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]interface{}{
		"title": "Page title",
		"links": []map[string]interface{}{
			map[string]interface{}{"text": "Item 1"},
			map[string]interface{}{"text": "Item 2"},
			map[string]interface{}{"text": "Item 3"},
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected: %v\nGot: %v", expected, actual)
	}
}

func TestJsonFromResult(t *testing.T) {
	assertResultFor(t, "title_and_list_links")
}

func TestSchemaOrgEvents(t *testing.T) {
	assertResultFor(t, "schema-org-events")
}

func assertResultFor(t *testing.T, name string) {
	result, err := scraperResult(name+".xml", name+".html")
	if err != nil {
		t.Fatal(err)
	}
	actual, expected, err := jsonExpectedFromFileAndActual(name+".json", result)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.EqualFold(actual, expected) {
		t.Fatalf("Expected: %#v\nGot: %#v", expected, actual)
	}
}

func jsonExpectedFromFileAndActual(expectedFile string, res map[string]interface{}) (actual, expected string, err error) {
	out, err := json.Marshal(res)
	actual = strings.TrimSpace(string(out))
	if err != nil {
		return
	}
	out, err = ioutil.ReadFile("testdata/" + expectedFile)
	expected = strings.TrimSpace(string(out))
	if err != nil {
		return
	}
	return
}

func scraperResult(scraperXML, scrapedHTML string) (res map[string]interface{}, err error) {
	scraper, err := NewScraperFromFile("testdata/" + scraperXML)
	if err != nil {
		return
	}
	html, err := os.Open("testdata/" + scrapedHTML)
	if err != nil {
		return
	}
	res, err = scraper.ScrapeFromReader(html)
	if err != nil {
		return
	}

	return
}
