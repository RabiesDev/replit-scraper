package scraper

import (
	"fmt"
	"github.com/RabiesDev/go-logger"
	"replit-scraper/internal/replit_client"
	"strings"
	"sync"
)

type Scraper struct {
	Client replit_client.Client
	Query  string
	Sort   string
	Exact  bool
}

func NewScraper(client replit_client.Client, query, sort string, exact bool) *Scraper {
	return &Scraper{
		Client: client,
		Query:  query,
		Sort:   sort,
		Exact:  exact,
	}
}

func (scraper *Scraper) RunSearch(page int, repositoryQueue chan<- []replit_client.Repository, group *sync.WaitGroup) {
	defer group.Done()
	logger := log.Default().WithColor()
	logger.Infoln(fmt.Sprintf("Searching for page %d...", page))
	repositories, err := scraper.Client.Search(scraper.Query, page, scraper.Sort, scraper.Exact)
	if err != nil {
		if strings.EqualFold(err.Error(), "page limited") {
			logger.Warnln("Page limited!")
			return
		}
		logger.Errorln(err.Error())
		return
	} else if len(repositories) == 0 {
		logger.Warnln("No repositories found")
		return
	}
	repositoryQueue <- repositories
}
