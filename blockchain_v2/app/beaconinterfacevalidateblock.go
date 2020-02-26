package app

import (
	"context"
	"fmt"

	"github.com/incognitochain/incognito-chain/blockchain_v2/types/beaconblockv2"
	"github.com/incognitochain/incognito-chain/blockchain_v2/types/blockinterface"
	consensus "github.com/incognitochain/incognito-chain/consensus_v2"
	"github.com/incognitochain/incognito-chain/consensus_v2/blsbftv2"
)

type ValidateBeaconBlockState struct {
	ctx     context.Context
	bc      *blockchainV2
	curView *BeaconView
	//app
	app []BeaconApp

	//tmp
	isPreSign bool
	newView   *BeaconView
}

func (beaconView *BeaconView) NewValidateState(ctx context.Context, createState *CreateBeaconBlockState) *ValidateBeaconBlockState {
	validateState := &ValidateBeaconBlockState{
		ctx:     ctx,
		bc:      beaconView.bc,
		curView: beaconView,
		newView: beaconView.CloneNewView().(*BeaconView),
		app:     []BeaconApp{},
	}

	//ADD YOUR APP HERE
	validateState.app = append(validateState.app, &BeaconCoreApp{Logger: beaconView.Logger, ValidateState: validateState, CreateState: createState})
	validateState.app = append(validateState.app, &BeaconBridgeApp{Logger: beaconView.Logger, validateState: validateState, createState: createState})
	validateState.app = append(validateState.app, &BeaconPDEApp{Logger: beaconView.Logger, ValidateState: validateState, CreateState: createState})

	createState.app = validateState.app
	createState.ctx = validateState.ctx
	createState.bc = validateState.bc
	createState.curView = validateState.curView
	createState.newView = validateState.newView

	return validateState
}

func (beaconView *BeaconView) ValidateBlockAndCreateNewView(ctx context.Context, block blockinterface.BlockInterface, isPreSign bool) (consensus.ChainViewInterface, error) {
	if block.GetHeader().GetVersion() == BEACON_BLOCK_VERSION_2 {

	}
	createState := &CreateBeaconBlockState{}
	validateState := beaconView.NewValidateState(ctx, createState)
	validateState.newView.Block = block.(blockinterface.BeaconBlockInterface)
	validateState.isPreSign = isPreSign

	//block has correct basic header
	//we have enough data to validate this block and get beaconblocks, crossshardblock, txToAdd confirm by proposed block
	for _, app := range validateState.app {
		err := app.preValidate()
		if err != nil {
			return nil, err
		}
	}

	if isPreSign {

		createState.createTimeStamp = validateState.newView.Block.GetHeader().GetTimestamp()
		createState.createTimeSlot = validateState.newView.Block.GetHeader().(blockinterface.BlockHeaderV2Interface).GetTimeslot()
		createState.proposer = validateState.newView.Block.GetHeader().GetProducer()

		for _, app := range validateState.app {
			if err := app.buildInstructionFromShardAction(); err != nil {
				return nil, err
			}
		}

		for _, app := range validateState.app {
			if err := app.buildInstructionByEpoch(); err != nil {
				return nil, err
			}
		}

		instructions := [][]string{}
		// instructions = append(instructions, createState.randomInstruction)
		instructions = append(instructions, createState.rewardInstByEpoch...)
		instructions = append(instructions, createState.validStakeInstructions...)
		instructions = append(instructions, createState.validStopAutoStakingInstructions...)
		instructions = append(instructions, createState.acceptedRewardInstructions...)
		instructions = append(instructions, createState.beaconSwapInstruction...)
		instructions = append(instructions, createState.shardAssignInst...)

		instructions = append(instructions, createState.bridgeInstructions...)
		instructions = append(instructions, createState.statefulInstructions...)

		createState.newBlock = &beaconblockv2.BeaconBlock{
			Body: beaconblockv2.BeaconBody{
				ShardState:   createState.shardStates,
				Instructions: instructions,
			},
		}

		//build header
		for _, app := range validateState.app {
			if err := app.buildHeader(); err != nil {
				return nil, err
			}
		}

		//compare block hash
		if !createState.newBlock.GetHeader().GetHash().IsEqual(validateState.newView.Block.GetHeader().GetHash()) {
			fmt.Println(createState.newBlock.GetHeader().GetHash().String(), validateState.newView.Block.GetHeader().GetHash().String())
			panic(1)
			return nil, nil
		}
	} else {
		//validate producer signature
		if err := (blsbftv2.BLSBFT{}.ValidateProducerSig(block)); err != nil {
			return nil, err
		}
		//validate committeee signature
		if err := (blsbftv2.BLSBFT{}.ValidateCommitteeSig(block, beaconView.GetCommittee())); err != nil {
			beaconView.Logger.Error("Validate fail!", beaconView.GetCommittee())
			panic(1)
			return nil, err
		}

		// validateState.newView = beaconView.CloneNewView().(*BeaconView)
		// for _, app := range validateState.app {
		// 	if err := app.updateNewViewFromBlock(block.(*BeaconBlock)); err != nil {
		// 		return nil, err
		// 	}
		// }
		// validateState.newView = validateState.newView
		//TODO: compare header content, with newview, necessary???
	}

	return validateState.newView, nil
}
