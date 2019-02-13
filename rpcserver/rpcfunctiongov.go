package rpcserver

import (
	"encoding/json"
	"errors"

	"github.com/ninjadotorg/constant/common"
	"github.com/ninjadotorg/constant/common/base58"
	"github.com/ninjadotorg/constant/metadata"
	"github.com/ninjadotorg/constant/rpcserver/jsonresult"
	"github.com/ninjadotorg/constant/transaction"
	"github.com/ninjadotorg/constant/wallet"
	"github.com/ninjadotorg/constant/wire"
)

func (rpcServer RpcServer) handleGetBondTypes(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	tempRes1 := jsonresult.GetBondTypeResultItem{
		StartSellingAt: 0,
		EndSellingAt:   500,
		Maturity:       700,
		BuyBackPrice:   110, // in constant
		BuyPrice:       105, // in constant
		TotalIssue:     1000,
		Available:      800,
	}
	tempRes2 := jsonresult.GetBondTypeResultItem{
		StartSellingAt: 0,
		EndSellingAt:   500,
		Maturity:       700,
		BuyBackPrice:   130, // in constant
		BuyPrice:       110, // in constant
		TotalIssue:     2000,
		Available:      800,
	}
	result := jsonresult.GetBondTypeResult{
		BondTypes: make(map[string]jsonresult.GetBondTypeResultItem),
	}

	result.BondTypes["fc8bbbd183f97ff6cc55a62b2ddceade8e93eed5cdf1240b42e223d589b29645"] = tempRes1

	result.BondTypes["fe7d3d124cf0309d8f1575982842b57266951a37a717a4d332a69eb176c409fa"] = tempRes2

	return result, nil
}

func (rpcServer RpcServer) handleGetCurrentSellingBondTypes(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	stabilityInfo := rpcServer.config.BlockChain.BestState.Beacon.StabilityInfo
	sellingBondsParam := stabilityInfo.GOVConstitution.GOVParams.SellingBonds
	buyPrice := uint64(0)
	bondID := sellingBondsParam.GetID()
	bondPriceFromOracle := stabilityInfo.Oracle.Bonds[bondID.String()]
	if bondPriceFromOracle == 0 {
		buyPrice = sellingBondsParam.BondPrice
	} else {
		buyPrice = bondPriceFromOracle
	}

	bondTypeRes := jsonresult.GetBondTypeResultItem{
		BondName:       sellingBondsParam.BondName,
		BondSymbol:     sellingBondsParam.BondSymbol,
		BondID:         bondID.String(),
		StartSellingAt: sellingBondsParam.StartSellingAt,
		EndSellingAt:   sellingBondsParam.StartSellingAt + sellingBondsParam.SellingWithin,
		Maturity:       sellingBondsParam.Maturity,
		BuyBackPrice:   sellingBondsParam.BuyBackPrice, // in constant
		BuyPrice:       buyPrice,                       // in constant
		TotalIssue:     sellingBondsParam.TotalIssue,
		Available:      sellingBondsParam.BondsToSell,
	}
	result := jsonresult.GetBondTypeResult{
		BondTypes: make(map[string]jsonresult.GetBondTypeResultItem),
	}
	result.BondTypes[bondID.String()] = bondTypeRes
	return result, nil
}

func (rpcServer RpcServer) handleGetGOVParams(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	// constitution := rpcServer.config.BlockChain.BestState[0].BestBlock.Header.GOVConstitution
	// govParams := constitution.GOVParams
	// results := make(map[string]interface{})
	// results["GOVParams"] = govParams
	// results["ExecuteDuration"] = constitution.ExecuteDuration
	// results["Explanation"] = constitution.Explanation
	// return results, nil
	return nil, nil
}

func (rpcServer RpcServer) handleGetGOVConstitution(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	// constitution := rpcServer.config.BlockChain.BestState[0].BestBlock.Header.GOVConstitution
	// return constitution, nil
	return nil, nil
}

func (rpcServer RpcServer) handleGetListGOVBoard(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	// res := ListPaymentAddressToListString(rpcServer.config.BlockChain.BestState[0].BestBlock.Header.GOVGovernor.BoardPaymentAddress)
	// return res, nil
	return nil, nil
}

