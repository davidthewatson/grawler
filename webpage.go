package Grawler

import (
	"bytes"
	"fmt"
	"strings"
)

// Data structure representing a webpage
// It stores irs url, referencing webpages and
// its static assets
type Webpage struct {
	url             string
	referencedPages *[]Webpage
	staticAssets    *[]string
}

// Given a url return a Webpage object
func buildWebpageFromUrl(url string) *Webpage {
	return &Webpage{
		url,
		&[]Webpage{},
		&[]string{},
	}
}

func (page *Webpage) String() string {
	var stringBuffer bytes.Buffer
	stringBuffer.WriteString(fmt.Sprintln(page.url))
	stringBuffer.WriteString("Static assets:")
	stringBuffer.WriteString(fmt.Sprintln(page.staticAssets))

	return stringBuffer.String()
}

// Perform a dfs walk
func (page *Webpage) walk(indent string) {
	fmt.Printf("|%s%s\n", indent, page.url)
	for index, asset := range *page.staticAssets {
		fmt.Printf("|%sStatic Asset %d: %s\n", indent, index, asset)
	}
	for _, page := range *page.referencedPages {
		page.walk(strings.Repeat(indent, 2))
	}
}
