package datafinder

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RabiesDev/go-logger"
	"github.com/RabiesDev/request-helper"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type CaptchaServiceFinder struct {
	Logger log.Logger
	Config map[string]interface{}
}

func (finder *CaptchaServiceFinder) Find(content string) ([]interface{}, error) {
	if !finder.IsActive() {
		return nil, nil
	}

	matchedServices := CaptchaServicePattern.FindStringSubmatch(content)
	if len(matchedServices) == 0 {
		return nil, errors.New("no matches found")
	}

	uniques := make(map[string]bool)
	var uniqueMatches []interface{}
	for i := 0; i < len(matchedServices)/4; i++ {
		matchedService := matchedServices[4*i+1]
		if !uniques[matchedService] && finder.IsValid(matchedService) {
			uniques[matchedService] = true

			serviceProvider, apiKey, balance := getCaptchaTokenDetails(matchedService)
			if serviceProvider == nil || apiKey == nil || balance == nil {
				return nil, errors.New("invalid captcha token")
			}

			formattedBalance := strconv.FormatFloat(*balance, 'f', 2, 64)
			coloredProvider := finder.Logger.ApplyColor([]byte(*serviceProvider), log.Green)
			coloredApiKey := finder.Logger.ApplyColor([]byte(*apiKey), log.Green)
			coloredBalance := finder.Logger.ApplyColor([]byte(formattedBalance), log.Green)
			if *balance < finder.Config["min_balance"].(float64) {
				coloredApiKey = finder.Logger.ApplyColor([]byte(*apiKey), log.Red)
				coloredBalance = finder.Logger.ApplyColor([]byte(formattedBalance), log.Red)
			}
			uniqueMatches = append(uniqueMatches, fmt.Sprintf("Provider=%s, Balance=%s, Key=%s", coloredProvider, coloredBalance, coloredApiKey))
		}
	}

	if len(uniqueMatches) == 0 {
		return nil, errors.New("no unique matches found")
	}
	return uniqueMatches, nil
}

func (finder *CaptchaServiceFinder) IsValid(match string) bool {
	if strings.Contains(match, "API_KEY") {
		return false
	}
	return true
}

func (finder *CaptchaServiceFinder) IsActive() bool {
	return finder.Config["active"].(bool)
}

func (finder *CaptchaServiceFinder) ToString() string {
	return "Captcha service"
}

func getCaptchaTokenDetails(match string) (*string, *string, *float64) {
	servicePattern := regexp.MustCompile("(.*)\\('(.*?)'\\)")
	matches := servicePattern.FindStringSubmatch(match)
	if len(matches) < 2 {
		return nil, nil, nil
	}

	provider := strings.ToLower(matches[1])
	apiKey := matches[2]
	if len(apiKey) == 0 {
		return nil, nil, nil
	}

	var balance float64
	if strings.EqualFold(provider, "twocaptcha") || strings.EqualFold(provider, "2captcha") {
		request := requests.Get(fmt.Sprintf("http://2captcha.com/res.php?key=%s&action=getbalance", apiKey))
		body, response, err := requests.DoAndReadString(&http.Client{}, request)
		if err != nil || response.StatusCode != 200 {
			return nil, nil, nil
		}

		balance, err = strconv.ParseFloat(body, 64)
		if err != nil {
			return nil, nil, nil
		}
	} else if strings.EqualFold(provider, "capmonster") {
		payload, err := json.Marshal(map[string]interface{}{
			"clientKey": apiKey,
		})
		if err != nil {
			return nil, nil, nil
		}

		request := requests.Post("https://api.capmonster.cloud/getBalance", bytes.NewReader(payload))
		requests.SetHeaders(request, map[string]string{
			"content-type": "application/json",
		})

		body, response, err := requests.DoAndReadByte(http.DefaultClient, request)
		if err != nil || response.StatusCode != 200 {
			return nil, nil, nil
		}

		var parsedBody map[string]interface{}
		if err = json.Unmarshal(body, &parsedBody); err != nil {
			return nil, nil, nil
		}

		balance = parsedBody["balance"].(float64)
	}
	return &provider, &apiKey, &balance
}
