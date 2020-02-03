package consensus

import (
	"context"
	"time"

	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/incognitokey"
	libp2p "github.com/libp2p/go-libp2p-peer"
)

type NodeSender interface {
	GetID() string
}
type BlockInterface interface {
	GetBlockType() string
	GetHeight() uint64
	GetBeaconHeight() uint64
	Hash() *common.Hash
	// AddValidationField(validateData string) error
	GetBlockProposer() string
	GetProducer() string
	GetValidationField() string
	GetRound() int
	GetRoundKey() string
	GetInstructions() [][]string
	GetConsensusType() string
	GetCurrentEpoch() uint64
	GetPreviousBlockHash() common.Hash
	GetTimeslot() uint64
	GetCreateTimeslot() uint64
	GetBlockTimestamp() int64
}

type BlockChainInterface interface {
	GetChain(chainName string) ChainViewManagerInterface
	GetAllChains() map[string]ChainViewManagerInterface
}

type BlockGenInterface interface {
	Start(cQuit chan struct{})
}
type NodeInterface interface {
	// RequestSyncBlockByHash(blockHash *common.Hash, isUnknownView bool, tipBlocksHash []common.Hash, peerID libp2p.ID) error
	// PushBlockToPeer(block common.BlockInterface, isShard bool, peerID libp2p.ID) error
	BroadCastBlock(blockInterface BlockInterface)
	PushMessageToChain(msg interface{}, chain ChainViewManagerInterface) error
	PushMessageToPeer(msg interface{}, peerId libp2p.ID) error
	RequestSyncBlock(nodeID string, fromView string, toView string)
	NotifyOutdatedView(nodeID string, latestView string)
	UpdateConsensusState(role string, userPbk string, currentShard *byte, beaconCommittee []string, shardCommittee map[byte][]string)
	IsEnableMining() bool
	GetMiningKeys() string
	GetPrivateKey() string
	DropAllConnections()
}

type ConsensusInterface interface {
	// NewInstance - Create a new instance of this consensus
	NewInstance(chain ChainViewManagerInterface, chainKey string, node NodeInterface, logger common.Logger) ConsensusInterface
	// GetConsensusName - retrieve consensus name
	GetConsensusName() string

	// Start - start consensus
	Start() error
	// Stop - stop consensus
	Stop() error
	// IsOngoing - check whether consensus is currently voting on a block
	// IsOngoing() bool //TO_BE_DELETE
	// ProcessBFTMsg - process incoming BFT message
	ProcessBFTMsg(msg interface{}, sender NodeSender)
	// ValidateProducerSig - validate a block producer signature
	ValidateProducerSig(block BlockInterface) error
	// ValidateCommitteeSig - validate a block committee signature
	ValidateCommitteeSig(block BlockInterface, committee []incognitokey.CommitteePublicKey) error

	// LoadUserKey - load user mining key
	LoadUserKey(miningKey string) error
	// LoadUserKeyFromIncPrivateKey - load user mining key from incognito privatekey
	LoadUserKeyFromIncPrivateKey(privateKey string) (string, error)
	// GetUserPublicKey - get user public key of loaded mining key
	GetUserPublicKey() *incognitokey.CommitteePublicKey
	// ValidateData - validate data with this consensus signature scheme
	ValidateData(data []byte, sig string, publicKey string) error
	// SignData - sign data with this consensus signature scheme
	SignData(data []byte) (string, error)
	// ExtractBridgeValidationData - extract bridge related field in validation data of block
	ExtractBridgeValidationData(block BlockInterface) ([][]byte, []int, error)
}

type BeaconManagerInterface interface {
	ChainViewManagerInterface
	GetAllCommittees() map[string]map[string][]incognitokey.CommitteePublicKey
	GetBeaconPendingList() []incognitokey.CommitteePublicKey
	GetShardsPendingList() map[string]map[string][]incognitokey.CommitteePublicKey
	GetShardsWaitingList() []incognitokey.CommitteePublicKey
	GetBeaconWaitingList() []incognitokey.CommitteePublicKey
}

type BeaconViewInterface interface {
	GetBestHeightOfShard(shardID byte) uint64
}

type ChainViewManagerInterface interface {
	GetChainName() string
	GetShardID() int
	GetGenesisTime() int64
	UnmarshalBlock(blockString []byte) (BlockInterface, error)
	GetBestView() ChainViewInterface
	GetFinalView() ChainViewInterface
	GetAllViews() map[string]ChainViewInterface
	GetViewByRange(from, to string) []ChainViewInterface
	GetViewByHash(common.Hash) (ChainViewInterface, error)
	ConnectBlockAndAddView(block BlockInterface) error
}

type ChainViewInterface interface {
	GetGenesisTime() int64
	GetConsensusConfig() string
	GetConsensusType() string
	GetBlkMinInterval() time.Duration
	GetBlkMaxCreateTime() time.Duration
	GetPubkeyRole(pubkey string, round int) (string, byte)
	GetCommittee() []incognitokey.CommitteePublicKey
	GetCommitteeHash() common.Hash
	GetCommitteeIndex(string) int
	GetBlock() BlockInterface
	GetHeight() uint64
	GetRound() int
	GetTimeStamp() int64
	GetTimeslot() uint64
	GetEpoch() uint64
	Hash() common.Hash
	GetPreviousViewHash() common.Hash

	CloneNewView() ChainViewInterface
	GetNextProposer(uint64) string
	//CreateNewViewFromBlock(BlockInterface) ChainViewInterface
	ValidateBlockAndCreateNewView(ctx context.Context, block BlockInterface, isPreSign bool) (ChainViewInterface, error)
	CreateNewBlock(context.Context, uint64, string) (BlockInterface, error)
	UnmarshalBlock(blockString []byte) (BlockInterface, error)
	GetRootTimeSlot() uint64
	CreateBlockFromOldBlockData(block BlockInterface) BlockInterface
}

type ConsensusMsgInterface interface {
	GetChainKey() string
}