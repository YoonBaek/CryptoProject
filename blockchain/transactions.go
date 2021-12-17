package blockchain

import (
	"errors"
	"time"

	"github.com/YoonBaek/CryptoProject/utils"
)

const minerReward int = 10

// Mempool
type mempool struct {
	Txs []*Tx
}

// blockchain을 생성할 때는 DB에서 한 번 불러와주는 로직이
// 반드시 필요했지만, Mempool은 그렇지 않기 때문에
// Singleton 패턴을 적용하지 않고 바로 Export 해준다.
var Mempool *mempool = &mempool{}

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

// 남은 잔고의 총액을 반환
func (b *blockchain) BalanceByAddr(addr string) (total int) {
	txOuts := b.TxOutsByAddr(addr)

	for _, txOut := range txOuts {
		total += txOut.Amount
	}
	return
}

// Tx의 메모리를 반환하는 함수
// 1. 남은 balance 보다 많은 amount를 요구하면 error
// 2. owner가 보유중인 과거 Tx들에서 amount를 떼 오고 total에 저장
// 3. 거스름 돈은 owner에게 다시 txout
func makeTx(from, to string, amount int) (*Tx, error) {
	b := BlockChain()
	if b.BalanceByAddr(from) < amount {
		return nil, errors.New("not enough money")
	}
	total, txIns, txOuts := 0, []*TxIn{}, []*TxOut{}
	pastTxOuts := b.TxOutsByAddr(from)
	for _, txOut := range pastTxOuts {
		if total > amount {
			break
		}
		txIns = append(txIns, &TxIn{txOut.Owner, txOut.Amount})
		total += txOut.Amount
	}
	change := total - amount
	// changes to prev owner
	txOut := &TxOut{from, change}
	txOuts = append(txOuts, txOut)
	// totals to new owner
	txOut = &TxOut{to, amount}
	txOuts = append(txOuts, txOut)

	tx := &Tx{
		Id:        "",
		TimeStamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("yoonbaek", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) Confirm() []*Tx {
	coinbase := makeCoinbaseTx("yoonbaek")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
