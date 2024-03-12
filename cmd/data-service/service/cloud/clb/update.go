/*
 *
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 混合云管理平台 (BlueKing - Hybrid Cloud Management System) available.
 * Copyright (C) 2024 THL A29 Limited,
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

package clb

import (
	"fmt"

	"hcm/pkg/api/core"
	"hcm/pkg/api/core/cloud/clb"
	protocloud "hcm/pkg/api/data-service/cloud"
	"hcm/pkg/criteria/enumor"
	"hcm/pkg/criteria/errf"
	"hcm/pkg/dal/dao/orm"
	"hcm/pkg/dal/dao/tools"
	"hcm/pkg/dal/dao/types"
	tableclb "hcm/pkg/dal/table/cloud/clb"
	tabletype "hcm/pkg/dal/table/types"
	"hcm/pkg/kit"
	"hcm/pkg/logs"
	"hcm/pkg/rest"
	"hcm/pkg/tools/converter"
	"hcm/pkg/tools/json"
	"hcm/pkg/tools/slice"

	"github.com/jmoiron/sqlx"
)

// BatchUpdateClb 批量跟新clb信息
func (svc *clbSvc) BatchUpdateClb(cts *rest.Contexts) (any, error) {
	vendor := enumor.Vendor(cts.PathParameter("vendor").String())
	if err := vendor.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	switch vendor {
	case enumor.TCloud:
		return batchUpdateClb[clb.TCloudClbExtension](cts, svc)

	default:
		return nil, fmt.Errorf("unsupport  vendor %s", vendor)
	}

}

func batchUpdateClb[T clb.Extension](cts *rest.Contexts, svc *clbSvc) (any, error) {

	req := new(protocloud.ClbExtBatchUpdateReq[T])
	if err := cts.DecodeInto(req); err != nil {
		return nil, errf.NewFromErr(errf.DecodeRequestFailed, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	clbIds := slice.Map(req.Clbs, func(one *protocloud.ClbExtUpdateReq[T]) string { return one.ID })

	extensionMap, err := svc.listClbExt(cts.Kit, clbIds)
	if err != nil {
		return nil, err
	}

	_, err = svc.dao.Txn().AutoTxn(cts.Kit, func(txn *sqlx.Tx, opt *orm.TxnOption) (any, error) {
		for _, lb := range req.Clbs {
			update := &tableclb.LoadBalancerTable{
				Name:    lb.Name,
				BkBizID: lb.BkBizID,
				Domain:  lb.Domain,
				Status:  lb.Status,

				PrivateIPv4Addresses: lb.PrivateIPv4Addresses,
				PrivateIPv6Addresses: lb.PrivateIPv6Addresses,
				PublicIPv4Addresses:  lb.PublicIPv4Addresses,
				PublicIPv6Addresses:  lb.PublicIPv6Addresses,

				CloudCreatedTime: lb.CloudCreatedTime,
				CloudStatusTime:  lb.CloudStatusTime,
				CloudExpiredTime: lb.CloudExpiredTime,
				Memo:             lb.Memo,
				Reviser:          cts.Kit.User,
			}

			if lb.Extension != nil {
				extension, exist := extensionMap[lb.ID]
				if !exist {
					continue
				}

				merge, err := json.UpdateMerge(lb.Extension, string(extension))
				if err != nil {
					return nil, fmt.Errorf("json UpdateMerge extension failed, err: %v", err)
				}
				update.Extension = tabletype.JsonField(merge)
			}

			if err := svc.dao.LoadBalancer().UpdateByIDWithTx(cts.Kit, txn, lb.ID, update); err != nil {
				logs.Errorf("update load balancer by id failed, err: %v, id: %s, rid: %s", err, lb.ID, cts.Kit.Rid)
				return nil, fmt.Errorf("update load balancer failed, err: %v", err)
			}
		}

		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (svc *clbSvc) listClbExt(kt *kit.Kit, ids []string) (map[string]tabletype.JsonField, error) {
	opt := &types.ListOption{
		Filter: tools.ContainersExpression("id", ids),
		Page:   &core.BasePage{Limit: core.DefaultMaxPageLimit},
	}

	resp, err := svc.dao.LoadBalancer().List(kt, opt)
	if err != nil {
		return nil, err
	}

	return converter.SliceToMap(resp.Details, func(t tableclb.LoadBalancerTable) (string, tabletype.JsonField) {
		return t.ID, t.Extension
	}), nil

}

// BatchUpdateClbBizInfo 批量更新业务信息
func (svc *clbSvc) BatchUpdateClbBizInfo(cts *rest.Contexts) (any, error) {
	req := new(protocloud.ClbBizBatchUpdateReq)
	if err := cts.DecodeInto(req); err != nil {
		return nil, errf.NewFromErr(errf.DecodeRequestFailed, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	updateFilter := tools.ContainersExpression("id", req.IDs)
	updateFiled := &tableclb.LoadBalancerTable{
		BkBizID: req.BkBizID,
		Reviser: cts.Kit.User,
	}
	return nil, svc.dao.LoadBalancer().Update(cts.Kit, updateFilter, updateFiled)
}