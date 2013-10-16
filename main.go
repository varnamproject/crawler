package main

import (
"flag"
"fmt"
)

var(
  config = flag.String("c","./config.json","Configuration file for crawler")
)

func main(){
  flag.Parse()
  siteConfigs := GetConfig(*config)
  fmt.Printf("No of sites to crawl : %d\n", len(siteConfigs))
}
