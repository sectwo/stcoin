package p2p

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sectwo/stcoin/blockchain"
	"github.com/sectwo/stcoin/utils"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
	MessageNewBlockNotify
	MessageNewTxNotify
	MessageNewPeerNotify
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: utils.ToJson(payload),
	}
	return utils.ToJson(m)
}

func sendNewestBlock(p *peer) {
	fmt.Printf("sending newest block to %s\n", p.key)
	b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)
	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m
}

func requestAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- m
}

func sendAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksResponse, blockchain.Blocks(blockchain.Blockchain()))
	p.inbox <- m
}

func notifyNewBlock(b *blockchain.Block, p *peer) {
	m := makeMessage(MessageNewBlockNotify, b)
	p.inbox <- m
}

func notifyNewTx(tx *blockchain.Tx, p *peer) {
	m := makeMessage(MessageNewTxNotify, tx)
	p.inbox <- m
}

func notifyNewPeer(address string, p *peer) {
	m := makeMessage(MessageNewPeerNotify, address)
	p.inbox <- m
}

func handleMsg(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		fmt.Printf("Recevied the newest block from %s\n", p.key)
		var payload blockchain.Block
		err := json.Unmarshal(m.Payload, &payload)
		utils.HandleErr(err)

		b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
		utils.HandleErr(err)
		if payload.Height >= b.Height {
			fmt.Printf("Requesting all block from %s\n", p.key)
			requestAllBlocks(p)
		} else {
			fmt.Printf("Sending newest block to %s\n", p.key)
			sendNewestBlock(p)
		}
	case MessageAllBlocksRequest:
		fmt.Printf("%s wants all the blocks.\n", p.key)
		sendAllBlocks(p)
	case MessageAllBlocksResponse:
		fmt.Printf("Recevied all the block from %s\n", p.key)
		var payload []*blockchain.Block
		err := json.Unmarshal(m.Payload, &payload)
		utils.HandleErr(err)
		blockchain.Blockchain().Replace(payload)
	case MessageNewBlockNotify:
		var payload *blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Blockchain().AddPeerBlock(payload)
	case MessageNewTxNotify:
		var payload *blockchain.Tx
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Mempool().AddPeerTx(payload)
	case MessageNewPeerNotify:
		var payload string
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		parts := strings.Split(payload, ":")
		AddPeer(parts[0], parts[1], parts[2], false)
	}
}
