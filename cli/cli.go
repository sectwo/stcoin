package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/sectwo/stcoin/explorer"
	"github.com/sectwo/stcoin/rest"
)

func usage() {
	fmt.Printf("Welcome to SECTWO COIN\n")
	fmt.Printf("Please use the following flags : \n")
	fmt.Printf("-port:		Set the PORT of the service\n")
	fmt.Printf("-mode:		Choose between 'html' and 'rest'\n\n")
	runtime.Goexit()
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 8000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	flag.Parse()

	switch *mode {
	case "rest":
		//start rest api
		rest.Start(*port)
	case "html":
		//start html explorer
		explorer.Start(*port)

	default:
		usage()
	}

	fmt.Println(*port, *mode)

}
