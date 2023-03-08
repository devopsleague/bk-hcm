/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 混合云管理平台 (BlueKing - Hybrid Cloud Management System) available.
 * Copyright (C) 2022 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * We undertake not to change the open source license (MIT license) applicable
 *
 * to the current version of the project delivered to anyone in the future.
 */

package logics

import (
	"errors"
	"fmt"

	"hcm/cmd/hc-service/logics/sync/vpc"
	cloudclient "hcm/cmd/hc-service/service/cloud-adaptor"
	"hcm/pkg/api/core"
	cloudcore "hcm/pkg/api/core/cloud"
	protocloud "hcm/pkg/api/data-service/cloud"
	hcservice "hcm/pkg/api/hc-service"
	dataclient "hcm/pkg/client/data-service"
	"hcm/pkg/criteria/enumor"
	"hcm/pkg/criteria/validator"
	"hcm/pkg/dal/dao/tools"
	"hcm/pkg/kit"
	"hcm/pkg/logs"
	"hcm/pkg/runtime/filter"
	"hcm/pkg/tools/slice"
)

// QueryVpcIDsAndSyncOption ...
type QueryVpcIDsAndSyncOption struct {
	Vendor            enumor.Vendor `json:"vendor" validate:"required"`
	AccountID         string        `json:"account_id" validate:"required"`
	CloudVpcIDs       []string      `json:"cloud_vpc_ids" validate:"required"`
	ResourceGroupName string        `json:"resource_group_name" validate:"omitempty"`
	Region            string        `json:"region" validate:"omitempty"`
}

// Validate QueryVpcIDsAndSyncOption
func (opt *QueryVpcIDsAndSyncOption) Validate() error {
	if err := validator.Validate.Struct(opt); err != nil {
		return err
	}

	if len(opt.CloudVpcIDs) == 0 {
		return errors.New("CloudVpcIDs is required")
	}

	if len(opt.CloudVpcIDs) > int(core.DefaultMaxPageLimit) {
		return fmt.Errorf("cloudIDs should <= %d", core.DefaultMaxPageLimit)
	}

	return nil
}

// QueryVpcIDsAndSync 查询vpc，如果不存在则同步完再进行查询.
func QueryVpcIDsAndSync(kt *kit.Kit, adaptor *cloudclient.CloudAdaptorClient,
	dataCli *dataclient.Client, opt *QueryVpcIDsAndSyncOption) (map[string]string, error) {

	cloudVpcIDs := slice.Unique(opt.CloudVpcIDs)
	listReq := &core.ListReq{
		Filter: tools.ContainersExpression("cloud_id", cloudVpcIDs),
		Page:   core.DefaultBasePage,
		Fields: []string{"id", "cloud_id"},
	}
	result, err := dataCli.Global.Vpc.List(kt.Ctx, kt.Header(), listReq)
	if err != nil {
		logs.Errorf("list vpc from db failed, err: %v, cloudIDs: %v, rid: %s", err, cloudVpcIDs, kt.Rid)
		return nil, err
	}

	existVpcMap := convVpcCloudIDMap(result)

	// 如果相等，则Vpc全部同步到了db
	if len(result.Details) == len(cloudVpcIDs) {
		return existVpcMap, nil
	}

	notExistCloudID := make([]string, 0)
	for _, cloudID := range cloudVpcIDs {
		if _, exist := existVpcMap[cloudID]; !exist {
			notExistCloudID = append(notExistCloudID, cloudID)
		}
	}

	// 如果有部分vpc不存在，则触发vpc同步
	switch opt.Vendor {
	case enumor.Aws:
		syncOpt := &vpc.SyncAwsOption{
			AccountID: opt.AccountID,
			Region:    opt.Region,
			CloudIDs:  notExistCloudID,
		}
		if _, err = vpc.AwsVpcSync(kt, adaptor, dataCli, syncOpt); err != nil {
			return nil, err
		}

	case enumor.TCloud:
		syncOpt := &vpc.SyncTCloudOption{
			AccountID: opt.AccountID,
			Region:    opt.Region,
			CloudIDs:  notExistCloudID,
		}
		if _, err = vpc.TCloudVpcSync(kt, syncOpt, adaptor, dataCli); err != nil {
			return nil, err
		}

	case enumor.HuaWei:
		syncOpt := &vpc.SyncHuaWeiOption{
			AccountID: opt.AccountID,
			Region:    opt.Region,
			CloudIDs:  notExistCloudID,
		}
		if _, err = vpc.HuaWeiVpcSync(kt, syncOpt, adaptor, dataCli); err != nil {
			return nil, err
		}

	case enumor.Azure:
		syncOpt := &hcservice.AzureResourceSyncReq{
			AccountID:         opt.AccountID,
			ResourceGroupName: opt.ResourceGroupName,
			CloudIDs:          notExistCloudID,
		}
		if _, err = vpc.AzureVpcSync(kt, syncOpt, adaptor, dataCli); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unknown %s vendor", opt.Vendor)
	}

	// 同步完，二次查询
	listReq = &core.ListReq{
		Filter: tools.ContainersExpression("cloud_id", notExistCloudID),
		Page:   core.DefaultBasePage,
		Fields: []string{"id", "cloud_id"},
	}
	notExistResult, err := dataCli.Global.Vpc.List(kt.Ctx, kt.Header(), listReq)
	if err != nil {
		logs.Errorf("list vpc from db failed, err: %v, cloudIDs: %v, rid: %s", err, notExistCloudID, kt.Rid)
		return nil, err
	}

	if len(notExistResult.Details) != len(cloudVpcIDs) {
		return nil, fmt.Errorf("some vpc can not sync, cloudIDs: %v", notExistCloudID)
	}

	for cloudID, id := range convVpcCloudIDMap(notExistResult) {
		existVpcMap[cloudID] = id
	}

	return existVpcMap, nil
}

