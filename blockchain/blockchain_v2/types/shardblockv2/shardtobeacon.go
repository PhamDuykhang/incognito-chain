package shardblockv2

import (
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/consensusheader"
	"github.com/incognitochain/incognito-chain/common"
)

type ShardToBeaconBlock struct {
	ConsensusHeader consensusheader.ConsensusHeader

	Instructions [][]string
	Header       ShardHeader
}

func (shardToBeaconBlock ShardToBeaconBlock) GetEpoch() uint64 {
	return shardToBeaconBlock.Header.Epoch
}

func (block ShardToBeaconBlock) GetValidationField() string {
	return block.ConsensusHeader.ValidationData
}

func (block ShardToBeaconBlock) GetHeight() uint64 {
	return block.Header.Height
}
func (block ShardToBeaconBlock) GetInstructions() [][]string {
	return block.Instructions
}

func (block ShardToBeaconBlock) GetProducer() string {
	return block.Header.Producer
}
func (block ShardToBeaconBlock) GetConsensusType() string {
	return block.Header.ConsensusType
}

func (block ShardToBeaconBlock) GetPreviousBlockHash() common.Hash {
	return block.Header.PreviousBlockHash
}

func (block ShardToBeaconBlock) GetBlockType() string {
	return "shardtobeacon"
}

func (block ShardToBeaconBlock) GetShardHeader() blockinterface.ShardHeaderInterface {
	return block.Header
}
