package worldmock

import (
	"github.com/Dharitri-org/drtg-core/core"
	vmcommon "github.com/Dharitri-org/drtg-vm-common"
)

// GenerateMockAddress simulates creation of a new address by the protocol.
//
// Not an actual blockchain hook, just a helper method.
func GenerateMockAddress(vmType []byte, creatorAddress []byte, creatorNonce uint64) []byte {
	result := make([]byte, 32)
	result[10] = 0x11
	result[11] = 0x11
	result[12] = 0x11
	result[13] = 0x11
	copy(result[14:29], creatorAddress)

	result[29] = byte(creatorNonce)

	copy(result[30:], creatorAddress[30:])

	if vmType == nil {
		vmType = DefaultVMType
	}

	copy(result[vmcommon.NumInitCharactersForScAddress-core.VMTypeLen:], vmType)
	return result
}
