package blockchain

import (
	"errors"
	"time"

	"github.com/sectwo/stcoin/utils"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{} // Memory only

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

type TxIn struct {
	TxID  string `json:"txId"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type UTxOut struct {
	TxID   string `json:"txId"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
			}
			// exists = input.TxID == uTxOut.TxID && input.Index == uTxOut.Index
		}
	}
	return exists
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

// From A to B , Amount $5
// A =[$5]=> B
// 1. A의 잔액 확인(A의 TxOut) 하고 잔액이 충분하지 않다면 종료
// 2. 충분하다면, A의 잔액을 이전 TxOut으로 부터 계산하여 현재 내 주머니의 돈을 계산
func makeTx(from, to string, amount int) (*Tx, error) {
	// BalanceByAddress로 부터 Balance가 충분하지 않다면 nil
	if Blockchain().BalanceByAddress(from) < amount {
		return nil, errors.New("Not enought money")
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	UTxOuts := Blockchain().UTxOutsByAddress(from)

	// 거래를 위한 TxIns 생성
	for _, uTxOut := range UTxOuts {
		if total > amount {
			break
		}
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	// 잔돈 계산(돌려줘야할 잔돈)
	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}
	// 거래 후 TxOut 생성(잔돈을 가진)
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)

	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()

	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	// Add Transaction to Mempool
	tx, err := makeTx("isaac", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := []*Tx{makeCoinbaseTx("isaac")}
	txs := m.Txs
	txs = append(txs, coinbase...)
	m.Txs = nil
	return txs
}
