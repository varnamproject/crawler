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
	files, err = filepath.Glob(filepath.Join(outputDir, "*.txt"))
	return
}

func findUnicodeWords(file string, code *unicode.RangeTable) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error while opening file %s\n", file)
		return
	}
	defer f.Close()
	wf, err := os.Create(file + ".words")
	if err != nil {
		fmt.Printf("Error while creating word file %s\n", file)
		return
	}
	defer wf.Close()
	scanner := bufio.NewScanner(f)
	w := bufio.NewWriter(wf)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		if isUnicodeWord(word, code) {
			fmt.Println(word)
			w.WriteString(word)
			w.WriteString("\n")
		}
	}
	w.Flush()
}

func isUnicodeWord(word string, code *unicode.RangeTable) bool {
	status := true
	for len(word) > 0 {
		r, size := utf8.DecodeRuneInString(word)
		if !unicode.Is(code, r) {
			status = false
		}
		word = word[size:]
	}
	return status
}

func genUnicodeWordFiles(outputDir string, code *unicode.RangeTable) {
	files, err := findAllTextFiles(outputDir)
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range files {
		findUnicodeWords(f, code)
	}
}
