package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/YoonBaek/CryptoProject/utils"
)

func Start() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	fmt.Println(privateKey.D)
}
