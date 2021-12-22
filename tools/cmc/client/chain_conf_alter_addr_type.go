/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package client

import (
	"fmt"
	"strconv"
	"strings"

	"chainmaker.org/chainmaker-go/tools/cmc/util"
	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/protocol/v2"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	sdkutils "chainmaker.org/chainmaker/sdk-go/v2/utils"
	"github.com/hokaccha/go-prettyjson"
	"github.com/spf13/cobra"
)

func alterAddrTypeCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alter-addr-type",
		Short: "alter addr type, eg. ChainMaker type, ZXL type etc.",
		Long:  "alter addr type, eg. ChainMaker type, ZXL type etc.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return alterAddrType()
		},
	}

	util.AttachFlags(cmd, flags, []string{
		flagUserSignKeyFilePath, flagUserSignCrtFilePath, flagUserTlsKeyFilePath, flagUserTlsCrtFilePath, flagOrgId,
		flagChainId, flagTimeout, flagSyncResult, flagAdminCrtFilePaths, flagAdminKeyFilePaths, flagAdminOrgIds,
	})

	util.AttachAndRequiredFlags(cmd, flags, []string{
		flagSdkConfPath, flagAddressType,
	})

	return cmd
}

func alterAddrType() error {
	var adminKeys []string
	var adminCrts []string
	var adminOrgs []string

	client, err := util.CreateChainClient(sdkConfPath, chainId, orgId, userTlsCrtFilePath, userTlsKeyFilePath,
		userSignCrtFilePath, userSignKeyFilePath)
	if err != nil {
		return err
	}
	defer client.Stop()

	if sdk.AuthTypeToStringMap[client.GetAuthType()] == protocol.PermissionedWithCert {
		if adminKeyFilePaths != "" {
			adminKeys = strings.Split(adminKeyFilePaths, ",")
		}
		if adminCrtFilePaths != "" {
			adminCrts = strings.Split(adminCrtFilePaths, ",")
		}
		if len(adminKeys) != len(adminCrts) {
			return fmt.Errorf(ADMIN_ORGID_KEY_CERT_LENGTH_NOT_EQUAL_FORMAT, len(adminKeys), len(adminCrts))
		}
	} else if sdk.AuthTypeToStringMap[client.GetAuthType()] == protocol.PermissionedWithKey {
		if adminKeyFilePaths != "" {
			adminKeys = strings.Split(adminKeyFilePaths, ",")
		}
		if adminOrgIds != "" {
			adminOrgs = strings.Split(adminOrgIds, ",")
		}
		if len(adminKeys) != len(adminOrgs) {
			return fmt.Errorf(ADMIN_ORGID_KEY_LENGTH_NOT_EQUAL_FORMAT, len(adminKeys), len(adminOrgs))
		}
	} else {
		if adminKeyFilePaths != "" {
			adminKeys = strings.Split(adminKeyFilePaths, ",")
		}
		if len(adminKeys) == 0 {
			return errAdminOrgIdKeyCertIsEmpty
		}
	}

	payload, err := client.CreateChainConfigAlterAddrTypePayload(strconv.FormatInt(int64(addressType), 10))
	if err != nil {
		return err
	}

	endorsers := make([]*common.EndorsementEntry, len(adminKeys))
	for i := range adminKeys {
		if sdk.AuthTypeToStringMap[client.GetAuthType()] == protocol.PermissionedWithCert {
			e, err := sdkutils.MakeEndorserWithPath(adminKeys[i], adminCrts[i], payload)
			if err != nil {
				return err
			}

			endorsers[i] = e
		} else if sdk.AuthTypeToStringMap[client.GetAuthType()] == protocol.PermissionedWithKey {
			e, err := sdkutils.MakePkEndorserWithPath(
				adminKeys[i],
				crypto.HashAlgoMap[client.GetHashType()],
				adminOrgs[i],
				payload,
			)
			if err != nil {
				return err
			}

			endorsers[i] = e
		} else {
			e, err := sdkutils.MakePkEndorserWithPath(
				adminKeys[i],
				crypto.HashAlgoMap[client.GetHashType()],
				"",
				payload,
			)
			if err != nil {
				return err
			}

			endorsers[i] = e
		}
	}

	// send
	resp, err := client.SendChainConfigUpdateRequest(payload, endorsers, timeout, syncResult)
	if err != nil {
		return err
	}

	err = util.CheckProposalRequestResp(resp, false)
	if err != nil {
		return err
	}

	output, err := prettyjson.Marshal(resp)
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
