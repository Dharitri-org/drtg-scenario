package worldmock

import (
	"fmt"
	"math/big"

	"github.com/Dharitri-org/drtg-core/core"
	"github.com/Dharitri-org/drtg-core/core/check"
	"github.com/Dharitri-org/drtg-core/data/dct"
	"github.com/Dharitri-org/drtg-core/data/vm"
	scenmodel "github.com/Dharitri-org/drtg-scenario/scenario/model"
	"github.com/Dharitri-org/drtg-scenario/worldmock/dctconvert"
	vmcommon "github.com/Dharitri-org/drtg-vm-common"
)

// GetTokenBalance returns the DCT balance of an account for the given token
// key (token keys are built from the token identifier using MakeTokenKey).
func (bf *BuiltinFunctionsWrapper) GetTokenBalance(address []byte, tokenIdentifier []byte, nonce uint64) (*big.Int, error) {
	account := bf.World.AcctMap.GetAccount(address)
	if check.IfNil(account) {
		return big.NewInt(0), nil
	}
	return dctconvert.GetTokenBalance(tokenIdentifier, nonce, account.Storage)
}

// GetTokenData gets the DCT information related to a token from the storage of an account
// (token keys are built from the token identifier using MakeTokenKey).
func (bf *BuiltinFunctionsWrapper) GetTokenData(address []byte, tokenIdentifier []byte, nonce uint64) (*dct.DCToken, error) {
	account := bf.World.AcctMap.GetAccount(address)
	if check.IfNil(account) {
		return &dct.DCToken{
			Value: big.NewInt(0),
		}, nil
	}
	systemAccStorage := make(map[string][]byte)
	systemAcc := bf.World.AcctMap.GetAccount(vmcommon.SystemAccountAddress)
	if systemAcc != nil {
		systemAccStorage = systemAcc.Storage
	}
	return account.GetTokenData(tokenIdentifier, nonce, systemAccStorage)
}

// SetTokenData sets the DCT information related to a token from the storage of an account
// (token keys are built from the token identifier using MakeTokenKey).
func (bf *BuiltinFunctionsWrapper) SetTokenData(address []byte, tokenIdentifier []byte, nonce uint64, tokenData *dct.DCToken) error {
	account := bf.World.AcctMap.GetAccount(address)
	if check.IfNil(account) {
		return nil
	}
	return account.SetTokenData(tokenIdentifier, nonce, tokenData)
}

// PerformDirectDCTTransfer calls the real DCTTransfer function immediately;
// only works for in-shard transfers for now, but it will be expanded to
// cross-shard.
// TODO rewrite to simulate what the SCProcessor does when executing a tx with
// data "DCTTransfer@token@value@contractfunc@contractargs..."
// TODO this function duplicates code from host.ExecuteDCTTransfer(), must refactor
func (bf *BuiltinFunctionsWrapper) PerformDirectDCTTransfer(
	sender []byte,
	receiver []byte,
	token []byte,
	nonce uint64,
	value *big.Int,
	callType vm.CallType,
	gasLimit uint64,
	gasPrice uint64,
) (uint64, error) {
	dctTransferInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    callType,
			GasPrice:    gasPrice,
			GasProvided: gasLimit,
			GasLocked:   0,
		},
		RecipientAddr:     receiver,
		Function:          core.BuiltInFunctionDCTTransfer,
		AllowInitFunction: false,
	}

	if nonce > 0 {
		dctTransferInput.Function = core.BuiltInFunctionDCTNFTTransfer
		dctTransferInput.RecipientAddr = dctTransferInput.CallerAddr
		nonceAsBytes := big.NewInt(0).SetUint64(nonce).Bytes()
		dctTransferInput.Arguments = append(dctTransferInput.Arguments, token, nonceAsBytes, value.Bytes(), receiver)
	} else {
		dctTransferInput.Arguments = append(dctTransferInput.Arguments, token, value.Bytes())
	}

	vmOutput, err := bf.ProcessBuiltInFunction(dctTransferInput)
	if err != nil {
		return 0, err
	}

	if vmOutput.ReturnCode != vmcommon.Ok {
		return 0, fmt.Errorf(
			"DCTtransfer failed: retcode = %d, msg = %s",
			vmOutput.ReturnCode,
			vmOutput.ReturnMessage)
	}

	return vmOutput.GasRemaining, nil
}

// PerformDirectMultiDCTTransfer -
func (bf *BuiltinFunctionsWrapper) PerformDirectMultiDCTTransfer(
	sender []byte,
	receiver []byte,
	dctTransfers []*scenmodel.DCTTxData,
	callType vm.CallType,
	gasLimit uint64,
	gasPrice uint64,
) (uint64, error) {
	nrTransfers := len(dctTransfers)
	nrTransfersAsBytes := big.NewInt(0).SetUint64(uint64(nrTransfers)).Bytes()

	multiTransferInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    callType,
			GasPrice:    gasPrice,
			GasProvided: gasLimit,
			GasLocked:   0,
		},
		RecipientAddr:     sender,
		Function:          core.BuiltInFunctionMultiDCTNFTTransfer,
		AllowInitFunction: false,
	}
	multiTransferInput.Arguments = append(multiTransferInput.Arguments, receiver, nrTransfersAsBytes)

	for i := 0; i < nrTransfers; i++ {
		token := dctTransfers[i].TokenIdentifier.Value
		nonceAsBytes := big.NewInt(0).SetUint64(dctTransfers[i].Nonce.Value).Bytes()
		value := dctTransfers[i].Value.Value

		multiTransferInput.Arguments = append(multiTransferInput.Arguments, token, nonceAsBytes, value.Bytes())
	}

	vmOutput, err := bf.ProcessBuiltInFunction(multiTransferInput)
	if err != nil {
		return 0, err
	}

	if vmOutput.ReturnCode != vmcommon.Ok {
		return 0, fmt.Errorf(
			"MultiDCTtransfer failed: retcode = %d, msg = %s",
			vmOutput.ReturnCode,
			vmOutput.ReturnMessage)
	}

	return vmOutput.GasRemaining, nil
}
