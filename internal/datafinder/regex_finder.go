package datafinder

import (
	"errors"
	"github.com/RabiesDev/go-logger"
	"regexp"
)

type RegexFinder struct {
	Logger       log.Logger
	Pattern      *regexp.Regexp
	ValidationFn func(match string) bool
	Config       map[string]interface{}
	TargetName   string
}

func (finder *RegexFinder) Find(content string) ([]interface{}, error) {
	if !finder.IsActive() {
		return nil, nil
	}

	matches := finder.Pattern.FindAllString(content, -1)
	if len(matches) == 0 {
		return nil, errors.New("no matches found")
	}

	uniques := make(map[string]bool)
	var uniqueMatches []interface{}

	for _, match := range matches {
		if !uniques[match] && finder.IsValid(match) {
			uniques[match] = true
			uniqueMatches = append(uniqueMatches, string(finder.Logger.ApplyColor([]byte(match), log.Green)))
		}
	}

	if len(uniqueMatches) == 0 {
		return nil, errors.New("no unique matches found")
	}
	return uniqueMatches, nil
}

func (finder *RegexFinder) IsValid(match string) bool {
	// Use the custom validation function if set, otherwise default to true
	if finder.ValidationFn != nil {
		return finder.ValidationFn(match)
	}
	return true
}

func (finder *RegexFinder) IsActive() bool {
	return finder.Config["active"].(bool)
}

func (finder *RegexFinder) ToString() string {
	return finder.TargetName
}
