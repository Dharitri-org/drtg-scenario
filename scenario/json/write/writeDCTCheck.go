package scenjsonwrite

import (
	oj "github.com/Dharitri-org/drtg-scenario/orderedjson"
	scenmodel "github.com/Dharitri-org/drtg-scenario/scenario/model"
)

func checkDCTDataToOJ(dctItems []*scenmodel.CheckDCTData, moreDCTTokensAllowed bool) *oj.OJsonMap {
	dctItemsOJ := oj.NewMap()
	for _, dctItem := range dctItems {
		dctItemsOJ.Put(dctItem.TokenIdentifier.Original, checkDCTItemToOJ(dctItem))
	}
	if moreDCTTokensAllowed {
		dctItemsOJ.Put("+", stringToOJ(""))
	}
	return dctItemsOJ
}

func checkDCTItemToOJ(dctItem *scenmodel.CheckDCTData) oj.OJsonObject {
	if isCompactCheckDCT(dctItem) {
		return checkBigIntToOJ(dctItem.Instances[0].Balance)
	}

	dctItemOJ := oj.NewMap()

	// instances
	if len(dctItem.Instances) > 0 {
		var convertedList []oj.OJsonObject
		for _, dctInstance := range dctItem.Instances {
			dctInstanceOJ := oj.NewMap()
			appendCheckDCTInstanceToOJ(dctInstance, dctInstanceOJ)
			convertedList = append(convertedList, dctInstanceOJ)
		}
		instancesOJList := oj.OJsonList(convertedList)
		dctItemOJ.Put("instances", &instancesOJList)
	}

	if len(dctItem.LastNonce.Original) > 0 {
		dctItemOJ.Put("lastNonce", checkUint64ToOJ(dctItem.LastNonce))
	}

	// roles
	if len(dctItem.Roles) > 0 {
		var convertedList []oj.OJsonObject
		for _, roleStr := range dctItem.Roles {
			convertedList = append(convertedList, &oj.OJsonString{Value: roleStr})
		}
		rolesOJList := oj.OJsonList(convertedList)
		dctItemOJ.Put("roles", &rolesOJList)
	}
	if len(dctItem.Frozen.Original) > 0 {
		dctItemOJ.Put("frozen", checkUint64ToOJ(dctItem.Frozen))
	}

	return dctItemOJ
}

func appendCheckDCTInstanceToOJ(dctInstance *scenmodel.CheckDCTInstance, targetOj *oj.OJsonMap) {
	targetOj.Put("nonce", uint64ToOJ(dctInstance.Nonce))

	if len(dctInstance.Balance.Original) > 0 {
		targetOj.Put("balance", checkBigIntToOJ(dctInstance.Balance))
	}
	if !dctInstance.Creator.Unspecified && len(dctInstance.Creator.Value) > 0 {
		targetOj.Put("creator", checkBytesToOJ(dctInstance.Creator))
	}
	if !dctInstance.Royalties.Unspecified && len(dctInstance.Royalties.Original) > 0 {
		targetOj.Put("royalties", checkUint64ToOJ(dctInstance.Royalties))
	}
	if !dctInstance.Hash.Unspecified && len(dctInstance.Hash.Value) > 0 {
		targetOj.Put("hash", checkBytesToOJ(dctInstance.Hash))
	}
	if !dctInstance.Uris.IsUnspecified() {
		targetOj.Put("uri", checkValueListToOJ(dctInstance.Uris))
	}
	if !dctInstance.Attributes.Unspecified && len(dctInstance.Attributes.Value) > 0 {
		targetOj.Put("attributes", checkBytesToOJ(dctInstance.Attributes))
	}
}

func isCompactCheckDCT(dctItem *scenmodel.CheckDCTData) bool {
	if len(dctItem.Instances) != 1 {
		return false
	}
	if len(dctItem.Instances[0].Nonce.Original) > 0 {
		return false
	}
	if len(dctItem.Roles) > 0 {
		return false
	}
	if len(dctItem.Frozen.Original) > 0 {
		return false
	}
	return true
}
