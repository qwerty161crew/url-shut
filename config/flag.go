package config

import (
	"flag"
	"fmt"
)

var FlagRunAddr string
var RedirectHost string

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "", "address and port to run server")
	flag.StringVar(&RedirectHost, "b", "", "address and port to redirect server")
	flag.Parse()
	fmt.Println(FlagRunAddr)

}
