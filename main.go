package main

import (
	"github.com/sectwo/stcoin/blockchain"
	"github.com/sectwo/stcoin/cli"
	"github.com/sectwo/stcoin/db"
)

func main() {
	defer db.Close()
	blockchain.Blockchain()
	cli.Start()
}
