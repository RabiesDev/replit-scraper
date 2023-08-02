package scraper

import (
	"fmt"
	"github.com/RabiesDev/go-logger"
	"github.com/RabiesDev/request-helper"
	"math"
	"math/rand"
	"net/url"
	"os"
	"regexp"
	"replit-scraper/internal/datafinder"
	"replit-scraper/internal/replit_client"
	"strings"
	"sync"
	"time"
)

type Scrapers struct {
	RepositoriesQueue chan []replit_client.Repository
	Logger            log.Logger
	Config            Config
	ScraperQueue      *CircularQueue[*Scraper]
	Finders           []datafinder.AnyTypeFinder
}

func NewScraperManager(proxies []string, sessions []string, config *Config) *Scrapers {
	scrapers := NewCircularQueue[*Scraper](len(sessions))
	for _, session := range sessions {
		var proxy *url.URL
		if len(proxies) > 0 {
			proxy, _ = requests.ParseProxy("socks5://" + proxies[rand.Intn(len(proxies))])
		}
		scrapers.Push(NewScraper(*replit_client.NewClient(proxy, session), config.Search.Query, config.Search.Sort, config.Search.Exact))
	}

	finderConfig := config.Scraper.Finders
	commonLogger := *log.NewLogger(os.Stdout, log.DefaultPrefix(), 2).WithColor()
	commonFinders := []datafinder.AnyTypeFinder{
		&datafinder.DiscordTokenFinder{
			Logger: commonLogger,
			Config: finderConfig["discord_token"].(map[string]interface{}),
		},
		&datafinder.ProxyFinder{
			Logger: commonLogger,
			Config: finderConfig["proxy"].(map[string]interface{}),
		},
		&datafinder.CaptchaServiceFinder{
			Logger: commonLogger,
			Config: finderConfig["captcha_service"].(map[string]interface{}),
		},
		&datafinder.RegexFinder{
			Logger:  commonLogger,
			Pattern: datafinder.EmailPattern,
			ValidationFn: func(match string) bool {
				domain := strings.ToLower(strings.Split(match, "@")[1])
				return !regexp.MustCompile(`\d`).MatchString(domain) && !strings.Contains(domain, "ppy") && !strings.Contains(domain, "example")
			},
			Config:     finderConfig["email"].(map[string]interface{}),
			TargetName: "Email",
		},
		&datafinder.RegexFinder{
			Logger:       commonLogger,
			Pattern:      datafinder.PasswordPattern,
			ValidationFn: nil,
			Config:       finderConfig["password"].(map[string]interface{}),
			TargetName:   "Password",
		},
		&datafinder.RegexFinder{
			Logger:       commonLogger,
			Pattern:      datafinder.PhonePattern,
			ValidationFn: nil,
			Config:       finderConfig["phone"].(map[string]interface{}),
			TargetName:   "Phone",
		},
		&datafinder.RegexFinder{
			Logger:       commonLogger,
			Pattern:      datafinder.OpenAiKeyPattern,
			ValidationFn: nil,
			Config:       finderConfig["openai_key"].(map[string]interface{}),
			TargetName:   "OpenAiKey",
		},
		&datafinder.RegexFinder{
			Logger:       commonLogger,
			Pattern:      datafinder.GoogleApiKeyPattern,
			ValidationFn: nil,
			Config:       finderConfig["google_api_key"].(map[string]interface{}),
			TargetName:   "GoogleApiKey",
		},
		&datafinder.RegexFinder{
			Logger:       commonLogger,
			Pattern:      datafinder.TelegramTokenPattern,
			ValidationFn: nil,
			Config:       finderConfig["telegram_token"].(map[string]interface{}),
			TargetName:   "Telegram token",
		},
	}

	return &Scrapers{
		RepositoriesQueue: make(chan []replit_client.Repository, len(sessions)*20),
		Logger:            commonLogger,
		Config:            *config,
		ScraperQueue:      scrapers,
		Finders:           commonFinders,
	}
}

func (scrapers *Scrapers) StartScrapers() {
	if scrapers.ScraperQueue.IsEmpty() {
		return
	}

	var group sync.WaitGroup
	runTimer := NewStopwatch()

	go scrapers.checkRepositories()
	for page := 0; page < int(math.Min(float64(scrapers.Config.Scraper.PageLimit), 20)); {
		if runTimer.Finish(time.Duration(scrapers.Config.Scraper.SearchDelay * 1000000)) {
			scraper := *scrapers.ScraperQueue.Shift()
			if scraper == nil {
				continue
			}

			group.Add(1)
			if scrapers.Config.Scraper.Parallel {
				go scraper.RunSearch(page+1, scrapers.RepositoriesQueue, &group)
			} else {
				scraper.RunSearch(page+1, scrapers.RepositoriesQueue, &group)
			}

			runTimer.Reset()
			page++
		}
	}

	group.Wait()
	close(scrapers.RepositoriesQueue)

	// After all scrapers have finished, check if there are any unprocessed repositories
	// in the RepositoriesQueue and keep calling checkRepositories until it's empty.
	for {
		if scrapers.RepositoriesQueue == nil {
			// RepositoriesQueue has been closed, all repositories have been processed.
			break
		}
	}
}

func (scrapers *Scrapers) checkRepositories() {
	processed := new(sync.Map)
	for repositories := range scrapers.RepositoriesQueue {
		for _, repository := range repositories {
			if _, ok := processed.Load(repository.Identifier); !ok {
				if scrapers.Config.Scraper.Finder {
					if scrapers.Config.Scraper.Massive {
						scrapers.Logger.Infoln(fmt.Sprintf("Reading %s(%s) directory...", scrapers.Logger.ApplyColor([]byte(repository.Title), log.Cyan), scrapers.Logger.ApplyColor([]byte(repository.UrlPath), log.Blue)))
						directories, err := repository.ReadDirectory(".")
						if err != nil {
							scrapers.Logger.Errorln(err.Error())
							continue
						}

						for _, directory := range directories {
							if directory.IsFile() {
								go func(repository replit_client.Repository, directory replit_client.Directory) {
									scrapers.Logger.Debugln(fmt.Sprintf("Validating %s...", directory.Path))
									scrapers.validateContent(repository, directory.Path, directory.Content)
								}(repository, directory)
							}
						}
					} else {
						scrapers.validateContent(repository, repository.FilePath, repository.Content)
					}
				} else {
					scrapers.Logger.Infoln(fmt.Sprintf("%s (%s)", scrapers.Logger.ApplyColor([]byte(repository.Title), log.Cyan), scrapers.Logger.ApplyColor([]byte(repository.UrlPath), log.Blue)))
				}
				processed.Store(repository.Title, true)
			}
		}
	}
}

func (scrapers *Scrapers) validateContent(repository replit_client.Repository, filePath string, content string) {
	for _, typeFinder := range scrapers.Finders {
		go func(repository replit_client.Repository, typeFinder datafinder.AnyTypeFinder, filePath string, content string) {
			if results, err := typeFinder.Find(content); err == nil && results != nil {
				for _, result := range results {
					coloredRepository := fmt.Sprintf("%s(%s)", scrapers.Logger.ApplyColor([]byte(repository.Title), log.Cyan), scrapers.Logger.ApplyColor([]byte(repository.UrlPath), log.Blue))
					scrapers.Logger.Infoln(fmt.Sprintf("%s found in %s of %s [%s]", scrapers.Logger.ApplyColor([]byte(typeFinder.ToString()), log.Bold), coloredRepository, filePath, result))
				}
			}
		}(repository, typeFinder, filePath, content)
	}
}
