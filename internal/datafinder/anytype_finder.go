package datafinder

type AnyTypeFinder interface {
	Find(content string) ([]interface{}, error)
	IsValid(match string) bool
	IsActive() bool
	ToString() string
}
