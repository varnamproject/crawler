package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"unicode"
	"unicode/utf8"
)

func findAllTextFiles(outputDir string) (files []string, err error) {
	files, err = filepath.Glob(filepath.Join(outputDir, "*.words"))
	return
}

func findUnicodeWords(file string, script *unicode.RangeTable, ch chan string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error while opening file %s\n", file)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		if isUnicodeWord(word, script) {
			ch <- word
		}
	}
}

func isUnicodeWord(word string, script *unicode.RangeTable) bool {
	status := true
	for len(word) > 0 {
		r, size := utf8.DecodeRuneInString(word)
		if !unicode.Is(script, r) {
			status = false
		}
		word = word[size:]
	}
	return status
}

func genUnicodeWordFiles(outputDir string, script *unicode.RangeTable, ch chan string) {
	defer close(ch)
	files, err := findAllTextFiles(outputDir)
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range files {
		findUnicodeWords(f, script, ch)
	}
}
