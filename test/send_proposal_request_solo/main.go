/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"chainmaker.org/chainmaker-go/accesscontrol"
	"chainmaker.org/chainmaker-go/common/ca"
	"chainmaker.org/chainmaker-go/common/crypto"
	"chainmaker.org/chainmaker-go/common/crypto/asym"
	"chainmaker.org/chainmaker-go/common/helper"
	acPb "chainmaker.org/chainmaker-go/pb/protogo/accesscontrol"
	apiPb "chainmaker.org/chainmaker-go/pb/protogo/api"
	commonPb "chainmaker.org/chainmaker-go/pb/protogo/common"
	discoveryPb "chainmaker.org/chainmaker-go/pb/protogo/discovery"
	"chainmaker.org/chainmaker-go/protocol"
	"chainmaker.org/chainmaker-go/utils"
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	logTempMarshalPayLoadFailed     = "marshal payload failed, %s"
	logTempUnmarshalBlockInfoFailed = "blockInfo unmarshal error %s\n"
	logTempSendTx                   = "send tx resp: code:%d, msg:%s, payload:%+v\n"
	logTempSendBlock                = "send tx resp: code:%d, msg:%s, blockInfo:%+v\n"
	fieldWithRWSet                  = "withRWSet"
)

const (
	CHAIN1         = "chain1"
	IP             = "localhost"
	Port           = 12301
	certPathPrefix = "../config"
	//WasmPath       = "wasm/evidence.wasm"
	WasmPath = "wasm/fact-rust-0.7.2.wasm"
	//userKeyPath    = certPathPrefix + "/certs/wx-org1/user/user.key"
	//userCrtPath    = certPathPrefix + "/certs/wx-org1/user/user.crt"
	userKeyPath  = certPathPrefix + "/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.key"
	userCrtPath  = certPathPrefix + "/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.crt"
	orgId        = "wx-org1.chainmaker.org"
	contractName = "contract13"
	runtimeType  = commonPb.RuntimeType_WASMER
	prePathFmt   = certPathPrefix + "/crypto-config/wx-org%s.chainmaker.org/user/admin1/"
)

//var caPaths = []string{certPathPrefix + "/certs/wx-org1/ca"}
var caPaths = []string{certPathPrefix + "/crypto-config/wx-org1.chainmaker.org/ca"}