func convVpcCloudIDMap(result *protocloud.VpcListResult) map[string]string {
	m := make(map[string]string, len(result.Details))
	for _, one := range result.Details {
		m[one.CloudID] = one.ID
	}
	return m
}

type vpcMeta struct {
	CloudID string
	ID      string
}

// QueryVpcIDsAndSyncForGcp 查询vpc，如果不存在则同步完再进行查询.
func QueryVpcIDsAndSyncForGcp(kt *kit.Kit, adaptor *cloudclient.CloudAdaptorClient,
	dataCli *dataclient.Client, accountID string, selfLinks []string) (map[string]vpcMeta, error) {

	sls := slice.Unique(selfLinks)
	listReq := &core.ListReq{
		Filter: &filter.Expression{
			Op: filter.And,
			Rules: []filter.RuleFactory{
				filter.AtomRule{Field: "extension.self_link", Op: filter.JSONIn.Factory(), Value: sls},
			},
		},
		Page:   core.DefaultBasePage,
		Fields: []string{"id", "cloud_id", "extension"},
	}
	result, err := dataCli.Gcp.Vpc.ListVpcExt(kt.Ctx, kt.Header(), listReq)
	if err != nil {
		logs.Errorf("list vpc from db failed, err: %v, selfLinks: %v, rid: %s", err, sls, kt.Rid)
		return nil, err
	}

	existVpcMap := convVpcSelfLinkMap(result)

	// 如果相等，则Vpc全部同步到了db
	if len(result.Details) == len(sls) {
		return existVpcMap, nil
	}

	notExistSelfLink := make([]string, 0)
	for _, cloudID := range sls {
		if _, exist := existVpcMap[cloudID]; !exist {
			notExistSelfLink = append(notExistSelfLink, cloudID)
		}
	}

	syncOpt := &vpc.SyncGcpOption{
		AccountID: accountID,
		SelfLinks: notExistSelfLink,
	}
	if _, err = vpc.GcpVpcSync(kt, syncOpt, adaptor, dataCli); err != nil {
		return nil, err
	}

	// 同步完，二次查询
	listReq = &core.ListReq{
		Filter: &filter.Expression{
			Op: filter.And,
			Rules: []filter.RuleFactory{
				filter.AtomRule{Field: "extension.self_link", Op: filter.JSONIn.Factory(), Value: notExistSelfLink},
			},
		},
		Page:   core.DefaultBasePage,
		Fields: []string{"id", "cloud_id", "extension"},
	}
	notExistResult, err := dataCli.Gcp.Vpc.ListVpcExt(kt.Ctx, kt.Header(), listReq)
	if err != nil {
		logs.Errorf("list vpc from db failed, err: %v, cloudIDs: %v, rid: %s", err, notExistSelfLink, kt.Rid)
		return nil, err
	}

	if len(notExistResult.Details) != len(sls) {
		return nil, fmt.Errorf("some vpc can not sync, selfLinks: %v", notExistSelfLink)
	}

	for cloudID, id := range convVpcSelfLinkMap(notExistResult) {
		existVpcMap[cloudID] = id
	}

	return existVpcMap, nil
}

func convVpcSelfLinkMap(result *protocloud.VpcExtListResult[cloudcore.GcpVpcExtension]) map[string]vpcMeta {
	m := make(map[string]vpcMeta, len(result.Details))
	for _, one := range result.Details {
		m[one.Extension.SelfLink] = vpcMeta{
			CloudID: one.CloudID,
			ID:      one.ID,
		}
	}
	return m
}