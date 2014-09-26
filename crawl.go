package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
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
	files          chan string
}

func (this *CrawlerExtender) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	if doc.Find(this.Section).Length() == 0 {
		fmt.Println("Nothing in this section")
		return nil, true
	}
	section := doc.Find(this.Section)
	title := fmt.Sprintf("%v", rand.Int63())
	body := section.Text()
	for _, skip := range this.skips {
		body = strings.Replace(body, skip, "", -1)
	}
	err := ioutil.WriteFile(this.outDir+"/"+title[:]+".txt", []byte(body), 0644)
	if err == nil {
		this.files <- this.outDir + "/" + title[:] + ".txt"
	}
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

func crawlSite(siteConfig SiteConfig) <-chan string {
	files := make(chan string, 10)
	crawler := new(CrawlerExtender)
	crawler.files = files
	crawler.Section = siteConfig.Section
	crawler.outDir = outDir
	crawler.skips = siteConfig.Skip
	crawler.isSectionLinks = siteConfig.IsSectionLinks
	opts := gocrawl.NewOptions(crawler)
	opts.CrawlDelay = 1 * time.Second

	opts.MaxVisits = siteConfig.Depth

	c := gocrawl.NewCrawlerWithOptions(opts)
	go func() {
		defer close(files)
		c.Run(siteConfig.Url)
	}()
	return files
}
