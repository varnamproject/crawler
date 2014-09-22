package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func findAllTextFiles(outputDir string) {
	files, err := filepath.Glob(filepath.Join(outputDir, "*.txt"))
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range files {
		fmt.Println(f)
		findUnicodeWords(f)
	}
}

func findUnicodeWords(file string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error while opening file %s\n", file)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
