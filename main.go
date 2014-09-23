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

func main() {
	flag.StringVar(&outDir, "o", "./output", "Output directory to save crawled data")
	flag.Parse()
	config := GetConfig(*configFile)
	prepareOutputDir()
	code, ok := unicode.Scripts[config.Unicode]
	if !ok {
		panic("Unable to find unicode with name " + config.Unicode)
	}
	if !*noCrawl {
		fmt.Printf("No of sites to crawl : %d\n", len(config.Sites))
		var wg sync.WaitGroup

		for _, siteConfig := range config.Sites {
			wg.Add(1)
			go crawlSite(siteConfig, &wg)
		}
		wg.Wait()
	}
	genUnicodeWordFiles(outDir, code)
}

func prepareOutputDir() {
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		fmt.Println("Unable to create output dir: ", err)
		os.Exit(1)
	}
}
