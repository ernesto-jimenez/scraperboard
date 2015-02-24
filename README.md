**This is an early release to gather feedback. The API and XML format will probably change. I would love to know your thoughts, so [email me](mailto:erjica@gmail.com) or [send me a tweet](https://twitter.com/ernesto_jimenez)**

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

	http.HandleFunc("/search", scraper.NewHTTPHandlerFunc(getUrl))
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
 * Implement output numbers and nulls

# Acknowledgements

Making this wouldn't have been so easy without the fantastic work from
[goquery](http://github.com/PuerkitoBio/goquery)
and [mapstructure](http://github.com/mitchellh/mapstructure)

# LICENSE

Copyright (c) 2015 Ernesto Jiménez

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
