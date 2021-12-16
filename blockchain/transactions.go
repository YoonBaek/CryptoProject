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
	Owner  string
	Amount int
}

type TxOut struct {
	Owner  string
	Amount int
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
