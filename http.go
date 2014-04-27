package scraperboard

import (
	"encoding/json"
	"net/http"
)

func (s *Scraper) HttpHandlerFunc(getUrl func(*http.Request) string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		res, err := s.ScrapeFromUrl(getUrl(req))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		out, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "text/json")
		w.Write(out)
	}
}