func main() {
	createContract := flag.Bool("c", true, "create Contract")
	//createContract := flag.Bool("c", false, "create Contract")
	flag.Parse()

	// conn, err := initGRPCConnect(false)
	conn, err := initGRPCConnect(true)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	client := apiPb.NewRpcNodeClient(conn)

	file, err := ioutil.ReadFile(userKeyPath)
	if err != nil {
		panic(err)
	}

	sk3, err := asym.PrivateKeyFromPEM(file, nil)
	if err != nil {
		panic(err)
	}

	// 1) 合约创建
	if *createContract {
		testCreate(sk3, &client, CHAIN1)
		time.Sleep(10 * time.Second)
		testInvoke(sk3, &client, CHAIN1)
		time.Sleep(5 * time.Second)
		testQuery2(sk3, &client, CHAIN1)
	}
	//testInvoke(sk3, &client, CHAIN1)
	// 冻结、解冻、吊销用户合约功能测试
	//testFreezeOrUnfreezeOrRevokeFlow(sk3, client)

	//////
	//////// 2)
	//txId := testInvoke(sk3, &client, CHAIN1)
	//time.Sleep(5 * time.Second)
	////////
	//////// 3) 合约查询
	//testQuery(sk3, &client, CHAIN1)
	////////
	////////// 4) 根据TxId查交易
	//testGetTxByTxId(sk3, &client, txId, CHAIN1)
	//
	// 5) 根据区块高度查区块，若height为-1，表示查当前区块
	//hash := testGetBlockByHeight(sk3, &client, CHAIN1, -1)
	//
	//// 6) 根据区块高度查区块（包含读写集），若height为-1，表示查当前区块
	//testGetBlockWithTxRWSetsByHeight(sk3, &client, CHAIN1, -1)
	//
	//// 7) 根据区块哈希查区块
	//testGetBlockByHash(sk3, &client, CHAIN1, hash)
	//
	//// 8) 根据区块哈希查区块（包含读写集）
	//testGetBlockWithTxRWSetsByHash(sk3, &client, CHAIN1, hash)
	//
	//// 9) 根据TxId查区块
	//testGetBlockByTxId(sk3, &client, txId, CHAIN1)
	//
	//// 10) 查询最新配置块
	//testGetLastConfigBlock(sk3, &client, CHAIN1)
	//
	//// 11) 查询最新区块
	//testGetLastBlock(sk3, &client, CHAIN1)
	//
	//// 12) 查询链信息
	//testGetChainInfo(sk3, &client, CHAIN1)
	//
	//// 12) 批量调用
	if false {
		var wg sync.WaitGroup
		for i := 0; i < 1; i++ {
			wg.Add(1)
			go func() {
				for j := 0; j < 10; j++ {
					var txId string
					txId = testInvoke(sk3, &client, CHAIN1)
					fmt.Printf("txId: %s\n", txId)
					time.Sleep(time.Millisecond * 100)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
	//
	////testCreate(sk3, &client, CHAIN1)
	////time.Sleep(5 * time.Second)
	//
	////// 5) 合约调用
	//for i := 0; i < 100; i++ {
	//	testInvoke(sk3, &client, CHAIN1)
	//}
	//
	////// 6) 合约查询
	////testQuery(sk3, &client, CHAIN1)
	//
	////7) 合约升级
	//testUpgrade(sk3, &client, CHAIN1)
	// 性能测试
	//testPerformanceModeTransfer(sk3, &client, CHAIN1)
	//time.Sleep(time.Second * 10)
	//testPerformanceModeBalance(sk3, &client, CHAIN1)
}

func testPerformanceModeTransfer(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, chainId string) {
	for j := 0; j < 5000; j++ {
		i := j % 1000
		txId := utils.GetRandTxId()
		// 构造Payload
		pairs := []*commonPb.KeyValuePair{
			{
				Key:   "from",
				Value: strconv.Itoa(i),
			},
			{
				Key:   "to",
				Value: strconv.Itoa(i + 1000),
			},
			{
				Key:   "amount",
				Value: "1",
			},
		}

		payload := &commonPb.TransactPayload{
			ContractName: contractName,
			Method:       "transfer",
			Parameters:   pairs,
		}

		payloadBytes, err := proto.Marshal(payload)
		if err != nil {
			log.Fatalf(logTempMarshalPayLoadFailed, err.Error())
		}

		resp := proposalRequest(sk3, client, commonPb.TxType_INVOKE_USER_CONTRACT,
			chainId, txId, payloadBytes)

		fmt.Printf(logTempSendTx, resp.Code, resp.Message, resp.ContractResult)
	}
}

func testPerformanceModeBalance(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, chainId string) {
	for i := 0; i < 2000; i++ {
		txId := utils.GetRandTxId()
		// 构造Payload
		pairs := []*commonPb.KeyValuePair{
			{
				Key:   "from",
				Value: strconv.Itoa(i),
			},
		}

		payload := &commonPb.TransactPayload{
			ContractName: contractName,
			Method:       "balance",
			Parameters:   pairs,
		}

		payloadBytes, err := proto.Marshal(payload)
		if err != nil {
			log.Fatalf(logTempMarshalPayLoadFailed, err.Error())
		}

		resp := proposalRequest(sk3, client, commonPb.TxType_QUERY_USER_CONTRACT,
			chainId, txId, payloadBytes)

		fmt.Printf(logTempSendTx, resp.Code, resp.Message, resp.ContractResult)
	}
}

func testFreezeOrUnfreezeOrRevokeFlow(sk3 crypto.PrivateKey, client apiPb.RpcNodeClient) {
	//执行合约
	testInvoke(sk3, &client, CHAIN1)
	testQuery2(sk3, &client, CHAIN1)
	time.Sleep(5 * time.Second)

	// 冻结
	testFreezeOrUnfreezeOrRevoke(sk3, &client, CHAIN1, commonPb.ManageUserContractFunction_FREEZE_CONTRACT.String())
	time.Sleep(5 * time.Second)
	testInvoke(sk3, &client, CHAIN1)
	testQuery2(sk3, &client, CHAIN1)
	time.Sleep(5 * time.Second)

	// 解冻
	testFreezeOrUnfreezeOrRevoke(sk3, &client, CHAIN1, commonPb.ManageUserContractFunction_UNFREEZE_CONTRACT.String())
	time.Sleep(5 * time.Second)
	testInvoke(sk3, &client, CHAIN1)
	testQuery2(sk3, &client, CHAIN1)
	time.Sleep(5 * time.Second)

	// 冻结
	testFreezeOrUnfreezeOrRevoke(sk3, &client, CHAIN1, commonPb.ManageUserContractFunction_FREEZE_CONTRACT.String())
	time.Sleep(5 * time.Second)
	testInvoke(sk3, &client, CHAIN1)
	testQuery2(sk3, &client, CHAIN1)
	time.Sleep(5 * time.Second)

	// 解冻
	testFreezeOrUnfreezeOrRevoke(sk3, &client, CHAIN1, commonPb.ManageUserContractFunction_UNFREEZE_CONTRACT.String())
	time.Sleep(5 * time.Second)
	testInvoke(sk3, &client, CHAIN1)
	testQuery2(sk3, &client, CHAIN1)
	time.Sleep(5 * time.Second)

	// 冻结
	//testFreezeOrUnfreezeOrRevoke(sk3, &client, CHAIN1, commonPb.ManageUserContractFunction_FREEZE_CONTRACT.String())
	//time.Sleep(5 * time.Second)
	// 吊销
	testFreezeOrUnfreezeOrRevoke(sk3, &client, CHAIN1, commonPb.ManageUserContractFunction_REVOKE_CONTRACT.String())
	time.Sleep(5 * time.Second)
	testInvoke(sk3, &client, CHAIN1)
	testQuery2(sk3, &client, CHAIN1)
	time.Sleep(5 * time.Second)

	testFreezeOrUnfreezeOrRevoke(sk3, &client, CHAIN1, commonPb.ManageUserContractFunction_FREEZE_CONTRACT.String())
	time.Sleep(5 * time.Second)

	testFreezeOrUnfreezeOrRevoke(sk3, &client, CHAIN1, commonPb.ManageUserContractFunction_UNFREEZE_CONTRACT.String())
	time.Sleep(5 * time.Second)
}

func testGetTxByTxId(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, txId, chainId string) {
	fmt.Printf("\n============ get tx by txId [%s] ============\n", txId)

	// 构造Payload
	pair := &commonPb.KeyValuePair{Key: "txId", Value: txId}
	var pairs []*commonPb.KeyValuePair
	pairs = append(pairs, pair)

	payloadBytes := constructPayload(commonPb.ContractName_SYSTEM_CONTRACT_QUERY.String(), "GET_TX_BY_TX_ID", pairs)

	resp := proposalRequest(sk3, client, commonPb.TxType_QUERY_SYSTEM_CONTRACT,
		chainId, txId, payloadBytes)

	fmt.Printf(logTempSendTx, resp.Code, resp.Message, resp.ContractResult)
}

func testGetBlockByTxId(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, txId, chainId string) {
	fmt.Printf("\n============ get block by txId [%s] ============\n", txId)

	// 构造Payload
	pairs := []*commonPb.KeyValuePair{
		{
			Key:   "txId",
			Value: txId,
		},
		{
			Key:   fieldWithRWSet,
			Value: "false",
		},
	}

	payloadBytes := constructPayload(commonPb.ContractName_SYSTEM_CONTRACT_QUERY.String(), "GET_BLOCK_BY_TX_ID", pairs)

	resp := proposalRequest(sk3, client, commonPb.TxType_QUERY_SYSTEM_CONTRACT,
		chainId, txId, payloadBytes)

	blockInfo := &commonPb.BlockInfo{}
	err := proto.Unmarshal(resp.ContractResult.Result, blockInfo)
	if err != nil {
		fmt.Printf(logTempUnmarshalBlockInfoFailed, err)
		os.Exit(0)
	}
	fmt.Printf(logTempSendBlock, resp.ContractResult.Code, resp.ContractResult.Message, blockInfo)
}

func testGetBlockByHeight(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, chainId string, height int64) string {
	fmt.Printf("\n============ get block by height [%d] ============\n", height)
	// 构造Payload

	pairs := []*commonPb.KeyValuePair{
		{
			Key:   "blockHeight",
			Value: strconv.FormatInt(height, 10),
		},
		{
			Key:   fieldWithRWSet,
			Value: "false",
		},
	}

	payloadBytes := constructPayload(commonPb.ContractName_SYSTEM_CONTRACT_QUERY.String(), "GET_BLOCK_BY_HEIGHT", pairs)

	resp := proposalRequest(sk3, client, commonPb.TxType_QUERY_SYSTEM_CONTRACT,
		chainId, "", payloadBytes)

	blockInfo := &commonPb.BlockInfo{}
	err := proto.Unmarshal(resp.ContractResult.Result, blockInfo)
	if err != nil {
		fmt.Printf(logTempUnmarshalBlockInfoFailed, err)
		os.Exit(0)
	}
	fmt.Printf(logTempSendBlock, resp.ContractResult.Code, resp.ContractResult.Message, blockInfo)

	return hex.EncodeToString(blockInfo.Block.Header.BlockHash)
}

func testGetBlockWithTxRWSetsByHeight(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, chainId string, height int64) string {
	fmt.Printf("\n============ get block with txRWsets by height [%d] ============\n", height)
	// 构造Payload
	pairs := []*commonPb.KeyValuePair{
		{
			Key:   "blockHeight",
			Value: strconv.FormatInt(height, 10),
		},
	}

	payloadBytes := constructPayload(commonPb.ContractName_SYSTEM_CONTRACT_QUERY.String(), "GET_BLOCK_WITH_TXRWSETS_BY_HEIGHT", pairs)

	resp := proposalRequest(sk3, client, commonPb.TxType_QUERY_SYSTEM_CONTRACT,
		chainId, "", payloadBytes)

	blockInfo := &commonPb.BlockInfo{}
	err := proto.Unmarshal(resp.ContractResult.Result, blockInfo)
	if err != nil {
		fmt.Printf(logTempUnmarshalBlockInfoFailed, err)
		os.Exit(0)
	}
	fmt.Printf(logTempSendBlock, resp.ContractResult.Code, resp.ContractResult.Message, blockInfo)

	return hex.EncodeToString(blockInfo.Block.Header.BlockHash)
}

func testGetBlockByHash(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, chainId string, hash string) {
	fmt.Printf("\n============ get block by hash [%s] ============\n", hash)
	// 构造Payload
	pairs := []*commonPb.KeyValuePair{
		{
			Key:   "blockHash",
			Value: hash,
		},
		{
			Key:   fieldWithRWSet,
			Value: "false",
		},
	}

	payloadBytes := constructPayload(commonPb.ContractName_SYSTEM_CONTRACT_QUERY.String(), "GET_BLOCK_BY_HASH", pairs)

	resp := proposalRequest(sk3, client, commonPb.TxType_QUERY_SYSTEM_CONTRACT,
		chainId, "", payloadBytes)

	blockInfo := &commonPb.BlockInfo{}
	err := proto.Unmarshal(resp.ContractResult.Result, blockInfo)
	if err != nil {
		fmt.Printf(logTempUnmarshalBlockInfoFailed, err)
		os.Exit(0)
	}
	fmt.Printf(logTempSendBlock, resp.ContractResult.Code, resp.ContractResult.Message, blockInfo)
}

func testGetBlockWithTxRWSetsByHash(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, chainId string, hash string) {
	fmt.Printf("\n============ get block with txRWsets by hash [%s] ============\n", hash)
	// 构造Payload
	pairs := []*commonPb.KeyValuePair{
		{
			Key:   "blockHash",
			Value: hash,
		},
	}

	payloadBytes := constructPayload(commonPb.ContractName_SYSTEM_CONTRACT_QUERY.String(), "GET_BLOCK_WITH_TXRWSETS_BY_HASH", pairs)

	resp := proposalRequest(sk3, client, commonPb.TxType_QUERY_SYSTEM_CONTRACT,
		chainId, "", payloadBytes)

	blockInfo := &commonPb.BlockInfo{}
	err := proto.Unmarshal(resp.ContractResult.Result, blockInfo)
	if err != nil {
		fmt.Printf(logTempUnmarshalBlockInfoFailed, err)
		os.Exit(0)
	}
	fmt.Printf(logTempSendBlock, resp.ContractResult.Code, resp.ContractResult.Message, blockInfo)
}

func testGetLastConfigBlock(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, chainId string) {
	fmt.Printf("\n============ get last config block ============\n")
	// 构造Payload
	pairs := []*commonPb.KeyValuePair{
		{
			Key:   fieldWithRWSet,
			Value: "true",
		},
	}

	payloadBytes := constructPayload(commonPb.ContractName_SYSTEM_CONTRACT_QUERY.String(), "GET_LAST_CONFIG_BLOCK", pairs)

	resp := proposalRequest(sk3, client, commonPb.TxType_QUERY_SYSTEM_CONTRACT,
		chainId, "", payloadBytes)

	blockInfo := &commonPb.BlockInfo{}
	err := proto.Unmarshal(resp.ContractResult.Result, blockInfo)
	if err != nil {
		fmt.Printf(logTempUnmarshalBlockInfoFailed, err)
		os.Exit(0)
	}
	fmt.Printf(logTempSendBlock, resp.ContractResult.Code, resp.ContractResult.Message, blockInfo)
}

func testGetLastBlock(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, chainId string) {
	fmt.Printf("\n============ get last block ============\n")
	// 构造Payload
	pairs := []*commonPb.KeyValuePair{
		{
			Key:   fieldWithRWSet,
			Value: "true",
		},
	}

	payloadBytes := constructPayload(commonPb.ContractName_SYSTEM_CONTRACT_QUERY.String(), "GET_LAST_BLOCK", pairs)

	resp := proposalRequest(sk3, client, commonPb.TxType_QUERY_SYSTEM_CONTRACT,
		chainId, "", payloadBytes)

	blockInfo := &commonPb.BlockInfo{}
	err := proto.Unmarshal(resp.ContractResult.Result, blockInfo)
	if err != nil {
		fmt.Printf(logTempUnmarshalBlockInfoFailed, err)
		os.Exit(0)
	}
	fmt.Printf(logTempSendBlock, resp.ContractResult.Code, resp.ContractResult.Message, blockInfo)
}

func testGetChainInfo(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, chainId string) {
	fmt.Printf("\n============ get chain info ============\n")
	// 构造Payload
	pairs := []*commonPb.KeyValuePair{}

	payloadBytes := constructPayload(commonPb.ContractName_SYSTEM_CONTRACT_QUERY.String(), "GET_CHAIN_INFO", pairs)

	resp := proposalRequest(sk3, client, commonPb.TxType_QUERY_SYSTEM_CONTRACT,
		chainId, "", payloadBytes)

	chainInfo := &discoveryPb.ChainInfo{}
	err := proto.Unmarshal(resp.ContractResult.Result, chainInfo)
	if err != nil {
		fmt.Printf("chainInfo unmarshal error %s\n", err)
		os.Exit(0)
	}
	fmt.Printf(logTempSendBlock, resp.ContractResult.Code, resp.ContractResult.Message, chainInfo)
}

func testCreate(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, chainId string) {

	txId := utils.GetRandTxId()

	fmt.Printf("\n============ create contract [%s] ============\n", txId)

	//wasmBin, _ := base64.StdEncoding.DecodeString(WasmPath)
	wasmBin, _ := ioutil.ReadFile(WasmPath)
	var pairs []*commonPb.KeyValuePair

	method := commonPb.ManageUserContractFunction_INIT_CONTRACT.String()

	payload := &commonPb.ContractMgmtPayload{
		ChainId: chainId,
		ContractId: &commonPb.ContractId{
			ContractName:    contractName,
			ContractVersion: "1.0.0",
			//RuntimeType:     commonPb.RuntimeType_GASM,
			RuntimeType: runtimeType,
		},
		Method:     method,
		Parameters: pairs,
		ByteCode:   wasmBin,
	}

	if endorsement, err := acSign(payload, []int{1, 2, 3, 4}); err == nil {
		payload.Endorsement = endorsement
	} else {
		log.Fatalf("testCreate failed to sign endorsement, %s", err.Error())
		os.Exit(0)
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalf(logTempMarshalPayLoadFailed, err.Error())
		os.Exit(0)
	}

	resp := proposalRequest(sk3, client, commonPb.TxType_MANAGE_USER_CONTRACT,
		chainId, txId, payloadBytes)

	fmt.Printf(logTempSendTx, resp.Code, resp.Message, resp.ContractResult)
}

func testFreezeOrUnfreezeOrRevoke(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, chainId string, method string) {

	txId := utils.GetRandTxId()

	fmt.Printf("\n============ freeze contract [%s] ============\n", txId)

	payload := &commonPb.ContractMgmtPayload{
		ChainId: chainId,
		ContractId: &commonPb.ContractId{
			ContractName: contractName,
			RuntimeType:  commonPb.RuntimeType_WASMER,
		},
		Method: method,
	}

	if endorsement, err := acSign(payload, []int{1, 2, 3, 4}); err == nil {
		payload.Endorsement = endorsement
	} else {
		log.Fatalf("testFreezeOrUnfreezeOrRevoke failed to sign endorsement, %s", err.Error())
		os.Exit(0)
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalf(logTempMarshalPayLoadFailed, err.Error())
		os.Exit(0)
	}

	resp := proposalRequest(sk3, client, commonPb.TxType_MANAGE_USER_CONTRACT, chainId, txId, payloadBytes)

	fmt.Printf(logTempSendTx, resp.Code, resp.Message, resp.ContractResult)
}

func testUpgrade(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, chainId string) {

	txId := utils.GetRandTxId()

	fmt.Printf("\n============ create contract [%s] ============\n", txId)

	wasmBin, _ := ioutil.ReadFile(WasmPath)
	var pairs []*commonPb.KeyValuePair

	method := commonPb.ManageUserContractFunction_UPGRADE_CONTRACT.String()

	payload := &commonPb.ContractMgmtPayload{
		ChainId: chainId,
		ContractId: &commonPb.ContractId{
			ContractName:    contractName,
			ContractVersion: "2.0.0",
			//RuntimeType:     commonPb.RuntimeType_GASM,
			RuntimeType: commonPb.RuntimeType_WASMER,
		},
		Method:     method,
		Parameters: pairs,
		ByteCode:   wasmBin,
	}
	if endorsement, err := acSign(payload, []int{1, 2, 3, 4}); err == nil {
		payload.Endorsement = endorsement
	} else {
		log.Fatalf("testUpgrade failed to sign endorsement, %s", err.Error())
		os.Exit(0)
	}
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalf(logTempMarshalPayLoadFailed, err.Error())
		os.Exit(0)
	}

	resp := proposalRequest(sk3, client, commonPb.TxType_MANAGE_USER_CONTRACT,
		chainId, txId, payloadBytes)

	fmt.Printf(logTempSendTx, resp.Code, resp.Message, resp.ContractResult)
}

func testInvoke(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, chainId string) string {
	txId := utils.GetRandTxId()
	fmt.Printf("\n============ invoke contract [%s] ============\n", txId)

	// 构造Payload
	pairs := []*commonPb.KeyValuePair{
		{
			Key:   "time",
			Value: "counter1",
		},
		{
			Key:   "file_hash",
			Value: "fileHash007",
		},
		{
			Key:   "file_name",
			Value: "counter1的是",
		},
		{
			Key:   "tx_id",
			Value: "counter3",
		},
		{
			Key:   "app_id",
			Value: "1",
		},
	}
	//for i := 0; i < 1000; i++ {
	//	val := &commonPb.KeyValuePair{
	//		Key:   "key_" + strconv.Itoa(i),
	//		Value: "counter1",
	//	}
	//	pairs = append(pairs, val)
	//}
	payload := &commonPb.TransactPayload{
		ContractName: contractName,
		Method:       "save",
		Parameters:   pairs,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalf(logTempMarshalPayLoadFailed, err.Error())
	}

	resp := proposalRequest(sk3, client, commonPb.TxType_INVOKE_USER_CONTRACT,
		chainId, txId, payloadBytes)

	fmt.Printf(logTempSendTx, resp.Code, resp.Message, resp.ContractResult)

	return txId
}

func testQuery2(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, chainId string) string {
	txId := utils.GetRandTxId()
	fmt.Printf("\n============ invoke contract [%s] ============\n", txId)

	// 构造Payload
	pairs := []*commonPb.KeyValuePair{
		{
			Key:   "time",
			Value: "counter21",
		},
		{
			Key:   "file_hash",
			Value: "fileHash007",
		},
		{
			Key:   "file_name",
			Value: "counter23",
		},
		{
			Key:   "tx_id",
			Value: "counter24",
		},
	}

	payload := &commonPb.TransactPayload{
		ContractName: contractName,
		Method:       "find_by_file_hash",
		Parameters:   pairs,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalf(logTempMarshalPayLoadFailed, err.Error())
	}

	resp := proposalRequest(sk3, client, commonPb.TxType_QUERY_USER_CONTRACT,
		chainId, txId, payloadBytes)

	fmt.Printf(logTempSendTx, resp.Code, resp.Message, resp.ContractResult)

	return txId
}

func testQuery(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, chainId string) {
	txId := utils.GetRandTxId()
	fmt.Printf("\n============ query contract [%s] ============\n", txId)

	// 构造Payload
	pairs := []*commonPb.KeyValuePair{
		{
			Key:   "key",
			Value: "counter1",
		},
	}

	payload := &commonPb.QueryPayload{
		ContractName: contractName,
		Method:       "query",
		Parameters:   pairs,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalf(logTempMarshalPayLoadFailed, err.Error())
		os.Exit(0)
	}

	resp := proposalRequest(sk3, client, commonPb.TxType_QUERY_USER_CONTRACT,
		chainId, txId, payloadBytes)

	fmt.Printf(logTempSendTx, resp.Code, resp.Message, resp.ContractResult)
}

func proposalRequest(sk3 crypto.PrivateKey, client *apiPb.RpcNodeClient, txType commonPb.TxType,
	chainId, txId string, payloadBytes []byte) *commonPb.TxResponse {

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(5*time.Second)))
	defer cancel()

	if txId == "" {
		txId = utils.GetRandTxId()
	}

	file, err := ioutil.ReadFile(userCrtPath)
	if err != nil {
		panic(err)
	}

	// 构造Sender
	//pubKeyString, _ := sk3.PublicKey().String()
	sender := &acPb.SerializedMember{
		OrgId:      orgId,
		MemberInfo: file,
		IsFullCert: true,
		//MemberInfo: []byte(pubKeyString),
	}

	// 构造Header
	header := &commonPb.TxHeader{
		ChainId:        chainId,
		Sender:         sender,
		TxType:         txType,
		TxId:           txId,
		Timestamp:      time.Now().Unix(),
		ExpirationTime: 0,
	}

	req := &commonPb.TxRequest{
		Header:    header,
		Payload:   payloadBytes,
		Signature: nil,
	}

	// 拼接后，计算Hash，对hash计算签名
	rawTxBytes, err := utils.CalcUnsignedTxRequestBytes(req)
	if err != nil {
		log.Fatalf("CalcUnsignedTxRequest failed, %s", err.Error())
		os.Exit(0)
	}

	fmt.Errorf("################ %s", string(sender.MemberInfo))

	signer := getSigner(sk3, sender)
	//signBytes, err := signer.Sign("SHA256", rawTxBytes)
	signBytes, err := signer.Sign("SM3", rawTxBytes)
	if err != nil {
		log.Fatalf("sign failed, %s", err.Error())
		os.Exit(0)
	}

	req.Signature = signBytes

	result, err := (*client).SendRequest(ctx, req)

	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok && statusErr.Code() == codes.DeadlineExceeded {
			fmt.Println("WARN: client.call err: deadline")
			os.Exit(0)
		}
		fmt.Printf("ERROR: client.call err: %v\n", err)
		os.Exit(0)
	}
	return result
}

