package datafinder

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RabiesDev/go-logger"
	"github.com/RabiesDev/request-helper"
	"net/http"
	"strconv"
	"strings"
)

type DiscordTokenFinder struct {
	Logger log.Logger
	Config map[string]interface{}
}

type GuildPreview struct {
	Identifier             string `json:"id"`
	Name                   string `json:"name"`
	ApproximateMemberCount int    `json:"approximate_member_count"`
}

func (finder *DiscordTokenFinder) Find(content string) ([]interface{}, error) {
	if !finder.IsActive() {
		return nil, nil
	}

	matchedDiscordTokens := DiscordTokenPattern.FindStringSubmatch(content)
	if len(matchedDiscordTokens) == 0 {
		return nil, errors.New("no matches found")
	}

	uniques := make(map[string]bool)
	var uniqueMatches []interface{}

	for _, token := range matchedDiscordTokens {
		if !uniques[token] && finder.IsValid(token) {
			uniques[token] = true

			botName, guilds, totalMembers := getTokenDetails(token, finder.Config["bot"].(bool))
			if botName == nil && guilds == 0 {
				finder.Logger.Errorln(fmt.Sprintf("Failed to retrieve token details (%s)", token))
				continue
			}

			coloredName := finder.Logger.ApplyColor([]byte(*botName), log.Green)
			coloredGuilds := finder.Logger.ApplyColor([]byte(strconv.Itoa(guilds)), log.Green)
			coloredTotalMember := finder.Logger.ApplyColor([]byte(strconv.Itoa(totalMembers)), log.Green)
			coloredToken := finder.Logger.ApplyColor([]byte(token), log.Green)
			uniqueMatches = append(uniqueMatches, string(finder.Logger.ApplyColor([]byte(fmt.Sprintf("Name=%s, Guilds=%s, TotalMember=%s, Token=%s", coloredName, coloredGuilds, coloredTotalMember, coloredToken)), log.Reset)))
		}
	}
	return uniqueMatches, nil
}

func (finder *DiscordTokenFinder) IsValid(match string) bool {
	return true
}

func (finder *DiscordTokenFinder) IsActive() bool {
	return finder.Config["active"].(bool)
}

func (finder *DiscordTokenFinder) ToString() string {
	if finder.Config["bot"].(bool) {
		return "Token"
	} else {
		return "Bot token"
	}
}

func getTokenDetails(token string, bot bool) (*string, int, int) {
	if !strings.Contains(token, "Bot") && bot {
		token = "Bot " + token
	}

	authorizeHeader := map[string]string{
		"authorization": fmt.Sprintf("%s", token),
	}

	request := requests.Get("https://discord.com/api/v10/users/@me")
	requests.SetHeaders(request, authorizeHeader)
	body, response, err := requests.DoAndReadByte(http.DefaultClient, request)
	if err != nil || response.StatusCode != 200 {
		return nil, 0, 0
	}

	var parsedBody map[string]interface{}
	if err = json.Unmarshal(body, &parsedBody); err != nil {
		return nil, 0, 0
	}
	botName := parsedBody["username"].(string)

	request = requests.Get("https://discord.com/api/v10/users/@me/guilds?with_counts=true")
	requests.SetHeaders(request, authorizeHeader)
	body, response, err = requests.DoAndReadByte(http.DefaultClient, request)
	if err != nil || response.StatusCode != 200 {
		return nil, 0, 0
	}

	var parsedGuilds []GuildPreview
	if err = json.Unmarshal(body, &parsedGuilds); err != nil {
		return nil, 0, 0
	}

	var totalMembers int
	for _, guild := range parsedGuilds {
		totalMembers += guild.ApproximateMemberCount
	}

	return &botName, len(parsedGuilds), totalMembers
}
