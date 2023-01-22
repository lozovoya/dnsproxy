package main

import (
	"dnproxier/server/app"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"log"
	"os"
	"strconv"
)

type Conf struct {
	Hosts []string
	Port  int `json:",default=53"`
}

func main() {
	var config Conf
	if err := conf.Load("/dnsproxy/config.yaml", &config); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if len(config.Hosts) == 0 {
		log.Printf("no upstream hosts")
		os.Exit(1)
	}
	port := fmt.Sprintf(":%s", strconv.Itoa(config.Port))
	dnsProxy, err := app.NewApp(port, config.Hosts)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	dnsProxy.ListenAndServe()
}
