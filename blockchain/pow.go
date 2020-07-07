package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math"
	"math/big"
	"time"
)

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func Newproof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-b.Difficulty))
	pow := &ProofOfWork{b, target}
	return pow
}
func (pow *ProofOfWork) InitData(nonce int, timestamp int64, difficulty int64) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			[]byte(pow.Block.Data),
			ToHex(int64(nonce)),
			ToHex((int64(pow.Block.Difficulty))),
			ToHex(int64(timestamp)),
		}, []byte{},
	)
	return data
}
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
func (pow *ProofOfWork) Run(lastTime int64) (int, []byte, int64) {
	var intHash big.Int
	var hash [32]byte
	nonce := 0
	var difficulty int64
	for nonce < math.MaxInt64 {
		timestamp := time.Now().Unix()
		difficulty = pow.Block.adjustDIfficulty(timestamp, lastTime)
		data := pow.InitData(nonce, timestamp, difficulty)

		hash = sha256.Sum256(data)
		// fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
		// fmt.Println()
	}
	return nonce, hash[:], difficulty
}
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := pow.InitData(pow.Block.Nonce, pow.Block.Timestamp, pow.Block.Difficulty)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	return intHash.Cmp(pow.Target) == -1
}

func (b *Block) adjustDIfficulty(timestamp int64, lastTime int64) int64 {

	var value int64

	difficulty := b.Difficulty
	if difficulty < 1 {
		return 1
	}
	difference := timestamp - lastTime

	if difference > MINE_RATE {
		value = difficulty - 1
	}

	if difficulty < MINE_RATE {
		value = difference + 1
	}

	b.Difficulty = value

	return value

}
func (pow *ProofOfWork) modifytarget(difficulty int64) {

	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))
	pow.Target = target
}
