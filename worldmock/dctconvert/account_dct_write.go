package dctconvert

import (
	"math/big"

	"github.com/Dharitri-org/drtg-core/core"
	"github.com/Dharitri-org/drtg-core/data/dct"
	scenmodel "github.com/Dharitri-org/drtg-scenario/scenario/model"
	"github.com/Dharitri-org/drtg-vm-common/builtInFunctions"
)

// MakeDCTUserMetadataBytes creates metadata byte slice
func MakeDCTUserMetadataBytes(frozen bool) []byte {
	metadata := &builtInFunctions.DCTUserMetadata{
		Frozen: frozen,
	}

	return metadata.ToBytes()
}

// WriteScenariosDCTToStorage writes the Scenarios DCT data to the provided storage map
func WriteScenariosDCTToStorage(dctData []*scenmodel.DCTData, destination map[string][]byte) error {
	for _, scenDCTData := range dctData {
		tokenIdentifier := scenDCTData.TokenIdentifier.Value
		isFrozen := scenDCTData.Frozen.Value > 0
		for _, instance := range scenDCTData.Instances {
			tokenNonce := instance.Nonce.Value
			tokenKey := makeTokenKey(tokenIdentifier, tokenNonce)
			tokenBalance := instance.Balance.Value
			var uris [][]byte
			for _, jsonUri := range instance.Uris.Values {
				uris = append(uris, jsonUri.Value)
			}
			tokenData := &dct.DCToken{
				Value:      tokenBalance,
				Type:       uint32(core.Fungible),
				Properties: MakeDCTUserMetadataBytes(isFrozen),
				TokenMetaData: &dct.MetaData{
					Name:       []byte{},
					Nonce:      tokenNonce,
					Creator:    instance.Creator.Value,
					Royalties:  uint32(instance.Royalties.Value),
					Hash:       instance.Hash.Value,
					URIs:       uris,
					Attributes: instance.Attributes.Value,
				},
			}
			err := setTokenDataByKey(tokenKey, tokenData, destination)
			if err != nil {
				return err
			}
		}
		err := SetLastNonce(tokenIdentifier, scenDCTData.LastNonce.Value, destination)
		if err != nil {
			return err
		}
		err = SetTokenRolesAsStrings(tokenIdentifier, scenDCTData.Roles, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetTokenData sets the DCT information related to a token into the storage of the account.
func setTokenDataByKey(tokenKey []byte, tokenData *dct.DCToken, destination map[string][]byte) error {
	marshaledData, err := dctDataMarshalizer.Marshal(tokenData)
	if err != nil {
		return err
	}
	destination[string(tokenKey)] = marshaledData
	return nil
}

// SetTokenData sets the token data
func SetTokenData(tokenIdentifier []byte, nonce uint64, tokenData *dct.DCToken, destination map[string][]byte) error {
	tokenKey := makeTokenKey(tokenIdentifier, nonce)
	return setTokenDataByKey(tokenKey, tokenData, destination)
}

// SetTokenRoles sets the specified roles to the account, corresponding to the given tokenIdentifier.
func SetTokenRoles(tokenIdentifier []byte, roles [][]byte, destination map[string][]byte) error {
	tokenRolesKey := makeTokenRolesKey(tokenIdentifier)
	tokenRolesData := &dct.DCTRoles{
		Roles: roles,
	}

	marshaledData, err := dctDataMarshalizer.Marshal(tokenRolesData)
	if err != nil {
		return err
	}

	destination[string(tokenRolesKey)] = marshaledData
	return nil
}

// SetTokenRolesAsStrings sets the specified roles to the account, corresponding to the given tokenIdentifier.
func SetTokenRolesAsStrings(tokenIdentifier []byte, rolesAsStrings []string, destination map[string][]byte) error {
	roles := make([][]byte, len(rolesAsStrings))
	for i := 0; i < len(roles); i++ {
		roles[i] = []byte(rolesAsStrings[i])
	}

	return SetTokenRoles(tokenIdentifier, roles, destination)
}

// SetLastNonce writes the last nonce of a specified DCT into the storage.
func SetLastNonce(tokenIdentifier []byte, lastNonce uint64, destination map[string][]byte) error {
	tokenNonceKey := makeLastNonceKey(tokenIdentifier)
	nonceBytes := big.NewInt(0).SetUint64(lastNonce).Bytes()
	destination[string(tokenNonceKey)] = nonceBytes
	return nil
}

// SetTokenBalance sets the DCT balance of the account, specified by the token
// key.
func SetTokenBalance(tokenIdentifier []byte, nonce uint64, balance *big.Int, destination map[string][]byte) error {
	tokenKey := makeTokenKey(tokenIdentifier, nonce)
	tokenData, err := getTokenDataByKey(tokenKey, destination, make(map[string][]byte))
	if err != nil {
		return err
	}

	if balance.Sign() < 0 {
		return errNegativeValue
	}

	tokenData.Value = balance
	return setTokenDataByKey(tokenKey, tokenData, destination)
}
