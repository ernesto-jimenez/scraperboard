package scraperboard

import (
	"io"

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

// ExtractFromReader scrapes the HTML served in the specified URL into a golang struct
func (s *Scraper) ExtractFromReader(r io.Reader, target interface{}) (err error) {
	res, err := s.ScrapeFromReader(r)
	if err != nil {
		return
	}

	return mapstructure.Decode(res, target)
}
