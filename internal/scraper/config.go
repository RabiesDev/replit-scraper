package scraper

import (
	"bufio"
	"encoding/json"
	"os"
)

type Config struct {
	Scraper struct {
		Sessions    string                 `json:"sessions"`
		Proxies     string                 `json:"proxies"`
		PageLimit   int                    `json:"page_limit"`
		SearchDelay int                    `json:"search_delay"`
		Parallel    bool                   `json:"parallel"`
		Finder      bool                   `json:"finder"`
		Massive     bool                   `json:"massive"`
		Finders     map[string]interface{} `json:"finders"`
	} `json:"scraper"`
	Search struct {
		Query string `json:"query"`
		Sort  string `json:"sort"`
		Exact bool   `json:"exact"`
	} `json:"search"`
}

func (config *Config) ParseSessions() []string {
	sessions, err := ReadLines(config.Scraper.Sessions)
	if err != nil {
		return nil
	}
	return sessions
}

func (config *Config) ParseProxies() []string {
	proxies, err := ReadLines(config.Scraper.Proxies)
	if err != nil {
		return nil
	}
	return proxies
}

func ReadLines(path string) (lines []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return lines, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines, nil
}

func ParseConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err = json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
