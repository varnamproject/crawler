package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"unicode"
)

var (
	configFile = flag.String("c", "./config.json", "Configuration file for crawler")
	noCrawl    = flag.Bool("n", false, "Dont crawl, only generate words")
	outDir     string
)

func init() {
	flag.StringVar(&outDir, "o", "./output", "Output directory to save crawled data")
}

func main() {
	flag.Parse()
	config := GetConfig(*configFile)
	prepareOutputDir()
	script := getUnicodeScript(config.Script)
	ch, done := initDb()
	if !*noCrawl {
		crawlAlllSites(config)
	}
	genUnicodeWordFiles(outDir, script, ch)
	<-done
}

func crawlAlllSites(config *Config) {
	fmt.Printf("No of sites to crawl : %d\n", len(config.Sites))
	var wg sync.WaitGroup
	for _, siteConfig := range config.Sites {
		wg.Add(1)
		go crawlSite(siteConfig, &wg)
	}
	wg.Wait()
}

func getUnicodeScript(scriptName string) *unicode.RangeTable {
	script, ok := unicode.Scripts[scriptName]
	if !ok {
		panic("Unable to find unicode script with name " + scriptName)
	}
	return script
}

func prepareOutputDir() {
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		panic(err)
	}
}
