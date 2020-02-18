package blockchain

import (
	"encoding/json"
	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/dataaccessobject/rawdb"
	"github.com/incognitochain/incognito-chain/dataaccessobject/statedb"
)

// ListCustomToken - return all custom token which existed in network
func (blockchain *BlockChain) ListPrivacyCustomTokenV2(shardID byte) (map[common.Hash]*statedb.TokenState, error) {
	tokenStates, err := statedb.ListPrivacyToken(blockchain.BestState.Shard[shardID].GetCopiedTransactionStateDB())
	if err != nil {
		return nil, err
	}
	delete(tokenStates, common.PRVCoinID)
	return tokenStates, nil
}
func (blockchain *BlockChain) GetAllCoinIDV2(shardID byte) ([]common.Hash, error) {
	tokenIDs := []common.Hash{}
	tokenStates, err := blockchain.ListPrivacyCustomTokenV2(shardID)
	if err != nil {
		return nil, err
	}
	for k, _ := range tokenStates {
		tokenIDs = append(tokenIDs, k)
	}
	brigdeTokenIDs, _, err := blockchain.GetAllBridgeTokens()
	if err != nil {
		return nil, err
	}
	for _, bridgeTokenID := range brigdeTokenIDs {
		if _, found := tokenStates[bridgeTokenID]; !found {
			tokenIDs = append(tokenIDs, bridgeTokenID)
		}
	}
	return tokenIDs, nil
}

// Check Privacy Custom token ID is existed
func (blockchain *BlockChain) PrivacyCustomTokenIDExistedV2(tokenID *common.Hash, shardID byte) bool {
	return statedb.PrivacyTokenIDExisted(blockchain.BestState.Shard[shardID].GetCopiedTransactionStateDB(), *tokenID)
}

func (blockchain *BlockChain) GetAllBridgeTokens() ([]common.Hash, []*rawdb.BridgeTokenInfo, error) {
	bridgeTokenIDs := []common.Hash{}
	allBridgeTokens := []*rawdb.BridgeTokenInfo{}
	bridgeStateDB := blockchain.BestState.Beacon.GetCopiedFeatureStateDB()
	allBridgeTokensBytes, err := statedb.GetAllBridgeTokens(bridgeStateDB)
	if err != nil {
		return bridgeTokenIDs, allBridgeTokens, err
	}
	err = json.Unmarshal(allBridgeTokensBytes, &allBridgeTokens)
	if err != nil {
		return bridgeTokenIDs, allBridgeTokens, err
	}
	for _, bridgeTokenInfo := range allBridgeTokens {
		bridgeTokenIDs = append(bridgeTokenIDs, *bridgeTokenInfo.TokenID)
	}
	return bridgeTokenIDs, allBridgeTokens, nil
}