package blockchain

import (
	"errors"
	"strings"
	"time"

	"github.com/YoonBaek/CryptoProject/blockchain/db"
	"github.com/YoonBaek/CryptoProject/utils"
)

var ErrNotFound = errors.New("block not found")

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"heihgt"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"Nonce"`
	Timestamp    int    `json:"timestamp"`
	TransActions []*Tx  `json:"transactions"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target) {
			// fmt.Printf("Target: %s\nHash: %s\nNonce: %d\nDifficulty: %d\n", target, hash, b.Nonce, b.Difficulty)
			b.Timestamp = int(time.Now().Unix())
			b.Hash = hash
			break
		}
		b.Nonce++
	}
}

func FindBlock(hash string) (*Block, error) {
	blockFromDB := db.LoadBlock(hash)
	if blockFromDB == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockFromDB)
	return block, nil
}

func createBlock(prevHash string, height int) *Block {
	block := Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: difficulty(BlockChain()),
		Nonce:      10,
	}
	block.mine()
	block.TransActions = Mempool.Confirm()
	block.persist()
	return &block
}
