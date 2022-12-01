package rlp

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestEncodeAndDecodeLogData(t *testing.T) {
	tests := []struct {
		name    string
		l       types.Log
		wantErr bool
	}{
		{
			name: "test 1",
			l: types.Log{
				Data: []byte{0x01, 0x02, 0x03, 0x04},
				Topics: []common.Hash{
					common.HexToHash("9DBB0E7DDA3E09710CE75B801ADDC87CF9D9C6C581641B3275FCA409AD086C62"),
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000599"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeLogData(tt.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeLogData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			decoded, err := DecodeLogData(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeLogData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Len(t, decoded.Data, len(tt.l.Data), "data length")
			assert.ElementsMatch(t, decoded.Data, tt.l.Data, "data array")
			if assert.Len(t, decoded.Topics, len(tt.l.Topics), "topics count") {
				for i := range decoded.Topics {
					h := common.BytesToHash(decoded.Topics[i][:])
					assert.Equal(t, h, tt.l.Topics[i], "topic[%d]", i)
				}
			}
		})
	}
}
