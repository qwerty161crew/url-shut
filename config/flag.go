package config

import (
	"flag"
	"fmt"
)

var FlagRunAddr string

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "", "address and port to run server")
	flag.Parse()
	fmt.Println(FlagRunAddr)

}
