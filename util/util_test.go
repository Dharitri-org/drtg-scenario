package util

import (
	"math/big"
	"testing"

	mj "github.com/Dharitri-org/drtg-scenario/model"
	"github.com/stretchr/testify/require"
)

func Test_CreateMultiTransferData_SingleTransfer(t *testing.T) {
	dctTransfers := make([]*mj.DCTTxData, 0)
	dctTransfer := &mj.DCTTxData{
		Nonce:           mj.JSONUint64{Value: 0},
		TokenIdentifier: mj.JSONBytesFromString{Value: []byte("TOK1-abcdef")},
		Value:           mj.JSONBigInt{Value: big.NewInt(100)},
	}
	dctTransfers = append(dctTransfers, dctTransfer)
	data := CreateMultiTransferData(
		[]byte("receiverAddress"),
		dctTransfers, "function1",
		[][]byte{
			[]byte("arg1"),
			[]byte("arg2")},
	)
	require.Equal(t, "MultiDCTNFTTransfer@726563656976657241646472657373@01@544f4b312d616263646566@@64@66756e6374696f6e31@61726731@61726732", string(data))
}

func Test_CreateMultiTransferData_MultipleTransfers(t *testing.T) {
	dctTransfers := make([]*mj.DCTTxData, 0)
	dctTransfer1 := &mj.DCTTxData{
		Nonce:           mj.JSONUint64{Value: 2},
		TokenIdentifier: mj.JSONBytesFromString{Value: []byte("TOK1-abcdef")},
		Value:           mj.JSONBigInt{Value: big.NewInt(100)},
	}
	dctTransfer2 := &mj.DCTTxData{
		Nonce:           mj.JSONUint64{Value: 14},
		TokenIdentifier: mj.JSONBytesFromString{Value: []byte("TOK2-abcdef")},
		Value:           mj.JSONBigInt{Value: big.NewInt(396)},
	}

	dctTransfers = append(dctTransfers, dctTransfer1, dctTransfer2)
	data := CreateMultiTransferData(
		[]byte("receiverAddress"),
		dctTransfers, "function1",
		[][]byte{
			[]byte("arg1"),
			[]byte("arg2")},
	)
	require.Equal(t, "MultiDCTNFTTransfer@726563656976657241646472657373@02@544f4b312d616263646566@02@64@544f4b322d616263646566@0e@018c@66756e6374696f6e31@61726731@61726732", string(data))
}
