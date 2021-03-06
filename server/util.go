package server

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func darkModeLinkFormatter(urlPath string) string {
	if strings.HasPrefix(urlPath, "/dark") {
		return strings.TrimPrefix(urlPath, "/dark")
	}
	return path.Join("/dark", urlPath)
}

func linkFormatter(link string, urlPath string) string {
	if strings.HasPrefix(urlPath, "/dark/eu") {
		return path.Join("/dark/eu", link)
	} else if strings.HasPrefix(urlPath, "/dark/us") {
		return path.Join("/dark/us", link)
	} else if strings.HasPrefix(urlPath, "/eu") {
		return path.Join("/eu", link)
	} else if strings.HasPrefix(urlPath, "/us") {
		return path.Join("/us", link)
	} else if strings.HasPrefix(urlPath, "/dark") {
		return path.Join("/dark", link)
	}
	return link
}

func findRegion(urlPath string) string {
	if strings.Contains(urlPath, "/eu/") {
		return "eu"
	} else if strings.Contains(urlPath, "/us/") {
		return "us"
	}
	return ""
}

func findDarkmode(urlPath string) bool {
	return strings.HasPrefix(urlPath, "/dark")
}

func findFileWithExtension(folder string, extension string) (string, error) {
	var files []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if strings.HasSuffix(file, extension) {
			return file, nil
		}
	}
	return "", errors.New("no file could be found")
}

func formatName(s string) string {
	if idx := strings.Index(s, "^"); idx != -1 {
		s = s[:idx]
	}
	return properTitle(s)
}

// From https://golangcookbook.com/chapters/strings/title/
func properTitle(input string) string {
	words := strings.Fields(input)
	smallwords := " a an on the to of for "

	for index, word := range words {
		if strings.Contains(smallwords, " "+word+" ") {
			words[index] = word
		} else {
			w := strings.Title(word)
			if strings.Contains(w, "'S") {
				w = strings.Replace(w, "'S", "'s", -1)
			}
			words[index] = w
		}
	}
	return strings.Join(words, " ")
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
