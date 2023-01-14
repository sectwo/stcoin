package main

import (
	explore "github.com/sectwo/stcoin/explorer"
	"github.com/sectwo/stcoin/rest"
)

func main() {
	go explore.StartExplorer(4000)
	rest.Start(8080)
}
