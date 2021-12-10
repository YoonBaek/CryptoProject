package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/YoonBaek/CryptoProject/explorer"
	"github.com/YoonBaek/CryptoProject/rest"
)

func usage() {
	fmt.Printf("Welcome to Crypto Project\n\n")
	fmt.Printf("Pealse use the following commands:\n\n")
	fmt.Printf("-mode:	Choose between 'rest' and 'explorer'\n")
	fmt.Printf("-port:	Sets the port number of the server\n\n")
	os.Exit(0)
}

func main() {
	mode := flag.String("mode", "rest", "Choose between 'rest' and 'explorer'")
	port := flag.Int("port", 4000, "Sets the port number of the server")

	switch *mode {
	case "rest":
		//start rest API
		rest.Start(*port)
	case "explorer":
		// start html explorer
		explorer.Start(*port)
	default:
		usage()
	}
}
