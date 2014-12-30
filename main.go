package main

import (
	"flag"
	"os"
	"sync"
	"unicode"
)

var (
	configFile string
	noCrawl    bool
	outDir     string
	keepFiles  bool
	script     string
)

func init() {
	flag.StringVar(&configFile, "c", "./config.json", "Configuration file for crawler")
	flag.StringVar(&outDir, "o", "./output", "Output directory to save crawled data")
	flag.StringVar(&script, "s", "", "Specify script, if this option provied it will ignore config flag")
	flag.BoolVar(&noCrawl, "no-crawl", false, "Dont crawl, only generate words")
	flag.BoolVar(&keepFiles, "k", false, "Keep crawled file for future")
}

func main() {
	flag.Parse()
	var config *Config
	if script == "" {
		config = GetConfig(configFile)
	} else {
		config = &Config{Script: script}
	}
	prepareOutputDir()
	clearOutputFiles()
	script := getUnicodeScript(config.Script)
	db := initDb()
	defer db.Close()
	ch, done := wordCollector(db)
	if noCrawl {
		genUnicodeWordFiles(outDir, script, ch)
	} else {
		fileChannels := crawlAlllSites(config)
		var wg sync.WaitGroup
		output := func(c <-chan string, keepFiles bool) {
			for file := range c {
				findUnicodeWords(file, script, ch)
				if !keepFiles {
					os.Remove(file)
				}
			}
			wg.Done()
		}
		wg.Add(len(fileChannels))
		for _, c := range fileChannels {
			go output(c, keepFiles)
		}
		wg.Wait()
		close(ch)
	}
	<-done
	generateVarnamFiles(db)
}

func clearOutputFiles() {
	os.Remove("words.db")
	os.Remove("output.txt")
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
