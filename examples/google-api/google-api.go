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

	scraper, _ := scraperboard.NewScraperFromString(scraperXML)

	http.HandleFunc("/search", scraper.HttpHandlerFunc(getUrl))
	fmt.Println("Started API server. You can test it in http://0.0.0.0:12345/search?q=scraperboard")
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
		os.Exit(-1)
	}
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
