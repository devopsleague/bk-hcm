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

package global

import (
	"hcm/pkg/api/core"
	dataservice "hcm/pkg/api/data-service"
	dataproto "hcm/pkg/api/data-service/cloud"
	"hcm/pkg/client/common"
	"hcm/pkg/kit"
	"hcm/pkg/rest"
)

// LoadBalancerClient is data service load balancer api client.
type LoadBalancerClient struct {
	client rest.ClientInterface
}

// NewLoadBalancerClient create a new load balancer api client.
func NewLoadBalancerClient(client rest.ClientInterface) *LoadBalancerClient {
	return &LoadBalancerClient{
		client: client,
	}
}

// ListLoadBalancer list load balancer.
func (cli *LoadBalancerClient) ListLoadBalancer(kt *kit.Kit, req *core.ListReq) (*dataproto.LbListResult, error) {
	return common.Request[core.ListReq, dataproto.LbListResult](cli.client, rest.POST, kt, req, "/load_balancers/list")
}

// BatchUpdateClbBizInfo update biz
func (cli *LoadBalancerClient) BatchUpdateClbBizInfo(kt *kit.Kit, req *dataproto.ClbBizBatchUpdateReq) error {
	return common.RequestNoResp[dataproto.ClbBizBatchUpdateReq](cli.client, rest.PATCH,
		kt, req, "/load_balancers/biz/batch/update")

}

// ListListener list listener.
func (cli *LoadBalancerClient) ListListener(kt *kit.Kit, req *core.ListReq) (
	*dataproto.ListenerListResult, error) {

	return common.Request[core.ListReq, dataproto.ListenerListResult](cli.client,
		rest.POST, kt, req, "/load_balancers/listeners/list")
}

// ListTarget list target.
func (cli *LoadBalancerClient) ListTarget(kt *kit.Kit, req *core.ListReq) (*dataproto.TargetListResult, error) {
	return common.Request[core.ListReq, dataproto.TargetListResult](
		cli.client, rest.POST, kt, req, "/load_balancers/targets/list")
}

// ListTargetGroup list target group.
func (cli *LoadBalancerClient) ListTargetGroup(kt *kit.Kit, req *core.ListReq) (*dataproto.TargetGroupListResult,
	error) {
	return common.Request[core.ListReq, dataproto.TargetGroupListResult](
		cli.client, rest.POST, kt, req, "/load_balancers/target_groups/list")
}

// BatchDelete 批量删除
func (cli *LoadBalancerClient) BatchDelete(kt *kit.Kit, req *dataproto.LoadBalancerBatchDeleteReq) error {
	return common.RequestNoResp[dataproto.LoadBalancerBatchDeleteReq](cli.client, rest.DELETE,
		kt, req, "/load_balancers/batch")
}

// ListTargetGroupListenerRel list target group listener rel.
func (cli *LoadBalancerClient) ListTargetGroupListenerRel(kt *kit.Kit, req *core.ListReq) (
	*dataproto.TargetListenerRuleRelListResult, error) {

	return common.Request[core.ListReq, dataproto.TargetListenerRuleRelListResult](
		cli.client, rest.POST, kt, req, "/target_group_listener_rels/list")
}

// DeleteTargetGroup delete target group.
func (cli *LoadBalancerClient) DeleteTargetGroup(kt *kit.Kit, req *core.ListReq) error {
	return common.RequestNoResp[core.ListReq](cli.client, rest.DELETE, kt, req, "/target_groups/batch")
}

// CreateTargetGroupListenerRel create target group listener rel.
func (cli *LoadBalancerClient) CreateTargetGroupListenerRel(kt *kit.Kit,
	req *dataproto.TargetGroupListenerRelCreateReq) (*core.BatchCreateResult, error) {

	return common.Request[dataproto.TargetGroupListenerRelCreateReq, core.BatchCreateResult](
		cli.client, rest.POST, kt, req, "/target_group_listener_rels/create")
}

// DeleteListener delete listener.
func (cli *LoadBalancerClient) DeleteListener(kt *kit.Kit, req *dataproto.LoadBalancerBatchDeleteReq) error {
	return common.RequestNoResp[dataproto.LoadBalancerBatchDeleteReq](
		cli.client, rest.DELETE, kt, req, "/listeners/batch")
}

// CreateResFlowLock create res flow lock.
func (cli *LoadBalancerClient) CreateResFlowLock(kt *kit.Kit, req *dataproto.ResFlowLockCreateReq) error {
	return common.RequestNoResp[dataproto.ResFlowLockCreateReq](
		cli.client, rest.POST, kt, req, "/res_flow_locks/create")
}

// DeleteResFlowLock delete res flow lock.
func (cli *LoadBalancerClient) DeleteResFlowLock(kt *kit.Kit, req *dataproto.ResFlowLockDeleteReq) error {
	return common.RequestNoResp[dataproto.ResFlowLockDeleteReq](
		cli.client, rest.DELETE, kt, req, "/res_flow_locks/batch")
}

// ListResFlowLock list res flow lock.
func (cli *LoadBalancerClient) ListResFlowLock(kt *kit.Kit, req *core.ListReq) (
	*dataproto.ResFlowLockListResult, error) {

	return common.Request[core.ListReq, dataproto.ResFlowLockListResult](
		cli.client, rest.POST, kt, req, "/res_flow_locks/list")
}

// BatchCreateResFlowRel batch create res flow rel.
func (cli *LoadBalancerClient) BatchCreateResFlowRel(kt *kit.Kit, req *dataproto.ResFlowRelBatchCreateReq) error {
	return common.RequestNoResp[dataproto.ResFlowRelBatchCreateReq](
		cli.client, rest.POST, kt, req, "/res_flow_rels/batch/create")
}

// BatchUpdateResFlowRel 批量更新腾讯云监听器
func (cli *LoadBalancerClient) BatchUpdateResFlowRel(kt *kit.Kit, req *dataproto.ResFlowRelBatchUpdateReq) error {
	return common.RequestNoResp[dataproto.ResFlowRelBatchUpdateReq](
		cli.client, rest.PATCH, kt, req, "/res_flow_rels/batch/update")
}

// BatchDeleteResFlowRel batch delete res flow rel.
func (cli *LoadBalancerClient) BatchDeleteResFlowRel(kt *kit.Kit, req *dataservice.BatchDeleteReq) error {
	return common.RequestNoResp[dataservice.BatchDeleteReq](cli.client, rest.DELETE, kt, req, "/res_flow_rels/batch")
}

// ListResFlowRel list res flow rel.
func (cli *LoadBalancerClient) ListResFlowRel(kt *kit.Kit, req *core.ListReq) (*dataproto.ResFlowRelListResult, error) {
	return common.Request[core.ListReq, dataproto.ResFlowRelListResult](
		cli.client, rest.POST, kt, req, "/res_flow_rels/list")
}

// ResFlowLock res flow lock.
func (cli *LoadBalancerClient) ResFlowLock(kt *kit.Kit, req *dataproto.ResFlowLockReq) error {
	return common.RequestNoResp[dataproto.ResFlowLockReq](cli.client, rest.POST, kt, req, "/res_flow_locks/lock")
}

// ResFlowUnLock res flow unlock.
func (cli *LoadBalancerClient) ResFlowUnLock(kt *kit.Kit, req *dataproto.ResFlowLockReq) error {
	return common.RequestNoResp[dataproto.ResFlowLockReq](cli.client, rest.POST, kt, req, "/res_flow_locks/unlock")
}