package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
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

func (b *Block) toBytes() []byte {
	var blockBuffer bytes.Buffer
	encoder := gob.NewEncoder(&blockBuffer)
	err := encoder.Encode(b)
	utils.HandleErr(err)
	return blockBuffer.Bytes()
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, b.toBytes())
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
