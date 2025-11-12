package util

import (
	"fmt"
	"log"
	"strings"
)

func ListReplacer(list *[]string, old string, new string) {
	var ret []string
	for _, item := range *list {
		ret = append(ret, strings.Replace(item, old, new, -1))
	}
	*list = ret
}

func ErrorPrinter(err error) {
	if err != nil {
		log.Println(err)
		stackTrace := fmt.Sprintf("%+v", err)
		lines := strings.Split(stackTrace, "\n")
		for _, line := range lines {
			if strings.Contains(line, "urlAPI") {
				log.Println(line)
			}
		}
	}
}

var excludedFiles = []string{
	".gitignore",
	".DS_Store",
	".ini",
	".yml",
	".yaml",
	".md",
	".txt",
	".json",
	".xml",
	".csv",
	".log",
}

func LinkFilter(links []string) []string {
	var ret []string
	for _, link := range links {
		excluded := false
		for _, exclude := range excludedFiles {
			if strings.HasSuffix(link, exclude) {
				excluded = true
				break
			}
		}
		if !excluded {
			ret = append(ret, link)
		}
	}
	return ret
}
