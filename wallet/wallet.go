package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/YoonBaek/CryptoProject/utils"
)

func Start() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	targetMessage := "you are amazing"
	hashedMessage := utils.Hash(targetMessage)

	hash, err := hex.DecodeString(hashedMessage)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash)
	utils.HandleErr(err)

	fmt.Println("R:", r)
	fmt.Println("S:", s)
}
