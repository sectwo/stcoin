package blockchain

import (
	"errors"
	"sync"
	"time"

	"github.com/sectwo/stcoin/utils"
	"github.com/sectwo/stcoin/wallet"
)

const (
	minerReward int = 50
)

type mempool struct {
	//Txs []*Tx
	Txs map[string]*Tx
	m   sync.Mutex
}

var m *mempool // Memory only
var memOnce sync.Once

func Mempool() *mempool {
	memOnce.Do(func() {
		m = &mempool{
			Txs: make(map[string]*Tx),
		}
	})
	return m
}

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type TxIn struct {
	TxID      string `json:"txId"`
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}

type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type UTxOut struct {
	TxID   string `json:"txId"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(t.Id, wallet.Wallet())
	}
}

func validate(tx *Tx) bool {
	valid := true
	for _, txIn := range tx.TxIns {
		prevTx := FindTx(Blockchain(), txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}
		address := prevTx.TxOuts[txIn.Index].Address
		valid = wallet.Verify(txIn.Signature, tx.Id, address)
		if !valid {
			break
		}
	}
	return valid
}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
Outer:
	for _, tx := range Mempool().Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer
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

var ErrorNoMoney = errors.New("Not enought money")
var ErrorNotValid = errors.New("Tx Invalid")

// From A to B , Amount $5
// A =[$5]=> B
// 1. A의 잔액 확인(A의 TxOut) 하고 잔액이 충분하지 않다면 종료
// 2. 충분하다면, A의 잔액을 이전 TxOut으로 부터 계산하여 현재 내 주머니의 돈을 계산
func makeTx(from, to string, amount int) (*Tx, error) {
	// BalanceByAddress로 부터 Balance가 충분하지 않다면 nil
	if BalanceByAddress(Blockchain(), from) < amount {
		return nil, ErrorNoMoney
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	UTxOuts := UTxOutsByAddress(Blockchain(), from)
	// 거래를 위한 TxIns 생성
	for _, uTxOut := range UTxOuts {
		if total >= amount {
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
	tx.sign()
	valid := validate(tx)
	if !valid {
		return nil, ErrorNotValid
	}

	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) (*Tx, error) {
	// Add Transaction to Mempool
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return nil, err
	}
	m.Txs[tx.Id] = tx
	return tx, nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := []*Tx{makeCoinbaseTx(wallet.Wallet().Address)}
	var txs []*Tx
	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}
	txs = append(txs, coinbase...)
	m.Txs = make(map[string]*Tx)
	return txs
}

func (m *mempool) AddPeerTx(tx *Tx) {
	m.m.Lock()
	defer m.m.Unlock()

	m.Txs[tx.Id] = tx
}
