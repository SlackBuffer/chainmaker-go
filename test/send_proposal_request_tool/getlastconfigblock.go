/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"

	commonPb "chainmaker.org/chainmaker/pb-go/common"
	"chainmaker.org/chainmaker/pb-go/syscontract"

	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"
)

func GetLastConfigBlockCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getLastConfigBlock",
		Short: "Get last config block",
		Long:  "Get last config block",
		RunE: func(_ *cobra.Command, _ []string) error {
			return getLastConfigBlock()
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&withRWSets, "with-rw-sets", "w", false, "specify whether return rw sets")

	return cmd
}

func getLastConfigBlock() error {
	// 构造Payload
	w := "false"
	if withRWSets {
		w = "true"
	}
	pairs := []*commonPb.KeyValuePair{
		{
			Key:   "withRWSet",
			Value: []byte(w),
		},
	}

	payloadBytes, err := constructQueryPayload(chainId, syscontract.SystemContract_CHAIN_QUERY.String(), "GET_LAST_CONFIG_BLOCK", pairs)
	if err != nil {
		return err
	}

	resp, err = proposalRequest(sk3, client, payloadBytes)
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
	fmt.Println(result.ToJsonString())

	return nil
}
