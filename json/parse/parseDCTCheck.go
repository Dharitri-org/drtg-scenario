package scenjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/Dharitri-org/drtg-scenario/model"
	oj "github.com/Dharitri-org/drtg-scenario/orderedjson"
)

func (p *Parser) processCheckDCTData(
	tokenName mj.JSONBytesFromString,
	dctDataRaw oj.OJsonObject) (*mj.CheckDCTData, error) {

	switch data := dctDataRaw.(type) {
	case *oj.OJsonString:
		// simple string representing balance "400,000,000,000"
		dctData := mj.CheckDCTData{
			TokenIdentifier: tokenName,
		}
		balance, err := p.processCheckBigInt(dctDataRaw, bigIntUnsignedBytes)
		if err != nil {
			return nil, fmt.Errorf("invalid DCT balance: %w", err)
		}
		dctData.Instances = []*mj.CheckDCTInstance{
			{
				Nonce:   mj.JSONUint64Zero(),
				Balance: balance,
			},
		}
		return &dctData, nil
	case *oj.OJsonMap:
		return p.processCheckDCTDataMap(tokenName, data)
	default:
		return nil, errors.New("invalid JSON object for DCT")
	}
}

// Map containing DCT fields, e.g.:
//
//	{
//		"instances": [ ... ],
//	 "lastNonce": "5",
//		"frozen": "true"
//	}
func (p *Parser) processCheckDCTDataMap(tokenName mj.JSONBytesFromString, dctDataMap *oj.OJsonMap) (*mj.CheckDCTData, error) {
	dctData := mj.CheckDCTData{
		TokenIdentifier: tokenName,
	}
	// var err error
	firstInstance := &mj.CheckDCTInstance{
		Nonce:      mj.JSONUint64Zero(),
		Balance:    mj.JSONCheckBigIntUnspecified(),
		Creator:    mj.JSONCheckBytesUnspecified(),
		Royalties:  mj.JSONCheckUint64Unspecified(),
		Hash:       mj.JSONCheckBytesUnspecified(),
		Uris:       mj.JSONCheckValueListUnspecified(),
		Attributes: mj.JSONCheckBytesUnspecified(),
	}
	firstInstanceLoaded := false
	var explicitInstances []*mj.CheckDCTInstance

	for _, kvp := range dctDataMap.OrderedKV {
		// it is allowed to load the instance directly, fields set to the first instance
		instanceFieldLoaded, err := p.tryProcessCheckDCTInstanceField(kvp, firstInstance)
		if err != nil {
			return nil, fmt.Errorf("invalid account DCT instance field: %w", err)
		}
		if instanceFieldLoaded {
			firstInstanceLoaded = true
		} else {
			switch kvp.Key {
			case "instances":
				explicitInstances, err = p.processCheckDCTInstances(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account DCT instances: %w", err)
				}
			case "lastNonce":
				dctData.LastNonce, err = p.processCheckUint64(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account DCT lastNonce: %w", err)
				}
			case "roles":
				dctData.Roles, err = p.processStringList(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account DCT roles: %w", err)
				}
			case "frozen":
				dctData.Frozen, err = p.processCheckUint64(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid DCT frozen flag: %w", err)
				}
			default:
				return nil, fmt.Errorf("unknown DCT data field: %s", kvp.Key)
			}
		}
	}

	if firstInstanceLoaded {
		if !p.AllowDctLegacyCheckSyntax {
			return nil, fmt.Errorf("wrong DCT check state syntax: instances in root no longer allowed")
		}
		dctData.Instances = []*mj.CheckDCTInstance{firstInstance}
	}
	dctData.Instances = append(dctData.Instances, explicitInstances...)

	return &dctData, nil
}

func (p *Parser) tryProcessCheckDCTInstanceField(kvp *oj.OJsonKeyValuePair, targetInstance *mj.CheckDCTInstance) (bool, error) {
	var err error
	switch kvp.Key {
	case "nonce":
		targetInstance.Nonce, err = p.processUint64(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid account nonce: %w", err)
		}
	case "balance":
		targetInstance.Balance, err = p.processCheckBigInt(kvp.Value, bigIntUnsignedBytes)
		if err != nil {
			return false, fmt.Errorf("invalid DCT balance: %w", err)
		}
	case "creator":
		targetInstance.Creator, err = p.parseCheckBytes(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid DCT NFT creator address: %w", err)
		}
	case "royalties":
		targetInstance.Royalties, err = p.processCheckUint64(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid DCT NFT royalties: %w", err)
		}
		if targetInstance.Royalties.Value > 10000 {
			return false, errors.New("invalid DCT NFT royalties: value exceeds maximum allowed 10000")
		}
	case "hash":
		targetInstance.Hash, err = p.parseCheckBytes(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid DCT NFT hash: %w", err)
		}
	case "uri":
		targetInstance.Uris, err = p.parseCheckValueList(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid DCT NFT URI: %w", err)
		}
	case "attributes":
		targetInstance.Attributes, err = p.parseCheckBytes(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid DCT NFT attributes: %w", err)
		}
	default:
		return false, nil
	}
	return true, nil
}

func (p *Parser) processCheckDCTInstances(dctInstancesRaw oj.OJsonObject) ([]*mj.CheckDCTInstance, error) {
	var instancesResult []*mj.CheckDCTInstance
	dctInstancesList, isList := dctInstancesRaw.(*oj.OJsonList)
	if !isList {
		return nil, errors.New("dct instances object is not a list")
	}
	for _, instanceItem := range dctInstancesList.AsList() {
		instanceAsMap, isMap := instanceItem.(*oj.OJsonMap)
		if !isMap {
			return nil, errors.New("JSON map expected as dct instances list item")
		}

		instance := mj.NewCheckDCTInstance()

		for _, kvp := range instanceAsMap.OrderedKV {
			instanceFieldLoaded, err := p.tryProcessCheckDCTInstanceField(kvp, instance)
			if err != nil {
				return nil, fmt.Errorf("invalid account DCT instance field in instances list: %w", err)
			}
			if !instanceFieldLoaded {
				return nil, fmt.Errorf("invalid account DCT instance field in instances list: `%s`", kvp.Key)
			}
		}

		instancesResult = append(instancesResult, instance)

	}

	return instancesResult, nil
}
