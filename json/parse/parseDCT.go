package scenjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/Dharitri-org/drtg-scenario/model"
	oj "github.com/Dharitri-org/drtg-scenario/orderedjson"
)

func (p *Parser) processDCTData(
	tokenName mj.JSONBytesFromString,
	dctDataRaw oj.OJsonObject) (*mj.DCTData, error) {

	switch data := dctDataRaw.(type) {
	case *oj.OJsonString:
		// simple string representing balance "400,000,000,000"
		dctData := mj.DCTData{
			TokenIdentifier: tokenName,
		}
		balance, err := p.processBigInt(dctDataRaw, bigIntUnsignedBytes)
		if err != nil {
			return nil, fmt.Errorf("invalid DCT balance: %w", err)
		}
		dctData.Instances = []*mj.DCTInstance{
			{
				Nonce:   mj.JSONUint64{Value: 0, Original: ""},
				Balance: balance,
			},
		}
		return &dctData, nil
	case *oj.OJsonMap:
		return p.processDCTDataMap(tokenName, data)
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
func (p *Parser) processDCTDataMap(tokenName mj.JSONBytesFromString, dctDataMap *oj.OJsonMap) (*mj.DCTData, error) {
	dctData := mj.DCTData{
		TokenIdentifier: tokenName,
	}
	firstInstance := &mj.DCTInstance{}
	firstInstanceLoaded := false
	var explicitInstances []*mj.DCTInstance

	for _, kvp := range dctDataMap.OrderedKV {
		// it is allowed to load the instance directly, fields set to the first instance
		instanceFieldLoaded, err := p.tryProcessDCTInstanceField(kvp, firstInstance)
		if err != nil {
			return nil, fmt.Errorf("invalid account DCT instance field: %w", err)
		}
		if instanceFieldLoaded {
			firstInstanceLoaded = true
		} else {
			switch kvp.Key {
			case "instances":
				explicitInstances, err = p.processDCTInstances(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account DCT instances: %w", err)
				}
			case "lastNonce":
				dctData.LastNonce, err = p.processUint64(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account DCT lastNonce: %w", err)
				}
			case "roles":
				dctData.Roles, err = p.processStringList(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account DCT roles: %w", err)
				}
			case "frozen":
				dctData.Frozen, err = p.processUint64(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid DCT frozen flag: %w", err)
				}
			default:
				return nil, fmt.Errorf("unknown DCT data field: %s", kvp.Key)
			}
		}
	}

	if firstInstanceLoaded {
		if !p.AllowDctLegacySetSyntax {
			return nil, fmt.Errorf("wrong DCT set state syntax: instances in root no longer allowed")
		}
		dctData.Instances = []*mj.DCTInstance{firstInstance}
	}
	dctData.Instances = append(dctData.Instances, explicitInstances...)

	return &dctData, nil
}

func (p *Parser) tryProcessDCTInstanceField(kvp *oj.OJsonKeyValuePair, targetInstance *mj.DCTInstance) (bool, error) {
	var err error
	switch kvp.Key {
	case "nonce":
		targetInstance.Nonce, err = p.processUint64(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid account nonce: %w", err)
		}
	case "balance":
		targetInstance.Balance, err = p.processBigInt(kvp.Value, bigIntUnsignedBytes)
		if err != nil {
			return false, fmt.Errorf("invalid DCT balance: %w", err)
		}
	case "creator":
		targetInstance.Creator, err = p.processStringAsByteArray(kvp.Value)
		if err != nil || len(targetInstance.Creator.Value) != 32 {
			return false, fmt.Errorf("invalid DCT NFT creator address: %w", err)
		}
	case "royalties":
		targetInstance.Royalties, err = p.processUint64(kvp.Value)
		if err != nil || targetInstance.Royalties.Value > 10000 {
			return false, fmt.Errorf("invalid DCT NFT royalties: %w", err)
		}
	case "hash":
		targetInstance.Hash, err = p.processStringAsByteArray(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid DCT NFT hash: %w", err)
		}
	case "uri":
		targetInstance.Uris, err = p.parseValueList(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid DCT NFT URI: %w", err)
		}
	case "attributes":
		targetInstance.Attributes, err = p.processSubTreeAsByteArray(kvp.Value)
		if err != nil {
			return false, fmt.Errorf("invalid DCT NFT attributes: %w", err)
		}
	default:
		return false, nil
	}
	return true, nil
}

func (p *Parser) processDCTInstances(dctInstancesRaw oj.OJsonObject) ([]*mj.DCTInstance, error) {
	var instancesResult []*mj.DCTInstance
	dctInstancesList, isList := dctInstancesRaw.(*oj.OJsonList)
	if !isList {
		return nil, errors.New("dct instances object is not a list")
	}
	for _, instanceItem := range dctInstancesList.AsList() {
		instanceAsMap, isMap := instanceItem.(*oj.OJsonMap)
		if !isMap {
			return nil, errors.New("JSON map expected as dct instances list item")
		}

		instance := &mj.DCTInstance{}

		for _, kvp := range instanceAsMap.OrderedKV {
			instanceFieldLoaded, err := p.tryProcessDCTInstanceField(kvp, instance)
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
