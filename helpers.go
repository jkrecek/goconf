package goconf

import (
	"path/filepath"
	"os"
	"strings"
	"log"
	"regexp"
)

var (
	FULL_PATH_REGEX = regexp.MustCompile("^\\/|[A-Z]:")
)

func getExecutePath() string {
	currentPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	return currentPath + string(os.PathSeparator)
}

func getValidFullDir(relative, absolutePart string) string {
	if !strings.HasSuffix(relative, "/") {
		relative += "/"
	}

	return getValidFullPath(relative, absolutePart)
}

func getValidFullPath(relative, absolutePart string) string {
	if FULL_PATH_REGEX.MatchString(relative) {
		return relative
	} else {
		return absolutePart + relative
	}
}
