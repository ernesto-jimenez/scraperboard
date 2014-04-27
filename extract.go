package scraperboard

import (
	"github.com/mitchellh/mapstructure"
)

func (s *Scraper) ExtractFromUrl(url string, target interface{}) (err error) {
	res, err := s.ScrapeFromUrl(url)
	if err != nil {
		return
	}

	return mapstructure.Decode(res, target)
}
