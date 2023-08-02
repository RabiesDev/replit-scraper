package datafinder

import (
	"errors"
	"github.com/RabiesDev/go-logger"
	"github.com/RabiesDev/request-helper"
	"net/http"
	"net/url"
	"strings"
)

type ProxyFinder struct {
	Logger log.Logger
	Config map[string]interface{}
}

func (finder *ProxyFinder) Find(content string) ([]interface{}, error) {
	if !finder.IsActive() {
		return nil, nil
	}

	matchedProxies := PaidProxyPattern.FindStringSubmatch(content)
	if len(matchedProxies) == 0 {
		return nil, errors.New("no matches found")
	}

	uniques := make(map[string]bool)
	var uniqueMatches []interface{}

	for i := 0; i < len(matchedProxies)/5; i++ {
		matchedProxy := matchedProxies[5*i]
		matchedProxy = formatProxy(matchedProxy)
		if !uniques[matchedProxy] && finder.IsValid(matchedProxy) {
			uniques[matchedProxy] = true
			uniqueMatches = append(uniqueMatches, string(finder.Logger.ApplyColor([]byte(matchedProxy), log.Green)))
		}
	}

	if len(uniqueMatches) == 0 {
		return nil, errors.New("no unique matches found")
	}
	return uniqueMatches, nil
}

func (finder *ProxyFinder) IsValid(match string) bool {
	if strings.Contains(strings.ToLower(match), "port") || strings.Contains(strings.ToLower(match), "example") {
		return false
	}
	parsedProxy, err := requests.ParseProxy(match)
	if err != nil || parsedProxy == nil || strings.Contains(parsedProxy.Host, "127.0.0.1") || strings.Contains(parsedProxy.Host, "localhost") {
		return false
	}

	request := requests.Get("https://www.google.com")
	response, err := requests.Do(&http.Client{
		Transport: &http.Transport{
			Proxy: func(request *http.Request) (*url.URL, error) {
				return parsedProxy, nil
			},
		},
	}, request)
	if err != nil || response.StatusCode != 200 {
		return false
	}
	return true
}

func (finder *ProxyFinder) IsActive() bool {
	return finder.Config["active"].(bool)
}

func (finder *ProxyFinder) ToString() string {
	return "Proxy"
}

func formatProxy(proxy string) string {
	proxy = strings.ReplaceAll(proxy, "\"", "")
	proxy = strings.ReplaceAll(proxy, "'", "")
	proxy = strings.ReplaceAll(proxy, "=", "")
	return proxy
}
