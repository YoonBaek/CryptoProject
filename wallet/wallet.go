package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/YoonBaek/CryptoProject/utils"
)

const (
	fileName string = "yoonbaek.wallet" // tmp
)

// follows singleton pattern
type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

func (w *wallet) createKey() {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	w.privateKey = privKey
}

func (w *wallet) restoreKey() {
	keyAsBytes, err := os.ReadFile(fileName)
	utils.HandleErr(err)
	restoredKey, err := x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	w.privateKey = restoredKey
}

func (w *wallet) getAddress(pubKey *ecdsa.PublicKey) {
	address := append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)
	w.Address = fmt.Sprintf("%x", address)
}

func Wallet() *wallet {
	// 해당 지갑이 메모리에 없으면
	if w == nil {
		w = &wallet{}
		// 해당 지갑 정보가 있는지 확인하기
		switch hasWallet() {
		case true: // 있다면 파일로부터 불러오기
			w.restoreKey()
		default: // 없다면 비공개키를 생성해 주고 파일에 저장
			w.createKey()
			persistKey(w.privateKey)
		}
		w.getAddress(&w.privateKey.PublicKey)
	}
	return w
}

func hasWallet() bool {
	_, err := os.Stat(fileName)
	return os.IsExist(err)
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = os.WriteFile(fileName, bytes, 0644)
	utils.HandleErr(err)
}
