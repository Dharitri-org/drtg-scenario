package scenjsonmodel

// DCTTxData models the transfer of tokens in a tx
type DCTTxData struct {
	TokenIdentifier JSONBytesFromString
	Nonce           JSONUint64
	Value           JSONBigInt
}

// DCTInstance models an instance of an NFT/SFT, with its own nonce
type DCTInstance struct {
	Nonce      JSONUint64
	Balance    JSONBigInt
	Creator    JSONBytesFromString
	Royalties  JSONUint64
	Hash       JSONBytesFromString
	Uris       JSONValueList
	Attributes JSONBytesFromTree
}

// DCTData models an account holding an DCT token
type DCTData struct {
	TokenIdentifier JSONBytesFromString
	Instances       []*DCTInstance
	LastNonce       JSONUint64
	Roles           []string
	Frozen          JSONUint64
}

// CheckDCTInstance checks an instance of an NFT/SFT, with its own nonce
type CheckDCTInstance struct {
	Nonce      JSONUint64
	Balance    JSONCheckBigInt
	Creator    JSONCheckBytes
	Royalties  JSONCheckUint64
	Hash       JSONCheckBytes
	Uris       JSONCheckValueList
	Attributes JSONCheckBytes
}

// NewCheckDCTInstance creates an instance with all fields unspecified.
func NewCheckDCTInstance() *CheckDCTInstance {
	return &CheckDCTInstance{
		Nonce:      JSONUint64Zero(),
		Balance:    JSONCheckBigIntUnspecified(),
		Creator:    JSONCheckBytesUnspecified(),
		Royalties:  JSONCheckUint64Unspecified(),
		Hash:       JSONCheckBytesUnspecified(),
		Uris:       JSONCheckValueListUnspecified(),
		Attributes: JSONCheckBytesUnspecified(),
	}
}

// CheckDCTData checks the DCT tokens held by an account
type CheckDCTData struct {
	TokenIdentifier JSONBytesFromString
	Instances       []*CheckDCTInstance
	LastNonce       JSONCheckUint64
	Roles           []string
	Frozen          JSONCheckUint64
}
