package store

import (
	"fmt"
	"regexp"
)

// TODO: Complete regex
var (
	INDEX_RGX = "^index-"
)

type Index string

func CreateIndex(name string) (Index, error) {
	match, _ := regexp.MatchString(INDEX_RGX, name)

	if !match {
		return "", fmt.Errorf("invalid index name: %s", name)
	}

	return Index(name), nil
}