func (rpcServer RpcServer) handleAppendListGOVBoard(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	// arrayParams := common.InterfaceSlice(params)
	// senderKey := arrayParams[0].(string)
	// paymentAddress, _ := rpcServer.GetPaymentAddressFromSenderKeyParams(senderKey)
	// rpcServer.config.BlockChain.BestState[0].BestBlock.Header.DCBGovernor.BoardPaymentAddress = append(rpcServer.config.BlockChain.BestState[0].BestBlock.Header.DCBGovernor.BoardPaymentAddress, *paymentAddress)
	// res := ListPaymentAddressToListString(rpcServer.config.BlockChain.BestState[0].BestBlock.Header.GOVGovernor.BoardPaymentAddress)
	// return res, nil
	return nil, nil
}

func (rpcServer RpcServer) handleCreateRawTxWithBuyBackRequest(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	arrayParams := common.InterfaceSlice(params)

	senderKeyParam := arrayParams[0]
	senderKey, err := wallet.Base58CheckDeserialize(senderKeyParam.(string))
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}
	senderKey.KeySet.ImportFromPrivateKey(&senderKey.KeySet.PrivateKey)
	paymentAddr := senderKey.KeySet.PaymentAddress
	tokenParamsRaw := arrayParams[4].(map[string]interface{})
	_, voutsAmount := transaction.CreateCustomTokenReceiverArray(tokenParamsRaw["TokenReceivers"])

	tokenIDStr := tokenParamsRaw["TokenID"].(string)
	tokenID, _ := common.Hash{}.NewHashFromStr(tokenIDStr)
	meta := metadata.NewBuyBackRequest(
		paymentAddr,
		uint64(voutsAmount),
		*tokenID,
		metadata.BuyBackRequestMeta,
	)
	customTokenTx, err := rpcServer.buildRawCustomTokenTransaction(params, meta)
	if err != nil {
		Logger.log.Error(err)
		return nil, NewRPCError(ErrUnexpected, err)
	}

	byteArrays, err := json.Marshal(customTokenTx)
	if err != nil {
		Logger.log.Error(err)
		return nil, NewRPCError(ErrUnexpected, err)
	}
	result := jsonresult.CreateTransactionResult{
		TxID:            customTokenTx.Hash().String(),
		Base58CheckData: base58.Base58Check{}.Encode(byteArrays, 0x00),
	}
	return result, nil
}

func (rpcServer RpcServer) handleCreateAndSendTxWithBuyBackRequest(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	data, err := rpcServer.handleCreateRawTxWithBuyBackRequest(params, closeChan)
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}

	tx := data.(jsonresult.CreateTransactionResult)
	base58CheckData := tx.Base58CheckData
	newParam := make([]interface{}, 0)
	newParam = append(newParam, base58CheckData)
	sendResult, err1 := rpcServer.handleSendRawCustomTokenTransaction(newParam, closeChan)
	if err1 != nil {
		return nil, NewRPCError(ErrUnexpected, err1)
	}
	result := jsonresult.CreateTransactionResult{
		TxID: sendResult.(jsonresult.CreateTransactionResult).TxID,
	}
	return result, nil
}

func (rpcServer RpcServer) handleCreateRawTxWithBuySellRequest(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	arrayParams := common.InterfaceSlice(params)

	// Req param #5: buy/sell request info
	buySellReq := arrayParams[4].(map[string]interface{})

	paymentAddressP := buySellReq["PaymentAddress"].(string)
	key, _ := wallet.Base58CheckDeserialize(paymentAddressP)
	tokenIDStr := buySellReq["TokenID"].(string)
	tokenID, _ := common.Hash{}.NewHashFromStr(tokenIDStr)
	amount := uint64(buySellReq["Amount"].(float64))
	buyPrice := uint64(buySellReq["BuyPrice"].(float64))
	metaType := metadata.BuyFromGOVRequestMeta
	meta := metadata.NewBuySellRequest(
		key.KeySet.PaymentAddress,
		*tokenID,
		amount,
		buyPrice,
		metaType,
	)
	normalTx, err := rpcServer.buildRawTransaction(params, meta)
	if err != nil {
		Logger.log.Error(err)
		return nil, NewRPCError(ErrUnexpected, err)
	}

	byteArrays, err1 := json.Marshal(normalTx)
	if err1 != nil {
		Logger.log.Error(err1)
		return nil, NewRPCError(ErrUnexpected, err1)
	}
	result := jsonresult.CreateTransactionResult{
		TxID:            normalTx.Hash().String(),
		Base58CheckData: base58.Base58Check{}.Encode(byteArrays, 0x00),
	}
	return result, nil
}

