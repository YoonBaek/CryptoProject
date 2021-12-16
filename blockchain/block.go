package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/YoonBaek/CryptoProject/blockchain/db"
	"github.com/YoonBaek/CryptoProject/utils"
)

// for test
const difficulty int = 2

type Block struct {
	Data       string `json:"data"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"prevHash,omitempty"`
	Height     int    `json:"heihgt"`
	Difficulty int    `json:"difficulty"`
	Nonce      int    `json:"Nonce"`
	Timestamp  int    `json:"timestamp"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

var ErrNotFound = errors.New("block not found")

func FindBlock(hash string) (*Block, error) {
	blockFromDB := db.LoadBlock(hash)
	if blockFromDB == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockFromDB)
	return block, nil
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

func createBlock(data, prevHash string, height int) *Block {
	block := Block{
		Data:       data,
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: BlockChain().difficulty(),
		Nonce:      10,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.Hash = hash
	block.mine()
	block.persist()
	return &block
}
