package dctconvert

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/Dharitri-org/drtg-core/core"
	"github.com/Dharitri-org/drtg-core/marshal"
)

// dctTokenKeyPrefix is the prefix of storage keys belonging to DCT tokens.
var dctTokenKeyPrefix = []byte(core.ProtectedKeyPrefix + core.DCTKeyIdentifier)

// dctRoleKeyPrefix is the prefix of storage keys belonging to DCT roles.
var dctRoleKeyPrefix = []byte(core.ProtectedKeyPrefix + core.DCTRoleIdentifier + core.DCTKeyIdentifier)

// dctNonceKeyPrefix is the prefix of storage keys belonging to DCT nonces.
var dctNonceKeyPrefix = []byte(core.ProtectedKeyPrefix + core.DCTNFTLatestNonceIdentifier)

// dctDataMarshalizer is the global marshalizer to be used for encoding/decoding DCT data
var dctDataMarshalizer = &marshal.GogoProtoMarshalizer{}

// errNegativeValue signals that a negative value has been detected and it is not allowed
var errNegativeValue = errors.New("negative value")

// makeTokenKey creates the storage key corresponding to the given tokenName.
func makeTokenKey(tokenName []byte, nonce uint64) []byte {
	nonceBytes := big.NewInt(0).SetUint64(nonce).Bytes()
	tokenKey := append(dctTokenKeyPrefix, tokenName...)
	tokenKey = append(tokenKey, nonceBytes...)
	return tokenKey
}

// makeTokenRolesKey creates the storage key corresponding to the roles for the
// given tokenName.
func makeTokenRolesKey(tokenName []byte) []byte {
	tokenRolesKey := append(dctRoleKeyPrefix, tokenName...)
	return tokenRolesKey
}

// makeLastNonceKey creates the storage key corresponding to the last nonce of
// the given tokenName.
func makeLastNonceKey(tokenName []byte) []byte {
	tokenNonceKey := append(dctNonceKeyPrefix, tokenName...)
	return tokenNonceKey
}

// isTokenKey returns true if the given storage key belongs to an DCT token.
func isTokenKey(key []byte) bool {
	return bytes.HasPrefix(key, dctTokenKeyPrefix)
}

// isRoleKey returns true if the given storage key belongs to an DCT role.
func isRoleKey(key []byte) bool {
	return bytes.HasPrefix(key, dctRoleKeyPrefix)
}

// isNonceKey returns true if the given storage key belongs to an DCT nonce.
func isNonceKey(key []byte) bool {
	return bytes.HasPrefix(key, dctNonceKeyPrefix)
}

// getTokenNameFromKey extracts the token name from the given storage key; it
// does not check whether the key is indeed a token key or not.
func getTokenNameFromKey(key []byte) []byte {
	return key[len(dctTokenKeyPrefix):]
}
