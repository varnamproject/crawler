package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type SiteConfig struct {
	Url     string   `json:url`
	Depth   int      `json:depth`
	Section string   `json:section`
	Skip    []string `json:skip`
}

type Config struct {
	Unicode string       `json:unicode`
	Sites   []SiteConfig `json:sites`
}

func GetConfig(path string) Config {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Unable to read file")
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Unable to parse json file")
	}
	return config
}
