package helper

import "github.com/metal-stack/metal-api/cmd/metal-api/internal/tags"

func ProcessTags(ts []string) ([]string, error) {
	t := tags.New(ts)
	return t.Unique(), nil
}
