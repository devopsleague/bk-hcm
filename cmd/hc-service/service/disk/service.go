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

package disk

import (
	"fmt"

	"hcm/cmd/hc-service/service/capability"
	cloudclient "hcm/cmd/hc-service/service/cloud-adaptor"
	"hcm/pkg/criteria/enumor"
	"hcm/pkg/criteria/errf"
	"hcm/pkg/rest"
)

// InitDiskService initial the disk service
func InitDiskService(cap *capability.Capability) {
	d := &diskAdaptor{
		adaptor: cap.CloudAdaptor,
	}

	h := rest.NewHandler()

	h.Add("CreateDisks", "POST", "/vendors/{vendor}/disks/create", d.CreateDisks)

	h.Load(cap.WebService)
}

type diskAdaptor struct {
	adaptor *cloudclient.CloudAdaptorClient
}

// CreateDisks 创建云硬盘
func (da *diskAdaptor) CreateDisks(cts *rest.Contexts) (interface{}, error) {
	vendor := enumor.Vendor(cts.Request.PathParameter("vendor"))
	if err := vendor.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	switch vendor {
	case enumor.TCloud:
		return TCloudCreateDisk(da, cts)
	case enumor.HuaWei:
		return HuaWeiCreateDisk(da, cts)
	case enumor.Aws:
		return AwsCreateDisk(da, cts)
	case enumor.Azure:
		return AzureCreateDisk(da, cts)
	case enumor.Gcp:
		return GcpCreateDisk(da, cts)
	default:
		return nil, fmt.Errorf("%s does not support the creation of cloud disks", vendor)
	}
}