func (rpcServer RpcServer) handleCreateAndSendTxWithBuySellRequest(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	data, err := rpcServer.handleCreateRawTxWithBuySellRequest(params, closeChan)
	if err != nil {
		return nil, err
	}
	tx := data.(jsonresult.CreateTransactionResult)
	base58CheckData := tx.Base58CheckData
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}
	newParam := make([]interface{}, 0)
	newParam = append(newParam, base58CheckData)
	sendResult, err := rpcServer.handleSendRawTransaction(newParam, closeChan)
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}
	result := jsonresult.CreateTransactionResult{
		TxID: sendResult.(jsonresult.CreateTransactionResult).TxID,
	}
	return result, nil
}

func (rpcServer RpcServer) buildRawVoteGOVBoardTransaction(
	params interface{},
) (*transaction.TxCustomToken, error) {
	arrayParams := common.InterfaceSlice(params)
	candidatePaymentAddress := arrayParams[len(arrayParams)-1].(string)
	account, _ := wallet.Base58CheckDeserialize(candidatePaymentAddress)

	receiversPaymentAddressParam := make(map[string]interface{})
	if arrayParams[1] != nil {
		receiversPaymentAddressParam = arrayParams[1].(map[string]interface{})
	}
	amount := getAmountVote(receiversPaymentAddressParam)

	metadata := metadata.NewVoteGOVBoardMetadata(account.KeySet.PaymentAddress, amount)
	tx, err := rpcServer.buildRawCustomTokenTransaction(params, metadata)
	return tx, err
}

func (rpcServer RpcServer) handleSendRawVoteBoardGOVTransaction(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	Logger.log.Info(params)
	arrayParams := common.InterfaceSlice(params)
	base58CheckData := arrayParams[0].(string)
	rawTxBytes, _, err := base58.Base58Check{}.Decode(base58CheckData)

	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}
	tx := transaction.TxCustomToken{}
	// Logger.log.Info(string(rawTxBytes))
	err = json.Unmarshal(rawTxBytes, &tx)
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}

	hash, txDesc, err := rpcServer.config.TxMemPool.MaybeAcceptTransaction(&tx)
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}

	Logger.log.Infof("there is hash of transaction: %s\n", hash.String())
	Logger.log.Infof("there is priority of transaction in pool: %d", txDesc.StartingPriority)

	// broadcast message
	txMsg, err := wire.MakeEmptyMessage(wire.CmdCustomToken)
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}

	txMsg.(*wire.MessageTx).Transaction = &tx
	rpcServer.config.Server.PushMessageToAll(txMsg)

	return tx.Hash(), nil
}

func (rpcServer RpcServer) handleCreateRawVoteGOVBoardTransaction(
	params interface{},
	closeChan <-chan struct{},
) (interface{}, *RPCError) {
	tx, err := rpcServer.buildRawVoteGOVBoardTransaction(params)
	if err != nil {
		Logger.log.Error(err)
		return nil, NewRPCError(ErrUnexpected, err)
	}

	byteArrays, err := json.Marshal(tx)
	if err != nil {
		Logger.log.Error(err)
		return nil, NewRPCError(ErrUnexpected, err)
	}
	result := jsonresult.CreateTransactionResult{
		//HexData:         hexData,
		TxID:            tx.Hash().String(),
		Base58CheckData: base58.Base58Check{}.Encode(byteArrays, 0x00),
	}
	return result, nil
}

func (rpcServer RpcServer) handleCreateAndSendVoteGOVBoardTransaction(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	data, err := rpcServer.handleCreateRawVoteGOVBoardTransaction(params, closeChan)
	if err != nil {
		return nil, err
	}
	tx := data.(jsonresult.CreateTransactionResult)
	base58CheckData := tx.Base58CheckData
	if err != nil {
		return nil, err
	}
	newParam := make([]interface{}, 0)
	newParam = append(newParam, base58CheckData)
	txId, err := rpcServer.handleSendRawVoteBoardGOVTransaction(newParam, closeChan)
	return txId, err
}

