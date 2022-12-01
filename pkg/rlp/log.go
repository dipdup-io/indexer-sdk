package rlp

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// Log -
type Log struct {
	Topics [][32]byte
	Data   []byte
}

// EncodeLogData -
func EncodeLogData(l types.Log) ([]byte, error) {
	topics := make([][32]byte, 0)
	for i := range l.Topics {
		var topic [32]byte
		copy(topic[:], l.Topics[i].Bytes())
		topics = append(topics, topic)
	}

	data := Log{
		Data:   l.Data,
		Topics: topics,
	}

	return rlp.EncodeToBytes(data)
}

// DecodeLogData -
func DecodeLogData(data []byte) (Log, error) {
	var l Log
	err := rlp.DecodeBytes(data, &l)
	return l, err
}
