package main

import (
	"fmt"
	"math/rand"
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
	Section        string
	outDir         string
	isSectionLinks bool
	skips          []string
}

func (this *CrawlerExtender) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {

	if doc.Find(this.Section).Length() == 0 {
		fmt.Println("Nothing in this section")
		return nil, true
	}

	section := doc.Find(this.Section)
	go func() {
		title := fmt.Sprintf("%v", rand.Int63())
		f, _ := os.Create(this.outDir + "/" + title[:] + ".txt")
		defer f.Close()
		body := section.Text()
		for _, skip := range this.skips {
			body = strings.Replace(body, skip, "", -1)
		}
		f.WriteString(body)
	}()
	if this.isSectionLinks {
		aTags := section.Find("a")
		links := make([]string, 10)
		for i := range aTags.Nodes {
			link, _ := aTags.Eq(i).Attr("href")
			links = append(links, link)
		}
		return links, false
	} else {
		return nil, true
	}
}

func crawlSite(siteConfig SiteConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	crawler := new(CrawlerExtender)
	crawler.Section = siteConfig.Section
	crawler.outDir = outDir
	crawler.skips = siteConfig.Skip
	crawler.isSectionLinks = siteConfig.IsSectionLinks
	opts := gocrawl.NewOptions(crawler)
	opts.CrawlDelay = 1 * time.Second

	opts.MaxVisits = siteConfig.Depth

	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run(siteConfig.Url)
}
