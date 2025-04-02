package config

import (
	"flag"
	"url-shortener/pkg/logger"
)

var FlagRunAddr string
var RedirectHost string
var FileUrl string

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "", "address and port to run server")
	flag.StringVar(&RedirectHost, "b", "", "address and port to redirect server")
	flag.StringVar(&FileUrl, "f", "", "file path")
	flag.Parse()
	logger.Info("Parse flags: ", FlagRunAddr, RedirectHost, FileUrl)
}
