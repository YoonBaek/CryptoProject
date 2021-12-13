package main

import (
	"github.com/YoonBaek/CryptoProject/blockchain"
	"github.com/YoonBaek/CryptoProject/blockchain/db"
	"github.com/YoonBaek/CryptoProject/cli"
)

func main() {
	defer db.Close()
	b := blockchain.BlockChain()
	b.AddBlock("First")
	b.AddBlock("Second")
	b.AddBlock("Third")
	cli.Start()
}
