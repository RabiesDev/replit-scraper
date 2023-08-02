package replit_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/RabiesDev/request-helper"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var SearchQuery = "query SearchPageSearchResults($options: SearchQueryOptions!) {\n  search(options: $options) {\n    ...SearchPageResults\n    ... on UserError {\n      message\n      __typename\n    }\n    ... on UnauthorizedError {\n      message\n      __typename\n    }\n    ... on TooManyRequestsError {\n      message\n      __typename\n    }\n    __typename\n  }\n}\n\nfragment SearchPageResults on SearchQueryResults {\n  userResults {\n    hitInfo {\n      ...HitInfo\n      __typename\n    }\n    results {\n      pageInfo {\n        ...PageInfo\n        __typename\n      }\n      items {\n        id\n        ...SearchPageResultsUser\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n  replResults {\n    hitInfo {\n      ...HitInfo\n      __typename\n    }\n    results {\n      pageInfo {\n        ...PageInfo\n        __typename\n      }\n      items {\n        id\n        ...SearchPageResultsRepl\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n  templateResults {\n    hitInfo {\n      ...HitInfo\n      __typename\n    }\n    results {\n      pageInfo {\n        ...PageInfo\n        __typename\n      }\n      items {\n        id\n        ...SearchPageResultsTemplate\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n  postResults {\n    hitInfo {\n      ...HitInfo\n      __typename\n    }\n    results {\n      pageInfo {\n        ...PageInfo\n        __typename\n      }\n      items {\n        id\n        ...SearchPageResultsPost\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n  docResults {\n    hitInfo {\n      ...HitInfo\n      __typename\n    }\n    results {\n      pageInfo {\n        ...PageInfo\n        __typename\n      }\n      items {\n        ...SearchPageResultsDoc\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n  fileResults {\n    hitInfo {\n      ...HitInfo\n      __typename\n    }\n    results {\n      pageInfo {\n        ...PageInfo\n        __typename\n      }\n      items {\n        ...SearchPageResultsFile\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n  __typename\n}\n\nfragment HitInfo on SearchQueryHitInfo {\n  totalHits\n  totalPages\n  __typename\n}\n\nfragment PageInfo on PageInfo {\n  hasPreviousPage\n  hasNextPage\n  nextCursor\n  previousCursor\n  __typename\n}\n\nfragment SearchPageResultsUser on User {\n  id\n  username\n  fullName\n  bio\n  image\n  url\n  ...UserRoles\n  __typename\n}\n\nfragment UserRoles on User {\n  roles(\n    only: [ADMIN, MODERATOR, PATRON, PYTHONISTA, DETECTIVE, LANGUAGE_JAMMER, FEATURED, REPLIT_REP, REPLIT_REP_EDU]\n  ) {\n    id\n    name\n    key\n    tagline\n    __typename\n  }\n  __typename\n}\n\nfragment SearchPageResultsRepl on Repl {\n  id\n  ...ReplPostReplCardRepl\n  __typename\n}\n\nfragment ReplPostReplCardRepl on Repl {\n  id\n  iconUrl\n  description(plainText: true)\n  ...ReplPostReplInfoRepl\n  ...ReplStatsRepl\n  ...ReplLinkRepl\n  tags {\n    id\n    ...PostsFeedNavTag\n    __typename\n  }\n  owner {\n    ... on Team {\n      id\n      username\n      url\n      image\n      __typename\n    }\n    ... on User {\n      id\n      username\n      url\n      image\n      __typename\n    }\n    __typename\n  }\n  __typename\n}\n\nfragment ReplPostReplInfoRepl on Repl {\n  id\n  title\n  description(plainText: true)\n  imageUrl\n  iconUrl\n  templateInfo {\n    label\n    iconUrl\n    __typename\n  }\n  __typename\n}\n\nfragment ReplStatsRepl on Repl {\n  id\n  likeCount\n  runCount\n  commentCount\n  __typename\n}\n\nfragment ReplLinkRepl on Repl {\n  id\n  url\n  nextPagePathname\n  __typename\n}\n\nfragment PostsFeedNavTag on Tag {\n  id\n  isOfficial\n  __typename\n}\n\nfragment SearchPageResultsTemplate on Repl {\n  id\n  ...TemplateReplCardRepl\n  __typename\n}\n\nfragment TemplateReplCardRepl on Repl {\n  id\n  iconUrl\n  templateCategory\n  title\n  description(plainText: true)\n  releasesForkCount\n  templateLabel\n  likeCount\n  url\n  owner {\n    ... on User {\n      id\n      ...TemplateReplCardFooterUser\n      __typename\n    }\n    ... on Team {\n      id\n      ...TemplateReplCardFooterTeam\n      __typename\n    }\n    __typename\n  }\n  deployment {\n    id\n    activeRelease {\n      id\n      __typename\n    }\n    __typename\n  }\n  publishedAs\n  __typename\n}\n\nfragment TemplateReplCardFooterUser on User {\n  id\n  username\n  image\n  url\n  __typename\n}\n\nfragment TemplateReplCardFooterTeam on Team {\n  id\n  username\n  image\n  url\n  __typename\n}\n\nfragment SearchPageResultsPost on Post {\n  id\n  ...ReplPostPost\n  ...ReplCardPostPost\n  ...OldPostPost\n  __typename\n}\n\nfragment ReplPostPost on Post {\n  id\n  title\n  timeCreated\n  isPinned\n  isAnnouncement\n  ...ReplViewPostActionPermissions\n  replComment {\n    id\n    body(removeMarkdown: true)\n    __typename\n  }\n  repl {\n    id\n    ...ReplViewReplActionsPermissions\n    ...ReplPostRepl\n    __typename\n  }\n  user {\n    id\n    ...ReplPostUserPostUser\n    __typename\n  }\n  recentReplComments {\n    id\n    ...ReplPostReplComment\n    __typename\n  }\n  __typename\n}\n\nfragment ReplViewPostActionPermissions on Post {\n  id\n  isHidden\n  __typename\n}\n\nfragment ReplViewReplActionsPermissions on Repl {\n  id\n  slug\n  lastPublishedAt\n  publishedAs\n  owner {\n    ... on User {\n      id\n      username\n      __typename\n    }\n    ... on Team {\n      id\n      username\n      __typename\n    }\n    __typename\n  }\n  templateReview {\n    id\n    promoted\n    __typename\n  }\n  currentUserPermissions {\n    publish\n    containerWrite\n    __typename\n  }\n  ...UnpublishReplRepl\n  ...ReplLinkRepl\n  __typename\n}\n\nfragment UnpublishReplRepl on Repl {\n  id\n  commentCount\n  likeCount\n  runCount\n  publishedAs\n  __typename\n}\n\nfragment ReplPostRepl on Repl {\n  id\n  ...ReplPostReplInfoRepl\n  ...LikeButtonRepl\n  ...ReplStatsRepl\n  tags {\n    id\n    ...PostsFeedNavTag\n    __typename\n  }\n  __typename\n}\n\nfragment LikeButtonRepl on Repl {\n  id\n  currentUserDidLike\n  likeCount\n  url\n  wasPosted\n  wasPublished\n  __typename\n}\n\nfragment ReplPostUserPostUser on User {\n  id\n  username\n  image\n  ...UserLinkUser\n  __typename\n}\n\nfragment UserLinkUser on User {\n  id\n  url\n  username\n  __typename\n}\n\nfragment ReplPostReplComment on ReplComment {\n  id\n  body\n  timeCreated\n  user {\n    id\n    ...ReplPostRecentCommentUser\n    __typename\n  }\n  __typename\n}\n\nfragment ReplPostRecentCommentUser on User {\n  id\n  username\n  image\n  ...UserLinkUser\n  __typename\n}\n\nfragment ReplCardPostPost on Post {\n  id\n  title\n  timeCreated\n  isPinned\n  isAnnouncement\n  ...ReplViewPostActionPermissions\n  repl {\n    id\n    ...ReplViewReplActionsPermissions\n    ...ReplCardPostRepl\n    __typename\n  }\n  recentReplComments {\n    id\n    ...ReplPostReplComment\n    __typename\n  }\n  user {\n    id\n    ...ReplPostUserPostUser\n    __typename\n  }\n  __typename\n}\n\nfragment ReplCardPostRepl on Repl {\n  id\n  ...LikeButtonRepl\n  ...ReplPostReplCardRepl\n  recentComments {\n    id\n    ...ReplPostReplComment\n    __typename\n  }\n  __typename\n}\n\nfragment OldPostPost on Post {\n  id\n  title\n  preview(removeMarkdown: true, length: 150)\n  url\n  commentCount\n  isPinned\n  isAnnouncement\n  timeCreated\n  ...PostLinkPost\n  user {\n    id\n    ...ReplPostUserPostUser\n    __typename\n  }\n  repl {\n    id\n    ...ReplPostRepl\n    __typename\n  }\n  board {\n    id\n    name\n    color\n    __typename\n  }\n  recentComments(count: 3) {\n    id\n    preview(removeMarkdown: true, length: 500)\n    timeCreated\n    user {\n      id\n      ...ReplPostRecentCommentUser\n      __typename\n    }\n    __typename\n  }\n  __typename\n}\n\nfragment PostLinkPost on Post {\n  id\n  url\n  __typename\n}\n\nfragment SearchPageResultsDoc on SearchResultIndexedDoc {\n  path\n  section\n  contents\n  contentMatches {\n    start\n    length\n    __typename\n  }\n  __typename\n}\n\nfragment SearchPageResultsFile on SearchResultIndexedFile {\n  repl {\n    id\n    title\n    iconUrl\n    url\n    owner {\n      ... on User {\n        id\n        image\n        username\n        __typename\n      }\n      ... on Team {\n        id\n        image\n        username\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n  filePath\n  fileContents\n  fileContentMatches {\n    start\n    length\n    __typename\n  }\n  __typename\n}"

type Client struct {
	Client    *http.Client
	SessionId string
}

func NewClient(proxy *url.URL, sessionId string) *Client {
	return &Client{
		Client: &http.Client{
			Transport: &http.Transport{
				Proxy: func(request *http.Request) (*url.URL, error) {
					return proxy, nil
				},
			},
		},
		SessionId: sessionId,
	}
}

func (client *Client) Search(query string, page int, sort string, exact bool) ([]Repository, error) {
	payload, err := json.Marshal([]interface{}{
		map[string]interface{}{
			"operationName": "SearchPageSearchResults",
			"variables": map[string]interface{}{
				"options": map[string]interface{}{
					"onlyCalculateHits": false,
					"categories":        []string{"Files"},
					"query":             query,
					"categorySettings": map[string]interface{}{
						"docs": struct{}{},
						"files": map[string]interface{}{
							"page": map[string]interface{}{
								"first": 10,
								"after": strconv.Itoa(page),
							},
							"sort":       sort,
							"exactMatch": exact,
							"myCode":     false,
						},
					},
				},
			},
			"query": SearchQuery,
		},
	})
	if err != nil {
		return nil, err
	}

	request := requests.Post("https://replit.com/graphql", bytes.NewReader(payload))
	requests.SetHeaders(request, map[string]string{
		"user-agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
		"content-type":     "application/json",
		"cookie":           "connect.sid=" + client.SessionId,
		"accept-language":  "en-US,en;q=0.5",
		"accept-encoding":  "gzip, deflate, br",
		"accept":           "*/*",
		"origin":           "https://replit.com",
		"referer":          "https://replit.com/search",
		"x-requested-with": "XMLHttpRequest",
		"sec-fetch-dest":   "empty",
		"sec-fetch-mode":   "cors",
		"sec-fetch-site":   "same-origin",
		"pragma":           "no-cache",
		"cache-control":    "no-cache",
	})
	body, _, err := requests.DoAndReadByte(client.Client, request)
	castBody := strings.ToLower(string(body))
	if err != nil {
		return nil, err
	} else if strings.Contains(castBody, "exceeded the daily quota") {
		return nil, errors.New("the daily quota for this query has been exceeded")
	} else if strings.Contains(castBody, "please try again later") {
		return nil, errors.New("rate limited has been applied")
	} else if strings.Contains(castBody, "please try a different page") {
		return nil, errors.New("page limited has been applied")
	} else if strings.Contains(castBody, "unauthorized") {
		return nil, errors.New("not authenticated")
	} else if strings.Contains(castBody, "internal server error") {
		return nil, errors.New("internal server error")
	}

	var parsedBody []interface{}
	if err = json.Unmarshal(body, &parsedBody); err != nil {
		return nil, err
	}

	var repositories []Repository
	fileResults, ok := parsedBody[0].(map[string]interface{})["data"].(map[string]interface{})["search"].(map[string]interface{})["fileResults"].(map[string]interface{})
	if !ok {
		return nil, errors.New("internal cast error")
	}

	results, ok := fileResults["results"].(map[string]interface{})
	if !ok {
		return nil, errors.New("internal cast error")
	}

	for _, matches := range results["items"].([]interface{}) {
		matchItem := matches.(map[string]interface{})
		repoInfo := matchItem["repl"].(map[string]interface{})
		repositories = append(repositories, Repository{
			SessionId:  client.SessionId,
			Identifier: repoInfo["id"].(string),
			Title:      repoInfo["title"].(string),
			UrlPath:    repoInfo["url"].(string),
			FilePath:   matchItem["filePath"].(string),
			Content:    matchItem["fileContents"].(string),
		})
	}
	return repositories, nil
}
