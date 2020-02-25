package app

import (
	"context"
	"github.com/incognitochain/incognito-chain/blockchain_v2/types/blockinterface"
	"github.com/incognitochain/incognito-chain/blockchain_v2/types/shardblockv2"
	"math/rand"
	"time"

	"github.com/incognitochain/incognito-chain/blockchain"
	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/metadata"
	"github.com/incognitochain/incognito-chain/privacy"
)

type CreateShardBlockState struct {
	ctx      context.Context
	bc       BlockChain
	curView  *ShardView
	newBlock blockinterface.ShardBlockInterface
	newView  *ShardView

	//tmp
	newConfirmBeaconHeight uint64
	newConfirmBeaconHash   common.Hash
	totalTxFee             map[common.Hash]uint64
	beaconBlocks           []common.BlockInterface
	crossShardBlocks       map[byte][]blockinterface.CrossShardBlockInterface
	createTimeStamp        int64
	createTimeSlot         uint64
	proposer               string
	newBlockEpoch          uint64
	//app
	app []ShardApp

	crossShardTx             map[byte][]blockchain.CrossTransaction
	txToRemoveFromPool       []metadata.Transaction
	txsToAddFromPool         []metadata.Transaction
	txsFromMetadataTx        []metadata.Transaction
	txsFromBeaconInstruction []metadata.Transaction
	errInstruction           [][]string
	stakingTx                map[string]string
	newShardPendingValidator []string

	instruction [][]string
}

func (s *ShardView) NewCreateState(ctx context.Context) *CreateShardBlockState {
	createState := &CreateShardBlockState{
		bc:       s.BC,
		curView:  s,
		newView:  s.CloneNewView().(*ShardView),
		ctx:      ctx,
		app:      []ShardApp{},
		newBlock: nil,
	}

	//ADD YOUR APP HERE
	createState.app = append(createState.app, &ShardCoreApp{Logger: s.Logger, CreateState: createState})
	createState.app = append(createState.app, &ShardBridgeApp{Logger: s.Logger, CreateState: createState})
	createState.app = append(createState.app, &ShardPDEApp{Logger: s.Logger, CreateState: createState})

	return createState
}

func (s *ShardView) CreateNewBlock(ctx context.Context, timeslot uint64, proposer string) (blockinterface.BlockInterface, error) {
	s.Logger.Criticalf("Creating Shard Block %+v at timeslot %v", s.GetHeight()+1, timeslot)
	createState := s.NewCreateState(ctx)
	createState.createTimeStamp = time.Now().Unix()
	createState.createTimeSlot = timeslot
	createState.proposer = proposer
	//pre processing
	for _, app := range createState.app {
		if err := preCreateBlock(); err != nil {
			return nil, err
		}
	}

	//build shardbody
	for _, app := range createState.app {
		if err := buildTxFromCrossShard(); err != nil {
			return nil, err
		}
	}

	for _, app := range createState.app {
		if err := buildTxFromMemPool(); err != nil {
			return nil, err
		}
	}

	for _, app := range createState.app {
		if err := buildResponseTxFromTxWithMetadata(); err != nil {
			return nil, err
		}
	}

	for _, app := range createState.app {
		if err := processBeaconInstruction(); err != nil {
			return nil, err
		}
	}

	for _, app := range createState.app {
		if err := generateInstruction(); err != nil {
			return nil, err
		}
	}

	createState.newBlock = &shardblockv2.ShardBlock{
		Body: shardblockv2.ShardBody{
			Transactions:      append(append(createState.txsToAddFromPool, createState.txsFromMetadataTx...), createState.txsFromBeaconInstruction...),
			CrossTransactions: createState.crossShardTx,
			Instructions:      createState.instruction,
		},
	}

	for _, app := range createState.app {
		if err := updateNewViewFromBlock(createState.newBlock); err != nil {
			return nil, err
		}
	}

	//build shard header
	for _, app := range createState.app {
		if err := buildHeader(); err != nil {
			return nil, err
		}
	}

	return createState.newBlock, nil
}

func createTempKeyset() privacy.PrivateKey {
	rand.Seed(time.Now().UnixNano())
	seed := make([]byte, 16)
	rand.Read(seed)
	return privacy.GeneratePrivateKey(seed)
}