func (rpcServer RpcServer) buildRawSubmitGOVProposalTransaction(
	params interface{},
) (*transaction.Tx, *RPCError) {
	arrayParams := common.InterfaceSlice(params)
	NParams := len(arrayParams)

	newParams := arrayParams[NParams-1].(map[string]interface{})
	tmp, err := rpcServer.GetPaymentAddressFromPrivateKeyParams(arrayParams[0].(string))
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}
	newParams["PaymentAddress"] = tmp

	meta := metadata.NewSubmitGOVProposalMetadataFromJson(newParams)
	params = setBuildRawBurnTransactionParams(params, common.FeeSubmitProposal)
	tx, err1 := rpcServer.buildRawTransaction(params, meta)
	if err1 != nil {
		return nil, err1
	}

	return tx, nil
}

func (rpcServer RpcServer) handleCreateRawSubmitGOVProposalTransaction(
	params interface{},
	closeChan <-chan struct{},
) (interface{}, *RPCError) {
	tx, err := rpcServer.buildRawSubmitGOVProposalTransaction(params)
	if err != nil {
		Logger.log.Error(err)
		return nil, NewRPCError(ErrUnexpected, err)
	}

	byteArrays, err1 := json.Marshal(tx)
	if err1 != nil {
		Logger.log.Error(err1)
		return nil, NewRPCError(ErrUnexpected, err1)
	}
	result := jsonresult.CreateTransactionResult{
		//HexData:         hexData,
		TxID:            tx.Hash().String(),
		Base58CheckData: base58.Base58Check{}.Encode(byteArrays, 0x00),
	}
	return result, nil
}

func (rpcServer RpcServer) handleSendRawSubmitGOVProposalTransaction(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	Logger.log.Info(params)
	arrayParams := common.InterfaceSlice(params)
	base58CheckData := arrayParams[0].(string)
	rawTxBytes, _, err := base58.Base58Check{}.Decode(base58CheckData)

	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}
	tx := transaction.Tx{}
	// Logger.log.Info(string(rawTxBytes))
	err = json.Unmarshal(rawTxBytes, &tx)
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}

	hash, txDesc, err1 := rpcServer.config.TxMemPool.MaybeAcceptTransaction(&tx)
	if err1 != nil {
		return nil, NewRPCError(ErrUnexpected, err1)
	}

	Logger.log.Infof("there is hash of transaction: %s\n", hash.String())
	Logger.log.Infof("there is priority of transaction in pool: %d", txDesc.StartingPriority)

	// broadcast message
	txMsg, err := wire.MakeEmptyMessage(wire.CmdTx)
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}

	txMsg.(*wire.MessageTx).Transaction = &tx
	rpcServer.config.Server.PushMessageToAll(txMsg)

	return tx.Hash(), nil
}

func (rpcServer RpcServer) handleCreateAndSendSubmitGOVProposalTransaction(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	data, err := rpcServer.handleCreateRawSubmitGOVProposalTransaction(params, closeChan)
	if err != nil {
		return nil, err
	}
	tx := data.(jsonresult.CreateTransactionResult)
	base58CheckData := tx.Base58CheckData
	if err != nil {
		return nil, err
	}
	newParam := make([]interface{}, 0)
	newParam = append(newParam, base58CheckData)
	txId, err := rpcServer.handleSendRawSubmitGOVProposalTransaction(newParam, closeChan)
	return txId, err
}

func (rpcServer RpcServer) handleCreateRawTxWithOracleFeed(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	arrayParams := common.InterfaceSlice(params)

	senderKeyParam := arrayParams[0]
	senderKey, err := wallet.Base58CheckDeserialize(senderKeyParam.(string))
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}
	senderKey.KeySet.ImportFromPrivateKey(&senderKey.KeySet.PrivateKey)
	feederAddr := senderKey.KeySet.PaymentAddress

	// Req param #4: oracle feed
	oracleFeed := arrayParams[4].(map[string]interface{})

	assetTypeBytes := []byte(oracleFeed["AssetType"].(string))
	assetType := common.Hash{}
	copy(assetType[:], assetTypeBytes)
	price := uint64(oracleFeed["Price"].(float64))
	metaType := metadata.OracleFeedMeta

	meta := metadata.NewOracleFeed(
		assetType,
		price,
		metaType,
		feederAddr,
	)

	normalTx, err := rpcServer.buildRawTransaction(params, meta)
	if err != nil {
		Logger.log.Error(err)
		return nil, NewRPCError(ErrUnexpected, err)
	}

	byteArrays, err1 := json.Marshal(normalTx)
	if err1 != nil {
		Logger.log.Error(err1)
		return nil, NewRPCError(ErrUnexpected, err1)
	}
	result := jsonresult.CreateTransactionResult{
		TxID:            normalTx.Hash().String(),
		Base58CheckData: base58.Base58Check{}.Encode(byteArrays, 0x00),
	}
	return result, nil
}

