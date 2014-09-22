package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

var (
	config = flag.String("c", "./config.json", "Configuration file for crawler")
	outDir string
)

func main() {
	flag.StringVar(&outDir, "o", "./output", "Output directory to save crawled data")
	flag.Parse()
	siteConfigs := GetConfig(*config)
	prepareOutputDir()
	fmt.Printf("No of sites to crawl : %d\n", len(siteConfigs))

	var wg sync.WaitGroup

	// for _, siteConfig := range siteConfigs {
	// 	wg.Add(1)
	// 	go crawlSite(siteConfig, &wg)
	// }

	findAllTextFiles(outDir)
	wg.Wait()
}

func prepareOutputDir() {
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		fmt.Println("Unable to create output dir: ", err)
		os.Exit(1)
	}
}

/*
Read All files
Find all malayalam words in the file
add it to the map
aggregate the map
*/
