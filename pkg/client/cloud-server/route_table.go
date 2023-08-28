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

package cloudserver

import (
	"context"
	"net/http"

	"hcm/pkg/api/core"
	routetable "hcm/pkg/api/data-service/cloud/route-table"
	"hcm/pkg/criteria/errf"
	"hcm/pkg/rest"
)

// RouteTableClient is data service subnet api client.
type RouteTableClient struct {
	client rest.ClientInterface
}

// NewRouteTable create a new subnet api client.
func NewRouteTable(client rest.ClientInterface) *RouteTableClient {
	return &RouteTableClient{
		client: client,
	}
}

// ListInRes list route table
func (v *RouteTableClient) ListInRes(ctx context.Context, h http.Header, req *core.ListReq) (
	*routetable.RouteTableListResult, error) {

	resp := new(routetable.RouteTableListResp)

	err := v.client.Post().
		WithContext(ctx).
		Body(req).
		SubResourcef("/route_tables/list").
		WithHeaders(h).
		Do().
		Into(resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != errf.OK {
		return nil, errf.New(resp.Code, resp.Message)
	}

	return resp.Data, nil
}

// ListInBiz list route tables under business
func (v *RouteTableClient) ListInBiz(ctx context.Context, h http.Header, bizID int64, req *core.ListReq) (
	*routetable.RouteTableListResult, error) {

	resp := new(routetable.RouteTableListResp)

	err := v.client.Post().
		WithContext(ctx).
		Body(req).
		SubResourcef("/bizs/%d/route_tables/list", bizID).
		WithHeaders(h).
		Do().
		Into(resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != errf.OK {
		return nil, errf.New(resp.Code, resp.Message)
	}

	return resp.Data, nil
}