func (rpcServer RpcServer) handleCreateAndSendTxWithOracleFeed(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	data, err := rpcServer.handleCreateRawTxWithOracleFeed(params, closeChan)
	if err != nil {
		return nil, err
	}
	tx := data.(jsonresult.CreateTransactionResult)
	base58CheckData := tx.Base58CheckData
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}
	newParam := make([]interface{}, 0)
	newParam = append(newParam, base58CheckData)
	sendResult, err := rpcServer.handleSendRawTransaction(newParam, closeChan)
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}
	result := jsonresult.CreateTransactionResult{
		TxID: sendResult.(jsonresult.CreateTransactionResult).TxID,
	}
	return result, nil
}

func (rpcServer RpcServer) handleCreateRawTxWithUpdatingOracleBoard(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	arrayParams := common.InterfaceSlice(params)
	// Req param #4: updating oracle board info
	updatingOracleBoard := arrayParams[4].(map[string]interface{})
	action := int8(updatingOracleBoard["Action"].(float64))
	oraclePubKeys := updatingOracleBoard["OraclePubKeys"].([]interface{})
	assertedOraclePKs := [][]byte{}
	for _, pk := range oraclePubKeys {
		assertedOraclePKs = append(assertedOraclePKs, []byte(pk.(string)))
	}
	signs := updatingOracleBoard["Signs"].(map[string]interface{})
	assertedSigns := map[string][]byte{}
	for k, s := range signs {
		assertedSigns[k] = []byte(s.(string))
	}
	metaType := metadata.UpdatingOracleBoardMeta
	meta := metadata.NewUpdatingOracleBoard(
		action,
		assertedOraclePKs,
		assertedSigns,
		metaType,
	)

	normalTx, err := rpcServer.buildRawTransaction(params, meta)
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}

	byteArrays, err1 := json.Marshal(normalTx)
	if err1 != nil {
		return nil, NewRPCError(ErrUnexpected, err1)
	}
	result := jsonresult.CreateTransactionResult{
		TxID:            normalTx.Hash().String(),
		Base58CheckData: base58.Base58Check{}.Encode(byteArrays, 0x00),
	}
	return result, nil
}

func (rpcServer RpcServer) handleCreateAndSendTxWithUpdatingOracleBoard(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	data, err := rpcServer.handleCreateRawTxWithUpdatingOracleBoard(params, closeChan)
	if err != nil {
		return nil, err
	}
	tx := data.(jsonresult.CreateTransactionResult)
	base58CheckData := tx.Base58CheckData
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}
	newParam := make([]interface{}, 0)
	newParam = append(newParam, base58CheckData)
	sendResult, err := rpcServer.handleSendRawTransaction(newParam, closeChan)
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}
	result := jsonresult.CreateTransactionResult{
		TxID: sendResult.(jsonresult.CreateTransactionResult).TxID,
	}
	return result, nil
}

func (self RpcServer) handleCreateRawTxWithSenderAddress(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	arrayParams := common.InterfaceSlice(params)

	hasPrivacy := int(arrayParams[3].(float64)) > 0
	if hasPrivacy {
		return nil, NewRPCError(ErrUnexpected, errors.New("Could not stick sender address to metadata when enabling privacy feature."))
	}

	senderKeyParam := arrayParams[0]
	senderKey, err := wallet.Base58CheckDeserialize(senderKeyParam.(string))
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}
	senderKey.KeySet.ImportFromPrivateKey(&senderKey.KeySet.PrivateKey)
	senderAddr := senderKey.KeySet.PaymentAddress
	metaType := metadata.WithSenderAddressMeta

	meta := metadata.NewWithSenderAddress(senderAddr, metaType)
	normalTx, err := self.buildRawTransaction(params, meta)
	if err != nil {
		Logger.log.Error(err)
		return nil, NewRPCError(ErrUnexpected, err)
	}

	byteArrays, err1 := json.Marshal(normalTx)
	if err1 != nil {
		Logger.log.Error(err1)
		return nil, NewRPCError(ErrUnexpected, err1)
	}
	result := jsonresult.CreateTransactionResult{
		TxID:            normalTx.Hash().String(),
		Base58CheckData: base58.Base58Check{}.Encode(byteArrays, 0x00),
	}
	return result, nil
}

