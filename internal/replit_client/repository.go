package replit_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RabiesDev/request-helper"
	"net/http"
	"net/url"
)

type Repository struct {
	SessionId  string
	Identifier string
	Title      string
	UrlPath    string
	FilePath   string
	Content    string
}

func (repository *Repository) ReadDirectory(path string) ([]Directory, error) {
	request := requests.Get(fmt.Sprintf("http://127.0.0.1:3000/directory?%s", url.Values{
		"repo":  {repository.Identifier},
		"token": {repository.SessionId},
		"path":  {path},
	}.Encode()))
	body, response, err := requests.DoAndReadByte(&http.Client{}, request)
	if err != nil {
		return nil, err
	} else if response.StatusCode != 200 {
		return nil, errors.New(string(body))
	}

	var directories []Directory
	if err = json.Unmarshal(body, &directories); err != nil {
		return nil, err
	}
	return directories, nil
}
