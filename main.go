package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
  "time"
  "strings"
  "net/http"
  "github.com/PuerkitoBio/gocrawl"
  "github.com/PuerkitoBio/goquery"
)

var (
	config = flag.String("c", "./config.json", "Configuration file for crawler")
	outDir = flag.String("o", "./output", "Output directory to save crawled data")
)

type CrawlerExtender struct {
  gocrawl.DefaultExtender
  Section string
  outDir string
}

func (this *CrawlerExtender) Visit(ctx *gocrawl.URLContext, res *http.Response, doc * goquery.Document) (interface{}, bool) {
  title := strings.TrimSpace(doc.Find("title").Text())

  if doc.Find(this.Section).Length() == 0 {
    return nil, true
  }

  f, _ := os.Create(this.outDir +"/" + title + ".txt")
  defer f.Close()

  body := doc.Find(this.Section).Text()
  f.WriteString(body)

  return nil, true
}


func main() {
	flag.Parse()
	siteConfigs := GetConfig(*config)
	prepareOutputDir()
	fmt.Printf("No of sites to crawl : %d\n", len(siteConfigs))

	var wg sync.WaitGroup
	for _, siteConfig := range siteConfigs {
    wg.Add(1)
		go crawlSite(siteConfig,&wg)
	}

	wg.Wait()
}

func prepareOutputDir() {
	if err := os.MkdirAll(*outDir, os.ModePerm); err != nil {
		fmt.Println("Unable to create output dir: ", err)
		os.Exit(1)
	}
}


func crawlSite(siteConfig SiteConfig,wg *sync.WaitGroup) {
 defer wg.Done()
 crawler := new(CrawlerExtender)
 crawler.Section = siteConfig.Section
 crawler.outDir = *outDir

 opts := gocrawl.NewOptions(crawler)
 opts.CrawlDelay = 1 * time.Second

 opts.MaxVisits = siteConfig.Depth

 // Create crawler and start at root of duckduckgo
 c := gocrawl.NewCrawlerWithOptions(opts)
 c.Run(siteConfig.Url)
}