func getSigner(sk3 crypto.PrivateKey, sender *acPb.SerializedMember) protocol.SigningMember {
	skPEM, err := sk3.String()
	if err != nil {
		log.Fatalf("get sk PEM failed, %s", err.Error())
	}
	//fmt.Printf("skPEM: %s\n", skPEM)

	m, err := accesscontrol.MockAccessControl().NewMemberFromCertPem(sender.OrgId, string(sender.MemberInfo))
	if err != nil {
		panic(err)
	}

	signer, err := accesscontrol.MockAccessControl().NewSigningMember(m, skPEM, "")
	if err != nil {
		panic(err)
	}
	return signer
}

func initGRPCConnect(useTLS bool) (*grpc.ClientConn, error) {
	url := fmt.Sprintf("%s:%d", IP, Port)

	if useTLS {
		tlsClient := ca.CAClient{
			ServerName: "chainmaker.org",
			CaPaths:    caPaths,
			CertFile:   userCrtPath,
			KeyFile:    userKeyPath,
		}

		c, err := tlsClient.GetCredentialsByCA()
		if err != nil {
			log.Fatalf("GetTLSCredentialsByCA err: %v", err)
			return nil, err
		}
		return grpc.Dial(url, grpc.WithTransportCredentials(*c))
	} else {
		return grpc.Dial(url, grpc.WithInsecure())
	}
}

