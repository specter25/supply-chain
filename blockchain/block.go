package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp  int64
	Hash       []byte
	PrevHash   []byte
	Nonce      int
	Height     int
	Data       string
	Difficulty int64
}

func Genesis() *Block {
	return CreateBlock("GENESIS-DATA", []byte{}, 0, INITIAL_DIFFICULTY, 0)
}

func CreateBlock(data string, PrevHash []byte, Height int, Difficulty int64, lasttime int64) *Block {
	block := &Block{time.Now().Unix(), []byte{}, PrevHash, 0, Height, data, Difficulty}
	pow := Newproof(block)
	nonce, hash, difficulty := pow.Run(lasttime)
	block.Hash = hash[:]
	block.Nonce = nonce
	block.Difficulty = difficulty
	return block
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)
	Handle(err)
	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	Handle(err)
	return &block
}
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
