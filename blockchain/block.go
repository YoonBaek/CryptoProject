package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/YoonBaek/CryptoProject/blockchain/db"
	"github.com/YoonBaek/CryptoProject/utils"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"heihgt"`
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

func createBlock(data, prevHash string, height int) *Block {
	block := Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.Hash = hash
	block.persist()
	return &block
}
