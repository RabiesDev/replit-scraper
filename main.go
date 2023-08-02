package main

import (
	"flag"
	"fmt"
	"github.com/RabiesDev/go-logger"
	"os"
	"replit-scraper/internal/scraper"
	"strconv"
	"time"
)

var Logger = log.NewLogger(os.Stdout, log.DefaultPrefix(), 1).WithColor()

func main() {
	path := flag.String("config", "config.json", "Config path")
	flag.Parse()

	config, err := scraper.ParseConfig(*path)
	if err != nil {
		Logger.Errorln(err)
		return
	}

	sessions := config.ParseSessions()
	proxies := config.ParseProxies()
	Logger.Infoln(fmt.Sprintf("%s: %s, %s: %s",
		string(Logger.ApplyColor([]byte("Sessions"), log.Bold)),
		string(Logger.ApplyColor([]byte(strconv.Itoa(len(sessions))), log.Green)),
		string(Logger.ApplyColor([]byte("Query"), log.Bold)),
		string(Logger.ApplyColor([]byte(config.Search.Query), log.Cyan)),
	))

	startTime := time.Now()
	scraperManager := scraper.NewScraperManager(proxies, sessions, config)
	scraperManager.StartScrapers()
	Logger.Debugln(fmt.Sprintf("Search completed in %d seconds", time.Since(startTime)/time.Second))
}
