package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"

	"github.com/YoonBaek/CryptoProject/utils"
)

const hashedMessage string = "11d2d024a8ba69cd45d1b9529016f973d0f4d78fb821c8ca67a2f8a25b09458b"

func Start() {
	// private, public key 생성
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	keyAsBytes, err := x509.MarshalECPrivateKey(privateKey)
	utils.HandleErr(err)
	fmt.Printf("%x\n", keyAsBytes)
	// 메시지 해싱
	hash, err := hex.DecodeString(hashedMessage)
	utils.HandleErr(err)
	// 메시지 서명하기
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash)
	utils.HandleErr(err)
	signature := append(r.Bytes(), s.Bytes()...)
	fmt.Printf("#%x\n", signature)
	// 메세지 검증하기
	ok := ecdsa.Verify(&privateKey.PublicKey, hash, r, s)
	fmt.Println(ok)
}
