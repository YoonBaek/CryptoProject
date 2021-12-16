package blockchain

import (
	"time"

	"github.com/YoonBaek/CryptoProject/utils"
)

const minerReward int = 10

// TransAction
type Tx struct {
	Id        string   `json:"id"`
	TimeStamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`  // 누가 보냈는지
	TxOuts    []*TxOut `json:"txOuts"` // 누가 받았는지
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

type TxIn struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"COINBASE", minerReward},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		Id:        "",
		TimeStamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

// TxOut들을 모아놓은 슬라이스 만들기
func (b *blockchain) txOuts() (txOuts []*TxOut) {
	for _, block := range b.Blocks() {
		for _, tx := range block.TransActions {
			txOuts = append(txOuts, tx.TxOuts...)
		}
	}
	return
}

// TxOut 슬라이스 중 address에 맞는 것 걸러주기 -> export
func (b *blockchain) TxOutsByAddr(addr string) (txOutsByAddr []*TxOut) {
	txOuts := b.txOuts()

	for _, txOut := range txOuts {
		if txOut.Owner != addr {
			continue
		}
		txOutsByAddr = append(txOutsByAddr, txOut)
	}
	return
}

func (b *blockchain) BalanceByAddr(addr string) (total int) {
	txOuts := b.TxOutsByAddr(addr)

	for _, txOut := range txOuts {
		total += txOut.Amount
	}
	return
}
