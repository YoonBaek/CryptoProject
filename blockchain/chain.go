package blockchain

import (
	"sync"

	"github.com/YoonBaek/CryptoProject/blockchain/db"
	"github.com/YoonBaek/CryptoProject/utils"
)

type blockchain struct {
	LatestHash string `json:"latestHash"`
	Height     int    `json:"height"`
}

var (
	b    *blockchain
	once sync.Once
)

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveCheckpoint(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.LatestHash, b.Height+1)
	b.LatestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func BlockChain() *blockchain {
	if b == nil {
		// GetBlockChain을 통해 b 인스턴스가 실행되는 상황을 가정해보자.
		// 그리고 수 많은 go routine 들이 GetBlockChain을 호출한다고 할 때,
		// b 인스턴스는 한 번만 제대로 초기화 되면 된다.
		// 그럴 때 동기 처리를 도와주는 sync 패키지를 활용한다.
		// once.Do()는 메서드 인자에 들어온 함수를 딱 한 번만 실행하는 기능이다.
		once.Do(func() {
			b = &blockchain{"", 0}
			checkpoint := db.LoadCheckpoint()
			if checkpoint == nil {
				b.AddBlock("Genesis Block")
				return
			}
			// restore b from bytes
			b.restore(checkpoint)
		})
	}
	return b
}

// 쭉 prevHash를 타고 가면서 이전 블록을 불러와 slice에 담고,
// 해당 슬라이스를 반환하기
func (b *blockchain) Blocks() (blocks []*Block) {
	hashMarker := b.LatestHash
	for {
		if hashMarker == "" {
			break
		}
		block, _ := FindBlock(hashMarker)
		blocks = append(blocks, block)
		hashMarker = block.PrevHash
	}
	return
}
