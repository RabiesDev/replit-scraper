package replit_client

type Directory struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func (directory *Directory) IsFile() bool {
	return len(directory.Content) > 0
}
