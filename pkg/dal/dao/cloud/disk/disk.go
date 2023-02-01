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

	"hcm/pkg/criteria/errf"
	"hcm/pkg/dal/dao"
	"hcm/pkg/dal/dao/orm"
	"hcm/pkg/dal/dao/tools"
	"hcm/pkg/dal/dao/types"
	"hcm/pkg/dal/dao/types/cloud"
	"hcm/pkg/dal/table"
	"hcm/pkg/dal/table/cloud/disk"
	"hcm/pkg/dal/table/utils"
	"hcm/pkg/kit"
	"hcm/pkg/logs"
	"hcm/pkg/runtime/filter"

	"github.com/jmoiron/sqlx"
	"hcm/pkg/api/core"
)

// DiskDao ...
type DiskDao struct {
	*dao.ObjectDaoManager
}

var _ dao.ObjectDao = new(DiskDao)

// Name 返回 Dao 描述对象的表名
func (diskDao *DiskDao) Name() table.Name {
	return disk.TableName
}

// BatchCreateWithTx 批量创建云盘数据
func (diskDao *DiskDao) BatchCreateWithTx(kt *kit.Kit, tx *sqlx.Tx, disks []*disk.DiskModel) ([]string, error) {
	if len(disks) == 0 {
		return nil, errf.New(errf.InvalidParameter, "disk model data is required")
	}

	sql := fmt.Sprintf(`INSERT INTO %s (%s)	VALUES(%s)`, diskDao.Name(), disk.DiskColumns.ColumnExpr(),
		disk.DiskColumns.ColonNameExpr(),
	)

	ids, err := diskDao.IDGen().Batch(kt, disk.TableName, len(disks))
	if err != nil {
		return nil, err
	}

	for idx, d := range disks {
		d.ID = ids[idx]
	}

	err = diskDao.Orm().Txn(tx).BulkInsert(kt.Ctx, sql, disks)
	if err != nil {
		return nil, fmt.Errorf("insert %s failed, err: %v", disk.TableName, err)
	}

	return ids, nil
}

