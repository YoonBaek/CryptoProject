package main

import (
	"github.com/YoonBaek/CryptoProject/blockchain"
	"github.com/YoonBaek/CryptoProject/blockchain/db"
	"github.com/YoonBaek/CryptoProject/cli"
)

func main() {
	defer db.Close()
	blockchain.BlockChain()
	cli.Start()
}
