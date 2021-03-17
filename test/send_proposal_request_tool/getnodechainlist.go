/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	commonPb "chainmaker.org/chainmaker-go/pb/protogo/common"
	discoveryPb "chainmaker.org/chainmaker-go/pb/protogo/discovery"
	"encoding/json"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"
)

func GetNodeChainListCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getNodeChainList",
		Short: "Get node chain list",
		Long:  "Get node chain list",
		RunE: func(_ *cobra.Command, _ []string) error {
			return getNodeChainList()
		},
	}

	return cmd
}

func getNodeChainList() error {
	// 构造Payload
	pairs := []*commonPb.KeyValuePair{}

	payloadBytes, err := constructPayload(commonPb.ContractName_SYSTEM_CONTRACT_QUERY.String(), "GET_NODE_CHAIN_LIST", pairs)
	if err != nil {
		return err
	}

	resp, err = proposalRequest(sk3, client, commonPb.TxType_QUERY_SYSTEM_CONTRACT,
		"system_chain", "", payloadBytes)
	if err != nil {
		return err
	}

	chainList := &discoveryPb.ChainList{}
	if err = proto.Unmarshal(resp.ContractResult.Result, chainList); err != nil {
		return err
	}
	result := &Result{
		Code:                  resp.Code,
		Message:               resp.Message,
		ContractResultCode:    resp.ContractResult.Code,
		ContractResultMessage: resp.ContractResult.Message,
		ChainList:             chainList,
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))

	return nil
}