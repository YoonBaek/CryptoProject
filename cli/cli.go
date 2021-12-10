package cli

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
	fmt.Printf("-mode:	Choose between 'rest' and 'explorer'\n	If you want both, 'greedy' is ready for you\n")
	fmt.Printf("-port:	Sets the port number of the server\n\n")
	os.Exit(0)
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}
	mode := flag.String("mode", "rest", "Choose between 'rest' and 'explorer'\nIf you want both, 'greedy' is ready for you\n")
	port := flag.Int("port", 4000, "Sets the port number of the server")
	flag.Parse()

	switch *mode {
	case "rest":
		//start rest API
		rest.Start(*port)
	case "explorer":
		// start html explorer
		explorer.Start(*port)
	case "greedy":
		go rest.Start(*port)
		explorer.Start(*port + 1)
	default:
		usage()
	}
}
