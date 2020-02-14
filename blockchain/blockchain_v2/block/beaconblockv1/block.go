package beaconblockv1

import (
	"encoding/json"
	"fmt"

	"github.com/incognitochain/incognito-chain/blockchain"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/block/blockinterface"
	"github.com/incognitochain/incognito-chain/common"
)

type BeaconBlock struct {
	ValidationData string `json:"ValidationData"`

	Body   BeaconBody
	Header BeaconHeader
}

func (beaconBlock BeaconBlock) GetHash() *common.Hash {
	return beaconBlock.Header.GetHash()
}

func (beaconBlock BeaconBlock) GetCurrentEpoch() uint64 {
	return beaconBlock.Header.Epoch
}

func (beaconBlock BeaconBlock) GetHeight() uint64 {
	return beaconBlock.Header.Height
}

func (beaconBlock *BeaconBlock) UnmarshalJSON(data []byte) error {
	tempBeaconBlock := &struct {
		ValidationData string `json:"ValidationData"`

		Header BeaconHeader
		Body   BeaconBody
	}{}
	err := json.Unmarshal(data, &tempBeaconBlock)
	if err != nil {
		return blockchain.NewBlockChainError(blockchain.UnmashallJsonShardBlockError, err)
	}
	beaconBlock.ValidationData = tempBeaconBlock.ValidationData
	beaconBlock.Header = tempBeaconBlock.Header
	beaconBlock.Body = tempBeaconBlock.Body
	return nil
}

func (beaconBlock *BeaconBlock) AddValidationField(validationData string) error {
	beaconBlock.ValidationData = validationData
	return nil
}
func (beaconBlock BeaconBlock) GetValidationField() string {
	return beaconBlock.ValidationData
}

func (beaconBlock BeaconBlock) GetRound() int {
	return beaconBlock.Header.Round
}
func (beaconBlock BeaconBlock) GetRoundKey() string {
	return fmt.Sprint(beaconBlock.Header.Height, "_", beaconBlock.Header.Round)
}

func (beaconBlock BeaconBlock) GetInstructions() [][]string {
	return beaconBlock.Body.Instructions
}

func (beaconBlock BeaconBlock) GetProducer() string {
	return beaconBlock.Header.Producer
}

func (beaconBlock BeaconBlock) GetProducerPubKeyStr() string {
	return beaconBlock.Header.ProducerPubKeyStr
}

func (beaconBlock BeaconBlock) GetConsensusType() string {
	return beaconBlock.Header.ConsensusType
}

func (beaconBlock BeaconBlock) GetHeader() blockinterface.BeaconHeaderInterface {
	return beaconBlock.Header
}

func (beaconBlock BeaconBlock) GetBody() blockinterface.BeaconBodyInterface {
	return beaconBlock.Body
}

func (beaconBlock BeaconBlock) GetVersion() int {
	return beaconBlock.Header.Version
}