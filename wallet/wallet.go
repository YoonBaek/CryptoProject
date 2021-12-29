package wallet

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/YoonBaek/CryptoProject/utils"
)

const (
	// 저번 커밋까지 미리 뽑아 둔 정보들
	// 얘네 들로 복구가 가능할지 지켜봄
	privateKey    string = "30770201010420e2048ed34e6edc3249b3335e7168e6d79310ed07cc8d4b294f33d72062910363a00a06082a8648ce3d030107a144034200048e8f85542b54328d72fadc35f98adc3b3c72cd2882f3ae3e455255d54f8f0081d52891295d63d9b97776a4e31f88f009dde111b97818290e1b58d58b5a77b250"
	hashedMessage string = "11d2d024a8ba69cd45d1b9529016f973d0f4d78fb821c8ca67a2f8a25b09458b"
	signature     string = "c860694e014e3bd28ca0bedaf30adf20b11b75f566f9779a4d6ab1e1bb3a24e21fd9bd0df49dd3660a6c8dd85640ed0cdca3df3ec07880fb94e5756be61d5004"
)

func Start() {
	// 키 복구하기'
	privBytes, err := hex.DecodeString(privateKey)
	utils.HandleErr(err)
	restoredPrivKey, err := x509.ParseECPrivateKey(privBytes)
	utils.HandleErr(err)

	fmt.Println(restoredPrivKey)

	// 공개키 복구하기
	restoredPubKey := restoredPrivKey.PublicKey

	// 서명 복구하기
	sigBytes, err := hex.DecodeString(signature)
	utils.HandleErr(err)

	mid := len(sigBytes) / 2
	rBytes := sigBytes[:mid]
	sBytes := sigBytes[mid:]

	var restoredR, restoredS = big.Int{}, big.Int{}
	restoredR.SetBytes(rBytes)
	restoredS.SetBytes(sBytes)

	fmt.Println(restoredR)
	fmt.Println(restoredS)

	// payload 복구하기
	restoredMsg, err := hex.DecodeString(hashedMessage)
	utils.HandleErr(err)
	// 복구한 키로 검증해보기
	fmt.Println(ecdsa.Verify(&restoredPubKey, restoredMsg, &restoredR, &restoredS))
}
