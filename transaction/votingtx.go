package transaction

import "github.com/ninjadotorg/constant/common"

type TxVoteDCBProposal struct {
	Tx
	VoteDCBProposalData VoteDCBProposalData
}

type TxVoteGovProposal struct {
	Tx
	VoteGovProposalData VoteGovProposalData
}

type VoteGovProposalData struct {
	GovProposalTXID *common.Hash
	AmountVoteToken uint32
}

type VoteDCBProposalData struct{
	DCBProposalTXID *common.Hash
	AmountVoteToken uint32
}

func (VoteDCBProposalData VoteDCBProposalData) Hash() *common.Hash {
	record := string(common.ToBytes(VoteDCBProposalData))
	record += string(VoteDCBProposalData.AmountVoteToken)
	hash := common.DoubleHashH([]byte(record))
	return &hash
}

func (VoteGovProposalData VoteGovProposalData) Hash() *common.Hash {
	record := string(common.ToBytes(VoteGovProposalData))
	record += string(VoteGovProposalData.AmountVoteToken)
	hash := common.DoubleHashH([]byte(record))
	return &hash
}

func (thisTx TxVoteDCBProposal) Hash() *common.Hash{
	record := string(common.ToBytes(thisTx.Tx.Hash()))
	record += string(common.ToBytes(thisTx.VoteDCBProposalData.Hash()))
	hash := common.DoubleHashH([]byte(record))
	return &hash
}

func (thisTx TxVoteGovProposal) Hash() *common.Hash{
	record := string(common.ToBytes(thisTx.Tx.Hash()))
	record += string(common.ToBytes(thisTx.VoteGovProposalData.Hash()))
	hash := common.DoubleHashH([]byte(record))
	return &hash
}

