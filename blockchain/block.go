package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sectwo/stcoin/db"
	"github.com/sectwo/stcoin/utils"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	Timestamp    int    `json:"timestamp"`
	Transections []*Tx  `json:"transections"`
}

func persistBlock(b *Block) {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

var ErrNotFound = errors.New("block not found")

func (b *Block) restore(data []byte) {
	utils.FromByte(b, data)
}

func FindBlock(hash string) (*Block, error) {
	blockByte := db.Block(hash)
	if blockByte == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockByte)
	return block, nil
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		hash := utils.Hash(b)
		fmt.Printf("\n\n\nTarget:%s\nHash:%s\nNonce:%d\n\n\n", target, hash, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			b.Timestamp = int(time.Now().Unix())
			break
		} else {
			b.Nonce++

		}
	}
}

func createBlock(prevHash string, height, diff int) *Block {
	block := &Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: diff,
		Nonce:      0,
	}
	block.mine()
	block.Transections = Mempool.TxToConfirm()
	persistBlock(block)
	return block
}
