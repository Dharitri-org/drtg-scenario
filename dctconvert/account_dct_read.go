package dctconvert

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/Dharitri-org/drtg-core/data/dct"
)

// MockDCTData groups together all instances of a token (same token name, different nonces).
type MockDCTData struct {
	TokenIdentifier []byte
	Instances       []*dct.DCToken
	LastNonce       uint64
	Roles           [][]byte
}

const (
	dctIdentifierSeparator  = "-"
	dctRandomSequenceLength = 6
)

// GetTokenBalance returns the DCT balance of the account, specified by the
// token key.
func GetTokenBalance(tokenIdentifier []byte, nonce uint64, source map[string][]byte) (*big.Int, error) {
	tokenData, err := GetTokenData(tokenIdentifier, nonce, source, make(map[string][]byte))
	if err != nil {
		return nil, err
	}

	return tokenData.Value, nil
}

// GetTokenData gets the DCT information related to a token from the storage of the account.
func GetTokenData(tokenIdentifier []byte, nonce uint64, source map[string][]byte, systemAccStorage map[string][]byte) (*dct.DCToken, error) {
	tokenKey := makeTokenKey(tokenIdentifier, nonce)
	return getTokenDataByKey(tokenKey, source, systemAccStorage)
}

func getTokenDataByKey(tokenKey []byte, source map[string][]byte, systemAccStorage map[string][]byte) (*dct.DCToken, error) {
	// default value copied from the protocol
	dctData := &dct.DCToken{
		Value: big.NewInt(0),
	}

	marshaledData := source[string(tokenKey)]
	if len(marshaledData) == 0 {
		return dctData, nil
	}

	err := dctDataMarshalizer.Unmarshal(dctData, marshaledData)
	if err != nil {
		return nil, err
	}

	marshaledData = systemAccStorage[string(tokenKey)]
	if len(marshaledData) == 0 {
		return dctData, nil
	}
	dctDataFromSystemAcc := &dct.DCToken{}
	err = dctDataMarshalizer.Unmarshal(dctDataFromSystemAcc, marshaledData)
	if err != nil {
		return nil, err
	}

	dctData.TokenMetaData = dctDataFromSystemAcc.TokenMetaData

	return dctData, nil
}

// GetTokenRoles returns the roles of the account for the specified tokenName.
func GetTokenRoles(tokenName []byte, source map[string][]byte) ([][]byte, error) {
	tokenRolesKey := makeTokenRolesKey(tokenName)
	tokenRolesData := &dct.DCTRoles{
		Roles: make([][]byte, 0),
	}

	marshaledData := source[string(tokenRolesKey)]
	if len(marshaledData) == 0 {
		return tokenRolesData.Roles, nil
	}

	err := dctDataMarshalizer.Unmarshal(tokenRolesData, marshaledData)
	if err != nil {
		return nil, err
	}

	return tokenRolesData.Roles, nil

}

// GetFullMockDCTData returns the information about all the DCT tokens held by the account.
func GetFullMockDCTData(source map[string][]byte, systemAccStorage map[string][]byte) (map[string]*MockDCTData, error) {
	resultMap := make(map[string]*MockDCTData)
	for key := range source {
		storageKeyBytes := []byte(key)
		if isTokenKey(storageKeyBytes) {
			tokenName, tokenInstance, err := loadMockDCTDataInstance(storageKeyBytes, source, systemAccStorage)
			if err != nil {
				return nil, err
			}
			if tokenInstance.Value.Sign() > 0 {
				resultObj := getOrCreateMockDCTData(tokenName, resultMap)
				resultObj.Instances = append(resultObj.Instances, tokenInstance)
			}
		} else if isNonceKey(storageKeyBytes) {
			tokenName := key[len(dctNonceKeyPrefix):]
			resultObj := getOrCreateMockDCTData(tokenName, resultMap)
			resultObj.LastNonce = big.NewInt(0).SetBytes(source[key]).Uint64()
		} else if isRoleKey(storageKeyBytes) {
			tokenName := key[len(dctRoleKeyPrefix):]
			roles, err := GetTokenRoles([]byte(tokenName), source)
			if err != nil {
				return nil, err
			}
			resultObj := getOrCreateMockDCTData(tokenName, resultMap)
			resultObj.Roles = roles
		}
	}

	return resultMap, nil
}

func extractTokenIdentifierAndNonceDCTWipe(args []byte) ([]byte, uint64) {
	argsSplit := bytes.Split(args, []byte(dctIdentifierSeparator))
	if len(argsSplit) < 2 {
		return args, 0
	}

	if len(argsSplit[1]) <= dctRandomSequenceLength {
		return args, 0
	}

	identifier := []byte(fmt.Sprintf("%s-%s", argsSplit[0], argsSplit[1][:dctRandomSequenceLength]))
	nonce := big.NewInt(0).SetBytes(argsSplit[1][dctRandomSequenceLength:])

	return identifier, nonce.Uint64()
}

// loads and prepared the DCT instance
func loadMockDCTDataInstance(tokenKey []byte, source map[string][]byte, systemAccStorage map[string][]byte) (string, *dct.DCToken, error) {
	tokenInstance, err := getTokenDataByKey(tokenKey, source, systemAccStorage)
	if err != nil {
		return "", nil, err
	}

	tokenNameFromKey := getTokenNameFromKey(tokenKey)
	tokenName, nonce := extractTokenIdentifierAndNonceDCTWipe(tokenNameFromKey)

	if tokenInstance.TokenMetaData == nil {
		tokenInstance.TokenMetaData = &dct.MetaData{
			Name:  tokenName,
			Nonce: nonce,
		}
	}

	return string(tokenName), tokenInstance, nil
}

func getOrCreateMockDCTData(tokenName string, resultMap map[string]*MockDCTData) *MockDCTData {
	resultObj := resultMap[tokenName]
	if resultObj == nil {
		resultObj = &MockDCTData{
			TokenIdentifier: []byte(tokenName),
			Instances:       nil,
			LastNonce:       0,
			Roles:           nil,
		}
		resultMap[tokenName] = resultObj
	}
	return resultObj
}
