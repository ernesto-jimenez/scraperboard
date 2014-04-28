package main

import (
	"fmt"
	"github.com/ernesto-jimenez/scraperboard"
	"net/http"
	"net/url"
	"os"
)

func main() {
	getURL := func(req *http.Request) string {
		query := req.URL.Query().Get("q")
		fmt.Println("Searching for:", query)
		return fmt.Sprintf("https://www.google.com/search?q=%s", url.QueryEscape(query))
	}

	scraper, err := scraperboard.NewScraperFromFile("google-scraper.xml")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	http.HandleFunc("/search", scraper.NewHTTPHandlerFunc(getURL))
	fmt.Println("Started API server. You can test it in http://0.0.0.0:12345/search?q=scraperboard")
	err = http.ListenAndServe(":12345", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
		os.Exit(-1)
	}
}
