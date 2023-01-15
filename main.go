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

	// difficulty := 3
	// target := strings.Repeat("0", difficulty)
	// nonce := 1
	// for {
	// 	hash := fmt.Sprintf("%x\n", sha256.Sum256([]byte("hello"+fmt.Sprint(nonce))))
	// 	fmt.Printf("Hash : %s\nTarget : %s\nNonce : %d\n", hash, target, nonce)
	// 	if strings.HasPrefix(hash, target) {
	// 		break
	// 	} else {
	// 		nonce++
	// 	}
	// }

	// fmt.Println(target)

	// hash := sha256.Sum256([]byte("hello"))
	// fmt.Printf("%x\n", hash)
}
