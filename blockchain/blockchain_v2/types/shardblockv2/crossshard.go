package shardblockv2

import (
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"

	"github.com/incognitochain/incognito-chain/blockchain"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/consensusheader"
	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/privacy"
)

type CrossShardBlock struct {
	ConsensusHeader consensusheader.ConsensusHeader

	Header          ShardHeader
	ToShardID       byte
	MerklePathShard []common.Hash
	// Cross Shard data for PRV
	CrossOutputCoin []privacy.OutputCoin
	// Cross Shard For Custom token privacy
	CrossTxTokenPrivacyData []blockchain.ContentCrossShardTokenPrivacyData
}

func (crossShardBlock CrossShardBlock) GetEpoch() uint64 {
	return crossShardBlock.Header.Epoch
}

func (crossShardBlock *CrossShardBlock) GetHash() *common.Hash {
	return crossShardBlock.Header.GetHash()
}
func (block CrossShardBlock) GetProducer() string {
	return block.Header.Producer
}

func (block CrossShardBlock) GetHeight() uint64 {
	return block.Header.Height
}

func (block CrossShardBlock) GetValidationField() string {
	return block.ConsensusHeader.ValidationData
}
func (block CrossShardBlock) GetInstructions() [][]string {
	return [][]string{}
}

func (block CrossShardBlock) GetConsensusType() string {
	return block.Header.ConsensusType
}

func (block CrossShardBlock) GetPreviousBlockHash() common.Hash {
	return block.Header.PreviousBlockHash
}

func (block CrossShardBlock) GetCrossOutputCoin() []privacy.OutputCoin {
	return block.CrossOutputCoin
}

func (block CrossShardBlock) GetCrossTxTokenPrivacyData() []blockchain.ContentCrossShardTokenPrivacyData {
	return block.CrossTxTokenPrivacyData
}

func (block CrossShardBlock) GetMerklePathShard() []common.Hash {
	return block.MerklePathShard
}

func (block CrossShardBlock) GetShardHeader() blockinterface.ShardHeaderInterface {
	return block.Header
}

func (block CrossShardBlock) GetBlockType() string {
	return "crossshard"
}

func (block CrossShardBlock) GetToShardID() byte {
	return block.ToShardID
}
