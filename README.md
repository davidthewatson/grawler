# Grawler
Grawler is a simple recursive and concurrent web crawler written in Golang. It recursively crawls starting from a given seed url and all 'child' urls matching user defined filter expression. While crawling urls, it internally builds a sitemap and collects user defined static assets.

## Installation and requirements
Grawler is tested under go1.8 and has a single dependency `golang.org/x/net/html`

Installing Grawler:

```
go get github.com/kdmu/grawler
```

## Usage example
Write down following source code and run `go run name.go > sitemap 2> trace`. This will output a seed url sitemap to `sitemap` file and a trace of crawling process.
```golang
package main

import (
	"regexp"

	"github.com/kdmu/grawler"
)

func main() {
  // Define a set of target static asset html tags
  // In our case we will collect image, css and javascript files
  targetTags := []string{"img", "link", "script"}

  // Our url filter, we will only craw internal links
  filterExpression := regexp.MustCompile("http://tomblomfield.com/.+")

  // Build a new crawler
  crawler := Grawler.NewCrawler("http://tomblomfield.com", targetTags, filterExpression)

  // Launch
  crawler.Start()

  // Make main thread wait all goroutines finish its work
  crawler.Wg.Wait()

  // Walk() method performs a dfs traverse from seed url and outputs a sitemap including assets to stdout
  crawler.Walk()
}
```

## Future Works
- [ ] JSON, XML sitemaps
- [ ] Follow robots.txt
- [ ] Apify
- [ ] Depth limit
- [ ] Logger file redirection

## License
MIT
