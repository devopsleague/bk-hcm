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

package eipcvmrel

import (
	datarelproto "hcm/pkg/api/data-service/cloud"
	"hcm/pkg/criteria/errf"
	"hcm/pkg/dal/dao/orm"
	tablecloud "hcm/pkg/dal/table/cloud"
	"hcm/pkg/rest"

	"github.com/jmoiron/sqlx"
)

// BatchCreateEipCvmRels ...
func (svc *relSvc) BatchCreateEipCvmRels(cts *rest.Contexts) (interface{}, error) {
	req := new(datarelproto.EipCvmRelBatchCreateReq)
	if err := cts.DecodeInto(req); err != nil {
		return nil, errf.NewFromErr(errf.DecodeRequestFailed, err)
	}

	if err := req.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	_, err := svc.dao.Txn().AutoTxn(cts.Kit, func(txn *sqlx.Tx, opt *orm.TxnOption) (interface{}, error) {
		rels := make([]*tablecloud.EipCvmRelModel, len(req.Rels))
		for idx, relReq := range req.Rels {
			rels[idx] = &tablecloud.EipCvmRelModel{
				CvmID:   relReq.CvmID,
				EipID:   relReq.EipID,
				Creator: cts.Kit.User,
			}
		}

		return nil, svc.dao.EipCvmRel().BatchCreateWithTx(cts.Kit, txn, rels)
	})

	return nil, err
}
