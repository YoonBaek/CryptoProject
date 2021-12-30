package blockchain

import (
	"errors"
	"time"

	"github.com/YoonBaek/CryptoProject/utils"
	"github.com/YoonBaek/CryptoProject/wallet"
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

// uTx를 구분하기 위해서는 해당 TxIn이 어느 Tx의 output에서 나왔는지,
// 그리고 해당 Tx의 output 중 어느 것인지에 대한 정보가 필요하다.
// 따라서 TxId와 Index랄 깆도록 뱐걍헤주고, 금액은 어차피 추출하고자 하는
// TxOut에서 가져오면 되기 때문에 제외해준다.
// 즉, 어느 TxOut에서 돈을 뽑아올지 보여주는 표지판이라고 보면 된다.
type TxIn struct {
	TxID  string `json:"transaction id"` // uTx 검증에 사용
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type UTxOut struct {
	TxId   string
	Index  int
	Amount int
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

// make coinbase Tx and Confirm transaction standby in mempool
func (m *mempool) Confirm() []*Tx {
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
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
		TimeStamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

// This function returns UnSpentTrasactionOutput.
// This fuction search all the blocks of the chain,
// looking for transactions of each block.
// Be aware of this function finds block recent to past,
// So this block can check used TxOut without a leak.
// If this function finds ID of TxIn, then it doesn't count
// TxOuts which have same ID.
func UtxOutsByAddr(addr string, b *blockchain) (UtxOuts []*UTxOut) {
	spentCheck := map[string]bool{}
	for _, block := range Blocks(b) {
		for _, tx := range block.TransActions {
			for _, txIn := range tx.TxIns {
				if txIn.Owner == addr {
					spentCheck[txIn.TxID] = true
				}
			}
			for index, txOut := range tx.TxOuts {
				if txOut.Owner != addr {
					continue
				}
				if _, exists := spentCheck[tx.Id]; !exists {
					criterion := false
					if criterion = isOnMempool(tx.Id, index); criterion {
						continue
					}
					UtxOuts = append(UtxOuts, &UTxOut{
						TxId:   tx.Id,
						Index:  index,
						Amount: txOut.Amount,
					})
				}
			}
		}
	}
	return
}

// 남은 잔고의 총액을 반환
func BalanceByAddr(addr string, b *blockchain) (total int) {
	txOuts := UtxOutsByAddr(addr, b)

	for _, txOut := range txOuts {
		total += txOut.Amount
	}
	return
}

// This function returns pointer of calculated Tx
func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddr(from, BlockChain()) < amount {
		return nil, errors.New("not enough money")
	}
	var txOuts []*TxOut
	var txIns []*TxIn

	total := 0
	for _, utxOut := range UtxOutsByAddr(from, BlockChain()) {
		if total >= amount {
			break
		}
		txIns = append(txIns, &TxIn{
			TxID:  utxOut.TxId,
			Index: utxOut.Index,
			Owner: from,
		})
		total += utxOut.Amount
	}
	if change := total - amount; change != 0 {
		txOuts = append(txOuts, &TxOut{from, change})
	}

	txOuts = append(txOuts, &TxOut{to, amount})
	tx := &Tx{
		Id:        "",
		TimeStamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
}

func isOnMempool(id string, index int) (exists bool) {
Outer:
	for _, tx := range Mempool.Txs {
		for _, txIn := range tx.TxIns {
			exists = txIn.TxID == id && txIn.Index == index
			break Outer
		}
	}
	return
}
