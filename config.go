package main

import (
  "io/ioutil"
  "encoding/json"
  "log"
)

type SiteConfig struct{
  Url string `json:url`
  Depth int `json:depth`
  Section string `json:section`
  Skip []string `json:skip`
}

type Config []SiteConfig

func GetConfig(path string) Config{
  data,err := ioutil.ReadFile(path)
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

