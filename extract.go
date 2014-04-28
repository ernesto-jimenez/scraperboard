package scraperboard

import (
	"github.com/mitchellh/mapstructure"
)

// ExtractFromURL scrapes the HTML served in the specified URL into a golang struct
func (s *Scraper) ExtractFromURL(url string, target interface{}) (err error) {
	res, err := s.ScrapeFromURL(url)
	if err != nil {
		return
	}

	return mapstructure.Decode(res, target)
}
