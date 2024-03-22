package scenjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/Dharitri-org/drtg-scenario/model"
	oj "github.com/Dharitri-org/drtg-scenario/orderedjson"
)

func (p *Parser) processTxDCT(txDctRaw oj.OJsonObject) ([]*mj.DCTTxData, error) {
	allDctData := make([]*mj.DCTTxData, 0)

	switch txDct := txDctRaw.(type) {
	case *oj.OJsonMap:
		if !p.AllowDctTxLegacySyntax {
			return nil, fmt.Errorf("wrong DCT Multi-Transfer format, list expected")
		}
		entry, err := p.parseSingleTxDctEntry(txDct)
		if err != nil {
			return nil, err
		}

		allDctData = append(allDctData, entry)
	case *oj.OJsonList:
		for _, txDctListItem := range txDct.AsList() {
			txDctMap, isMap := txDctListItem.(*oj.OJsonMap)
			if !isMap {
				return nil, fmt.Errorf("wrong DCT Multi-Transfer format")
			}

			entry, err := p.parseSingleTxDctEntry(txDctMap)
			if err != nil {
				return nil, err
			}

			allDctData = append(allDctData, entry)
		}
	default:
		return nil, fmt.Errorf("wrong DCT transfer format, expected list")
	}

	return allDctData, nil
}

func (p *Parser) parseSingleTxDctEntry(dctTxEntry *oj.OJsonMap) (*mj.DCTTxData, error) {
	dctData := mj.DCTTxData{}
	var err error

	for _, kvp := range dctTxEntry.OrderedKV {
		switch kvp.Key {
		case "tokenIdentifier":
			dctData.TokenIdentifier, err = p.processStringAsByteArray(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid DCT token name: %w", err)
			}
		case "nonce":
			dctData.Nonce, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, errors.New("invalid account nonce")
			}
		case "value":
			dctData.Value, err = p.processBigInt(kvp.Value, bigIntUnsignedBytes)
			if err != nil {
				return nil, fmt.Errorf("invalid DCT balance: %w", err)
			}
		default:
			return nil, fmt.Errorf("unknown transaction DCT data field: %s", kvp.Key)
		}
	}

	return &dctData, nil
}