// Update 更新云盘信息
func (diskDao *DiskDao) Update(kt *kit.Kit, filterExpr *filter.Expression, updateData *disk.DiskModel) error {
	if filterExpr == nil {
		return errf.New(errf.InvalidParameter, "filter expr is nil")
	}

	whereExpr, err := filterExpr.SQLWhereExpr(tools.DefaultSqlWhereOption)
	if err != nil {
		return err
	}

	opts := utils.NewFieldOptions().AddBlankedFields("memo").AddIgnoredFields(types.DefaultIgnoredFields...)
	setExpr, toUpdate, err := utils.RearrangeSQLDataWithOption(updateData, opts)
	if err != nil {
		return fmt.Errorf("prepare parsed sql set filter expr failed, err: %v", err)
	}

	sql := fmt.Sprintf(`UPDATE %s %s %s`, diskDao.Name(), setExpr, whereExpr)

	_, err = diskDao.Orm().AutoTxn(kt, func(txn *sqlx.Tx, opt *orm.TxnOption) (interface{}, error) {
		effected, err := diskDao.Orm().Txn(txn).Update(kt.Ctx, sql, toUpdate)
		if err != nil {
			logs.ErrorJson("update disk failed, err: %v, filter: %s, rid: %v", err, filterExpr, kt.Rid)
			return nil, err
		}

		if effected == 0 {
			logs.ErrorJson("update disk, but record not found, filter: %v, rid: %v", filterExpr, kt.Rid)
			return nil, errf.New(errf.RecordNotFound, orm.ErrRecordNotFound.Error())
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	return nil
}

// UpdateByIDWithTx 根据 ID 更新单条数据
func (diskDao *DiskDao) UpdateByIDWithTx(kt *kit.Kit, tx *sqlx.Tx, diskID string, updateData *disk.DiskModel) error {
	opts := utils.NewFieldOptions().AddBlankedFields("memo").AddIgnoredFields(types.DefaultIgnoredFields...)
	setExpr, toUpdate, err := utils.RearrangeSQLDataWithOption(updateData, opts)
	if err != nil {
		return fmt.Errorf("prepare parsed sql set filter expr failed, err: %v", err)
	}

	sql := fmt.Sprintf(`UPDATE %s %s where id = :id`, diskDao.Name(), setExpr)

	toUpdate["id"] = diskID
	_, err = diskDao.Orm().Txn(tx).Update(kt.Ctx, sql, toUpdate)
	if err != nil {
		logs.ErrorJson("update disk failed, err: %v, id: %s, rid: %v", err, diskID, kt.Rid)
		return err
	}

	return nil
}

// List 根据条件查询云盘列表
func (diskDao *DiskDao) List(kt *kit.Kit, opt *types.ListOption) (*cloud.ListDisk, error) {
	if opt == nil {
		return nil, errf.New(errf.InvalidParameter, "list disk options is nil")
	}

	if err := opt.Validate(filter.NewExprOption(filter.RuleFields(disk.DiskColumns.ColumnTypes())),
		core.DefaultPageOption); err != nil {
		return nil, err
	}

	whereOpt := tools.DefaultSqlWhereOption
	whereExpr, err := opt.Filter.SQLWhereExpr(whereOpt)
	if err != nil {
		return nil, err
	}

	if opt.Page.Count {
		// this is a count request, then do count operation only.
		sql := fmt.Sprintf(`SELECT COUNT(*) FROM %s %s`, diskDao.Name(), whereExpr)
		count, err := diskDao.Orm().Do().Count(kt.Ctx, sql)
		if err != nil {
			logs.ErrorJson("count disk failed, err: %v, filter: %s, rid: %s", err, opt.Filter, kt.Rid)
			return nil, err
		}
		return &cloud.ListDisk{Count: &count}, nil
	}
	pageExpr, err := types.PageSQLExpr(opt.Page, types.DefaultPageSQLOption)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf(`SELECT %s FROM %s %s %s`, disk.DiskColumns.FieldsNamedExpr(opt.Fields), diskDao.Name(),
		whereExpr, pageExpr)

	details := make([]*disk.DiskModel, 0)
	if err = diskDao.Orm().Do().Select(kt.Ctx, &details, sql); err != nil {
		return nil, err
	}

	result := &cloud.ListDisk{Details: details}

	return result, nil
}

// DeleteWithTx 删除云盘
func (diskDao *DiskDao) DeleteWithTx(kt *kit.Kit, tx *sqlx.Tx, filterExpr *filter.Expression) error {
	if filterExpr == nil {
		return errf.New(errf.InvalidParameter, "filter expr is required")
	}

	whereExpr, err := filterExpr.SQLWhereExpr(tools.DefaultSqlWhereOption)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf(`DELETE FROM %s %s`, diskDao.Name(), whereExpr)
	if err = diskDao.Orm().Txn(tx).Delete(kt.Ctx, sql); err != nil {
		logs.ErrorJson("delete disk failed, err: %v, filter: %s, rid: %s", err, filterExpr, kt.Rid)
		return err
	}

	return nil
}

// Count 根据条件统计云盘数量
func (diskDao *DiskDao) Count(kt *kit.Kit, opt *types.CountOption) (*cloud.CountDisk, error) {
	if opt == nil {
		return nil, errf.New(errf.InvalidParameter, "count disk options is nil")
	}

	exprOption := filter.NewExprOption(filter.RuleFields(disk.DiskColumns.ColumnTypes()))
	if err := opt.Validate(exprOption); err != nil {
		return nil, err
	}
	whereOpt := tools.DefaultSqlWhereOption
	whereExpr, err := opt.Filter.SQLWhereExpr(whereOpt)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf(`SELECT COUNT(*) FROM %s %s`, diskDao.Name(), whereExpr)
	count, err := diskDao.Orm().Do().Count(kt.Ctx, sql)
	if err != nil {
		return nil, err
	}
	return &cloud.CountDisk{Count: count}, nil
}
