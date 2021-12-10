package blockchain

import (
	"crypto/sha256"
	"fmt"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
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
	return &block
}
