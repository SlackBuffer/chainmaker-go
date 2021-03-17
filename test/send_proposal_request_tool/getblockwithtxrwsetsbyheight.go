/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	commonPb "chainmaker.org/chainmaker-go/pb/protogo/common"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"
)

func GetBlockWithRWSetsByHeightCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getBlockWithRWSetsByHeight",
		Short: "Get block with RW sets by height",
		Long:  "Get block with RW sets by height",
		RunE: func(_ *cobra.Command, _ []string) error {
			return getBlockWithRWSetsByHeight()
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&height, "height", "H", -1, "specify block height")

	return cmd
}

func getBlockWithRWSetsByHeight() error {
	// 构造Payload
	pairs := []*commonPb.KeyValuePair{
		{
			Key:   "blockHeight",
			Value: strconv.Itoa(height),
		},
	}

	payloadBytes, err := constructPayload(commonPb.ContractName_SYSTEM_CONTRACT_QUERY.String(), "GET_BLOCK_WITH_TXRWSETS_BY_HEIGHT", pairs)
	if err != nil {
		return err
	}

	resp, err = proposalRequest(sk3, client, commonPb.TxType_QUERY_SYSTEM_CONTRACT,
		chainId, "", payloadBytes)
	if err != nil {
		return err
	}

	blockInfo := &commonPb.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return err
	}
	result := &Result{
		Code:                  resp.Code,
		Message:               resp.Message,
		ContractResultCode:    resp.ContractResult.Code,
		ContractResultMessage: resp.ContractResult.Message,
		BlockInfo:             blockInfo,
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))

	return nil
}