func constructPayload(contractName, method string, pairs []*commonPb.KeyValuePair) []byte {
	payload := &commonPb.QueryPayload{
		ContractName: contractName,
		Method:       method,
		Parameters:   pairs,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalf(logTempMarshalPayLoadFailed, err.Error())
		os.Exit(0)
	}

	return payloadBytes
}

func acSign(msg *commonPb.ContractMgmtPayload, orgIdList []int) ([]*commonPb.EndorsementEntry, error) {
	msg.Endorsement = nil
	bytes, _ := proto.Marshal(msg)

	signers := make([]protocol.SigningMember, 0)
	for _, orgId := range orgIdList {

		numStr := strconv.Itoa(orgId)
		path := fmt.Sprintf(prePathFmt, numStr) + "admin1.sign.key"
		file, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		sk, err := asym.PrivateKeyFromPEM(file, nil)
		if err != nil {
			panic(err)
		}

		userCrtPath := fmt.Sprintf(prePathFmt, numStr) + "admin1.sign.crt"
		file2, err := ioutil.ReadFile(userCrtPath)
		fmt.Println("node", orgId, "crt", string(file2))
		if err != nil {
			panic(err)
		}

		// 获取peerId
		peerId, err := helper.GetLibp2pPeerIdFromCert(file2)
		fmt.Println("node", orgId, "peerId", peerId)

		// 构造Sender
		sender1 := &acPb.SerializedMember{
			OrgId:      "wx-org" + numStr + ".chainmaker.org",
			MemberInfo: file2,
			IsFullCert: true,
		}

		signer := getSigner(sk, sender1)
		signers = append(signers, signer)
	}

	return accesscontrol.MockSignWithMultipleNodes(bytes, signers, "SHA256")
}