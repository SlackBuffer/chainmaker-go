/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"

	commonPb "chainmaker.org/chainmaker/pb-go/common"

	"chainmaker.org/chainmaker/pb-go/syscontract"

	"chainmaker.org/chainmaker-go/utils"
	"github.com/spf13/cobra"
)

func FreezeContractCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "freezeContract",
		Short: "Freeze Contract",
		Long:  "Freeze Contract",
		RunE: func(_ *cobra.Command, _ []string) error {
			return freezeContract()
		},
	}

	return cmd
}

func freezeContract() error {
	txId := utils.GetRandTxId()

	method := syscontract.ContractManageFunction_FREEZE_CONTRACT.String()
	var pairs []*commonPb.KeyValuePair
	pairs = append(pairs, &commonPb.KeyValuePair{
		Key:   syscontract.UpgradeContract_CONTRACT_NAME.String(),
		Value: []byte(contractName),
	})
	payload, _ := constructInvokePayload(chainId, contractName, method, pairs)

	endorsement, err := acSign(payload)
	if err != nil {
		return err
	}

	resp, err := proposalRequestWithMultiSign(sk3, client, payload, endorsement)
	if err != nil {
		return err
	}

	result := &Result{
		Code:    resp.Code,
		Message: resp.Message,
		TxId:    txId,
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))

	return nil
}

func UnfreezeContractCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unfreezeContract",
		Short: "Unfreeze Contract",
		Long:  "Unfreeze Contract",
		RunE: func(_ *cobra.Command, _ []string) error {
			return unfreezeContract()
		},
	}

	return cmd
}

func unfreezeContract() error {
	txId := utils.GetRandTxId()

	method := syscontract.ContractManageFunction_UNFREEZE_CONTRACT.String()
	var pairs []*commonPb.KeyValuePair
	pairs = append(pairs, &commonPb.KeyValuePair{
		Key:   syscontract.UpgradeContract_CONTRACT_NAME.String(),
		Value: []byte(contractName),
	})
	payload, _ := constructInvokePayload(chainId, contractName, method, pairs)

	endorsement, err := acSign(payload)
	if err != nil {
		return err
	}

	resp, err := proposalRequestWithMultiSign(sk3, client, payload, endorsement)
	if err != nil {
		return err
	}

	result := &Result{
		Code:    resp.Code,
		Message: resp.Message,
		TxId:    txId,
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))

	return nil
}

func RevokeContractCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revokeContract",
		Short: "Revoke Contract",
		Long:  "Revoke Contract",
		RunE: func(_ *cobra.Command, _ []string) error {
			return RevokeContract()
		},
	}

	return cmd
}

func RevokeContract() error {

	method := syscontract.ContractManageFunction_REVOKE_CONTRACT.String()
	var pairs []*commonPb.KeyValuePair
	pairs = append(pairs, &commonPb.KeyValuePair{
		Key:   syscontract.UpgradeContract_CONTRACT_NAME.String(),
		Value: []byte(contractName),
	})
	payload, _ := constructInvokePayload(chainId, contractName, method, pairs)

	endorsement, err := acSign(payload)
	if err != nil {
		return err
	}

	resp, err := proposalRequestWithMultiSign(sk3, client, payload, endorsement)
	if err != nil {
		return err
	}

	result := &Result{
		Code:    resp.Code,
		Message: resp.Message,
		TxId:    txId,
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))

	return nil
}
