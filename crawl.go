package main

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

type CrawlerExtender struct {
	gocrawl.DefaultExtender
	Section string
	outDir  string
	skips   []string
}

func (this *CrawlerExtender) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	title := strings.TrimSpace(doc.Find("title").Text())

	if doc.Find(this.Section).Length() == 0 {
		return nil, true
	}

	f, _ := os.Create(this.outDir + "/" + title + ".txt")
	defer f.Close()

	section := doc.Find(this.Section)
	body := section.Text()
	for _, skip := range this.skips {
		body = strings.Replace(body, skip, "", -1)
	}
	f.WriteString(body)
	aTags := section.Find("a")
	links := make([]string, 10)
	for i := range aTags.Nodes {
		link, _ := aTags.Eq(i).Attr("href")
		links = append(links, link)
	}
	return links, false
}

func crawlSite(siteConfig SiteConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	crawler := new(CrawlerExtender)
	crawler.Section = siteConfig.Section
	crawler.outDir = *outDir
	crawler.skips = siteConfig.Skip

	opts := gocrawl.NewOptions(crawler)
	opts.CrawlDelay = 1 * time.Second

	opts.MaxVisits = siteConfig.Depth

	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run(siteConfig.Url)
}
