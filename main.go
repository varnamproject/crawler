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
		fileChannels := crawlAlllSites(config)
		var wg sync.WaitGroup
		output := func(c <-chan string) {
			for file := range c {
				findUnicodeWords(file, script, ch)
			}
			wg.Done()
		}
		wg.Add(len(fileChannels))
		for _, c := range fileChannels {
			go output(c)
		}
		wg.Wait()
		close(ch)
	}
	// genUnicodeWordFiles(outDir, script, ch)
	<-done
}

func crawlAlllSites(config *Config) []<-chan string {
	fmt.Printf("No of sites to crawl : %d\n", len(config.Sites))
	fileChannels := make([]<-chan string, len(config.Sites))
	for i, siteConfig := range config.Sites {
		fileChannels[i] = crawlSite(siteConfig)
	}
	return fileChannels
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
