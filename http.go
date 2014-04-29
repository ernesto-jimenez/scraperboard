package scraperboard

import (
	"encoding/json"
	"net/http"
)

// NewHTTPHandlerFunc constructs an http.HandlerFunc function to expose a JSON API from the scraper. It takes a getURL function which will take each request and return the URL that should be scrapped.
// e.g.: http.Request could contain a Query Parameter with the ID to be scrapped.
func (s *Scraper) NewHTTPHandlerFunc(getURL func(*http.Request) string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		res, err := s.ScrapeFromURL(getURL(req))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		out, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "text/json; charset=UTF-8")
		w.Write(out)
	}
}
