Scraperboard allows you to define scrapers declaratively.

The key feautres include:
 - Easily extract structured data from HTML websites
 - Generate JSON from HTML based on the defined scraper
 - Create REST APIs to serve the scraped JSON

# How to declare a scraper

### Extract results from Google search

```xml
<Scraper>
  <Each name="results" selector="#search ol > li">
    <Property name="title" selector="h3 a"/>
    <Property name="url" selector="h3 a">
      <Filter type="first"/>
      <Filter type="attr" argument="href"/>
      <Filter type="regex" argument="q=([^&amp;]+)"/>
    </Property>
  </Each>
</Scraper>
```

# Simple API

### Creating an JSON REST API from a scraper
```go
package main

import (
	"fmt"
	"github.com/ernesto-jimenez/scraperboard"
	"net/http"
	"net/url"
	"os"
)

func main() {
	getUrl := func(req *http.Request) string {
		query := req.URL.Query().Get("q")
		fmt.Println("Searching for:", query)
		return fmt.Sprintf("https://www.google.com/search?q=%s", url.QueryEscape(query))
	}

	scraper, err := scraperboard.NewScraperFromFile("google-scraper.xml")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	http.HandleFunc("/search", scraper.HttpHandlerFunc(getUrl))
	fmt.Println("Started API server. You can test it in http://0.0.0.0:12345/search?q=scraperboard")
	err = http.ListenAndServe(":12345", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
		os.Exit(-1)
	}
}
```

### Extracting scrapped data into Go structs

```go
package main

import (
	"flag"
	"fmt"
	"github.com/ernesto-jimenez/scraperboard"
	"net/url"
	"strings"
)

func main() {
	flag.Parse()

	query := strings.Join(flag.Args(), " ")
	searchUrl := fmt.Sprintf("https://www.google.com/search?q=%s", url.QueryEscape(query))

	scraper, _ := scraperboard.NewScraperFromString(scraperXML)

	var response Response
	scraper.ExtractFromUrl(searchUrl, &response)

	for _, result := range response.Results {
		fmt.Printf("%s:\n\t%s\n", result.Title, result.Url)
	}
}

type Response struct {
	Results []Result
}

type Result struct {
	Title string
	Url   string
}

var scraperXML string = `
	<Scraper>
		<Each name="results" selector="#search ol > li">
			<Property name="title" selector="h3 a"/>
			<Property name="url" selector="h3 a">
				<Filter type="first"/>
				<Filter type="attr" argument="href"/>
				<Filter type="regex" argument="q=([^&amp;]+)"/>
			</Property>
		</Each>
	</Scraper>
`
```

# Working examples

 * [A command line tool to extract top results from a Google
   search](examples/google-cli/google-cli.go)
 * [Create a REST API to query top results from a google
   search](examples/google-api)
 * Create a REST API to return structured data from website with
   schema.org markup

# To Do

 * Implement scraping string arrays
 * Document XML document format
 * Validate XML conforms to the format
 * More documentation and examples
 * Implement support for custom filters?
 * Implement scraping numbers and nulls

# Acknowledgements

Making this wouldn't have been so easy without the fantastic work from
[goquery](http://github.com/PuerkitoBio/goquery)
and [mapstructure](http://github.com/mitchellh/mapstructure)
