package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

const (
	INITIAL_DIFFICULTY = 3
	STARTING_BALANCE   = 1000
	MINE_RATE          = 1000
)

type Block struct {
	Timestamp int64
	Hash      []byte
	PrevHash  []byte
	Nonce     int
	// Height       int
	Data       string
	Difficulty int
}

func Genesis() *Block {
	return CreateBlock("GENESIS-DATA", []byte{})
}

func CreateBlock(data string, PrevHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte{}, PrevHash, 0, data, INITIAL_DIFFICULTY}
	pow := Newproof(block, 3)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
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
