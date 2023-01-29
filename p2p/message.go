package p2p

import (
	"encoding/json"

	"github.com/sectwo/stcoin/blockchain"
	"github.com/sectwo/stcoin/utils"
)

type MessageKind int

const (
	//	MessageNewestBlock       MessageKind = 1
	//	MessageAllBlocksRequest  MessageKind = 2
	//	MessageAllBlocksResponse MessageKind = 3

	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
)

type Message struct {
	Kind    MessageKind
	payload []byte
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		payload: utils.ToBytes(payload),
	}
	return utils.ToJson(m)
}

func sendNewestBlock(p *peer) {
	b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)
	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m
}

func handleMsg(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		var payload blockchain.Block
		err := json.Unmarshal(m.payload, &payload)
		utils.HandleErr(err)

	}
}
