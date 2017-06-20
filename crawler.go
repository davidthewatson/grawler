package Grawler

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sync"

	"golang.org/x/net/html"
)

type Crawler struct {
	urlFilter          *regexp.Regexp
	lock               *sync.Mutex
	seedPage           *Webpage
	Wg                 *sync.WaitGroup
	visitedUrls        map[string]bool
	//objectiveAssetTags map[string]bool
}

func NewCrawler(seedUrl string, tags []string, filterExpression *regexp.Regexp) *Crawler {
	crawler := new(Crawler)
	crawler.lock = new(sync.Mutex)
	crawler.urlFilter = filterExpression
	crawler.seedPage = buildWebpageFromUrl(seedUrl)
	crawler.Wg = new(sync.WaitGroup)
	crawler.visitedUrls = make(map[string]bool)
	//crawler.objectiveAssetTags = map[string]bool{}
// 	for _, tag := range tags {
// 		crawler.objectiveAssetTags[tag] = true
// 	}

	return crawler
}

func (c *Crawler) Start() {
	c.Wg.Add(1)
	c.craw(c.seedPage)
}

func (c *Crawler) fetch(URL string) (response *http.Response) {
	response, err := http.Get(URL)
	if err != nil {
		logInfo(fmt.Sprintf("Error %v fetching %v, skipping", err, URL))
		return
	}

	if response.StatusCode != 200 {
		logInfo(fmt.Sprintf("Error fetching %v, server returned status code %d", URL, response.StatusCode))
	}

	return
}

func (c *Crawler) craw(wp *Webpage) {
	defer c.Wg.Done()
	logInfo(fmt.Sprintf("Crawling %v", wp.url))

	response := c.fetch(wp.url)
	defer response.Body.Close()

	logInfo(fmt.Sprintf("Parsing content of %v", wp.url))
	c.parse(response, wp)

	logInfo(fmt.Sprintf("Finished crawling %v. Found %d href and %d static assets",
		wp.url, len(*wp.referencedPages), len(*wp.staticAssets)))
}

func (c *Crawler) parse(response *http.Response, wp *Webpage) {
	tokenizer := html.NewTokenizer(response.Body)
	for {
		tokenType := tokenizer.Next()
		switch {
		case tokenType == html.ErrorToken:
			return

		case tokenType == html.StartTagToken:
			token := tokenizer.Token()
			// Href starting tag
			if token.Data == "a" {
				c.collectAndCrawRefUrls(&token, wp)
			}
// 			} else if _, isObjectiveTag := c.objectiveAssetTags[token.Data]; isObjectiveTag {
// 				c.collectStaticAssets(&token, wp)
// 			}
		}
	}
}

// Collect and craw href of current crawling webpage matching crawler
// filter expression
func (c *Crawler) collectAndCrawRefUrls(token *html.Token, wp *Webpage) {
	for _, attr := range token.Attr {
		// get href link
		if attr.Key == "href" {
			// Rebuild relative links
			URL := c.normalizeUrl(attr.Val)

			// Check if it's an internal line
			if c.isAllowedHref(URL) {
				refPage := buildWebpageFromUrl(URL)
				*wp.referencedPages = append(*wp.referencedPages, *refPage)
				c.lock.Lock()
				_, isVisited := c.visitedUrls[URL]
				c.lock.Unlock()
				if !isVisited {
					c.lock.Lock()
					c.visitedUrls[URL] = true
					c.lock.Unlock()
					c.Wg.Add(1)
					go c.craw(refPage)
				}
			}
		}
	}
}

// Collect current crawling webpage static assets
func (c *Crawler) collectStaticAssets(token *html.Token, wp *Webpage) {
	for _, attr := range token.Attr {
		if attr.Key == "src" || attr.Key == "href" {
			*wp.staticAssets = append(*wp.staticAssets, attr.Val)
		}
	}
}

func (c *Crawler) normalizeUrl(href string) string {
	URL, _ := url.Parse(href)
	baseUrl, _ := url.Parse(c.seedPage.url)
	URL = baseUrl.ResolveReference(URL)

	return URL.String()
}

func (c *Crawler) isAllowedHref(URL string) bool {
	return c.urlFilter.MatchString(URL)
}

func (c *Crawler) Walk() {
	c.seedPage.walk("-")
}