func (self RpcServer) handleCreateAndSendTxWithSenderAddress(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	data, err := self.handleCreateRawTxWithSenderAddress(params, closeChan)
	if err != nil {
		return nil, err
	}
	tx := data.(jsonresult.CreateTransactionResult)
	base58CheckData := tx.Base58CheckData
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}
	newParam := make([]interface{}, 0)
	newParam = append(newParam, base58CheckData)
	sendResult, err := self.handleSendRawTransaction(newParam, closeChan)
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}
	result := jsonresult.CreateTransactionResult{
		TxID: sendResult.(jsonresult.CreateTransactionResult).TxID,
	}
	return result, nil
}

func (rpcServer RpcServer) handleGetCurrentSellingGOVTokens(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	stabilityInfo := rpcServer.config.BlockChain.BestState.Beacon.StabilityInfo
	sellingGOVTokensParam := stabilityInfo.GOVConstitution.GOVParams.SellingGOVTokens
	buyPrice := uint64(0)
	govTokenPriceFromOracle := stabilityInfo.Oracle.GOVToken
	if govTokenPriceFromOracle == 0 {
		buyPrice = sellingGOVTokensParam.GOVTokenPrice
	} else {
		buyPrice = govTokenPriceFromOracle
	}

	result := jsonresult.GetCurrentSellingGOVTokens{
		GOVTokenID:     common.GOVTokenID.String(),
		StartSellingAt: sellingGOVTokensParam.StartSellingAt,
		EndSellingAt:   sellingGOVTokensParam.StartSellingAt + sellingGOVTokensParam.SellingWithin,
		BuyPrice:       buyPrice, // in constant
		TotalIssue:     sellingGOVTokensParam.TotalIssue,
		Available:      sellingGOVTokensParam.GOVTokensToSell,
	}
	return result, nil
}

func (rpcServer RpcServer) handleCreateRawTxWithBuyGOVTokensRequest(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	arrayParams := common.InterfaceSlice(params)

	// Req param #5: buy gov tokens request info
	buyGOVTokensReq := arrayParams[4].(map[string]interface{})

	paymentAddressP := buyGOVTokensReq["PaymentAddress"].(string)
	key, _ := wallet.Base58CheckDeserialize(paymentAddressP)
	tokenIDStr := buyGOVTokensReq["TokenID"].(string)
	tokenID, _ := common.Hash{}.NewHashFromStr(tokenIDStr)
	amount := uint64(buyGOVTokensReq["Amount"].(float64))
	buyPrice := uint64(buyGOVTokensReq["BuyPrice"].(float64))
	metaType := metadata.BuyGOVTokenRequestMeta
	meta := metadata.NewBuyGOVTokenRequest(
		key.KeySet.PaymentAddress,
		*tokenID,
		amount,
		buyPrice,
		metaType,
	)
	normalTx, err := rpcServer.buildRawTransaction(params, meta)
	if err != nil {
		Logger.log.Error(err)
		return nil, NewRPCError(ErrUnexpected, err)
	}

	byteArrays, err1 := json.Marshal(normalTx)
	if err1 != nil {
		Logger.log.Error(err1)
		return nil, NewRPCError(ErrUnexpected, err1)
	}
	result := jsonresult.CreateTransactionResult{
		TxID:            normalTx.Hash().String(),
		Base58CheckData: base58.Base58Check{}.Encode(byteArrays, 0x00),
	}
	return result, nil
}

func (rpcServer RpcServer) handleCreateAndSendTxWithBuyGOVTokensRequest(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	data, err := rpcServer.handleCreateRawTxWithBuyGOVTokensRequest(params, closeChan)
	if err != nil {
		return nil, err
	}
	tx := data.(jsonresult.CreateTransactionResult)
	base58CheckData := tx.Base58CheckData
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}
	newParam := make([]interface{}, 0)
	newParam = append(newParam, base58CheckData)
	sendResult, err := rpcServer.handleSendRawTransaction(newParam, closeChan)
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}
	result := jsonresult.CreateTransactionResult{
		TxID: sendResult.(jsonresult.CreateTransactionResult).TxID,
	}
	return result, nil
}
