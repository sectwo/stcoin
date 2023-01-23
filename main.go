package main

import (
	"github.com/sectwo/stcoin/cli"
	"github.com/sectwo/stcoin/db"
)

func main() {
	//	blockchain.Blockchain()
	defer db.Close()
	cli.Start()

	//wallet.Wallet()
}
