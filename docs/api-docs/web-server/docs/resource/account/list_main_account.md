### 描述

- 该接口提供版本：v1.6.0+。
- 该接口所需权限：二级账号查看权限。
- 该接口功能描述：获取二级账号列表。

### URL

POST /api/v1/account/main_accounts/list

### 输入参数

| 参数名称   | 参数类型   | 必选  | 描述     |
|--------|--------|-----|--------|
| filter | object | 是   | 查询过滤条件 |
| page   | object | 是   | 分页设置   |

#### filter

| 参数名称  | 参数类型        | 必选  | 描述                                                              |
|-------|-------------|-----|-----------------------------------------------------------------|
| op    | enum string | 是   | 操作符（枚举值：and、or）。如果是and，则表示多个rule之间是且的关系；如果是or，则表示多个rule之间是或的关系。 |
| rules | array       | 是   | 过滤规则，最多设置5个rules。如果rules为空数组，op（操作符）将没有作用，代表查询全部数据。             |

#### rules[n] （详情请看 rules 表达式说明）

| 参数名称  | 参数类型        | 必选 | 描述                                          |
|-------|-------------|----|---------------------------------------------|
| field | string      | 是  | 查询条件Field名称，具体可使用的用于查询的字段及其说明请看下面 - 查询参数介绍  |
| op    | enum string | 是  | 操作符（枚举值：eq、neq、gt、gte、le、lte、in、nin、cs、cis） |
| value | 可变类型        | 是  | 查询条件Value值                                  |

##### rules 表达式说明：

##### 1. 操作符

| 操作符 | 描述                                        | 操作符的value支持的数据类型                              |
|-----|-------------------------------------------|-----------------------------------------------|
| eq  | 等于。不能为空字符串                                | boolean, numeric, string                      |
| neq | 不等。不能为空字符串                                | boolean, numeric, string                      |
| gt  | 大于                                        | numeric，时间类型为字符串（标准格式："2006-01-02T15:04:05Z"） |
| gte | 大于等于                                      | numeric，时间类型为字符串（标准格式："2006-01-02T15:04:05Z"） |
| lt  | 小于                                        | numeric，时间类型为字符串（标准格式："2006-01-02T15:04:05Z"） |
| lte | 小于等于                                      | numeric，时间类型为字符串（标准格式："2006-01-02T15:04:05Z"） |
| in  | 在给定的数组范围中。value数组中的元素最多设置100个，数组中至少有一个元素  | boolean, numeric, string                      |
| nin | 不在给定的数组范围中。value数组中的元素最多设置100个，数组中至少有一个元素 | boolean, numeric, string                      |
| cs  | 模糊查询，区分大小写                                | string                                        |
| cis | 模糊查询，不区分大小写                               | string                                        |

##### 2. 协议示例

查询 name 是 "Jim" 且 age 大于18小于30 且 servers 类型是 "api" 或者是 "web" 的数据。

```json
{
  "op": "and",
  "rules": [
    {
      "field": "name",
      "op": "eq",
      "value": "Jim"
    },
    {
      "field": "age",
      "op": "gt",
      "value": 18
    },
    {
      "field": "age",
      "op": "lt",
      "value": 30
    },
    {
      "field": "servers",
      "op": "in",
      "value": [
        "api",
        "web"
      ]
    }
  ]
}
```

#### page

| 参数名称  | 参数类型   | 必选  | 描述                                                                                                                                                  |
|-------|--------|-----|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| count | bool   | 是   | 是否返回总记录条数。 如果为true，查询结果返回总记录条数 count，但查询结果详情数据 details 为空数组，此时 start 和 limit 参数将无效，且必需设置为0。如果为false，则根据 start 和 limit 参数，返回查询结果详情数据，但总记录条数 count 为0 |
| start | uint32 | 否   | 记录开始位置，start 起始值为0                                                                                                                                  |
| limit | uint32 | 否   | 每页限制条数，最大500，不能为0                                                                                                                                   |
| sort  | string | 否   | 排序字段，返回数据将按该字段进行排序                                                                                                                                  |
| order | string | 否   | 排序顺序（枚举值：ASC、DESC）                                                                                                                                  |

#### 查询参数介绍：

| 参数名称                | 参数类型   | 描述         |
|---------------------|--------|------------|
| id                  | string | 资源ID       |
| name                | string | 名称         |
| vendor              | string | 云厂商        |
| cloud_id            | string | 云ID        |
| email               | string | 邮箱         |
| managers            | json   | 主账号管理者列表   |
| bak_managers        | json   | 主账号备份管理者列表 |
| site                | string | 站点         |
| business_type       | string | 业务类型       |
| status              | string | 状态         |
| parent_account_name | string | 父账号名称      |
| parent_account_id   | string | 父账号ID      |
| dept_id             | int    | 部门ID       |
| bk_biz_id           | int    | 业务ID       |
| op_product_id       | int    | 运营产品ID     |
| memo                | string | 备注         |
| extension           | json   | 扩展字段       |
| creator             | string | 创建者        |
| reviser             | string | 修改者        |
| created_at          | string | 创建时间       |
| updated_at          | string | 修改时间       |


接口调用者可以根据以上参数自行根据查询场景设置查询规则。


### 响应数据
```
{
    "code": 0,
    "message": "",
    "data": {
        "count": 10,
        "details": [
            {
                "id": "xxxx",                           // id
                "vendor": "aws",                        // string,云厂商
                "email": "xxxx@tencent.com",            // 邮箱
				"cloud_id": "xxx",						  // 云账号id
                "parent_account_name": "xxx",           // 所属一级账号名
                "parent_account_id": "xxxx",            // 所属一级账号id
                "site": "international",                // string,站点
        		"business_type": "internal",            // string,业务类型
                "managers": ["xxx","xxx"],        // string,负责人
                "bak_managers": ["xxx","xxx"],    // string,备份负责人
                "dept_id": 1234,                        // int,组织架构ID
                "op_product_id": 1234,                  // int,运营产品ID
                "bk_biz_id": 1312,                      // int,业务ID
                "status": "xxxx",                       // string,账号状态
                "memo": "xxxxx",                        // string,备忘录
                "creator": "xx",                        // string,创建者
                "reviser": "",                          // string,修改者
                "created_at": "",                       // string,创建时间
                "updated_at": "",                        // string,修改时间
            }
            //...
        ]
    }
}
```

### 响应参数说明

| 参数名称    | 参数类型   | 描述   |
|---------|--------|------|
| code    | int32  | 状态码  |
| message | string | 请求信息 |
| data    | array  | 响应数据 |

#### data


| 参数名称    | 参数类型  | 描述    |
|---------|-------|-------|
| count   | int32 | 总记录条数 |
| details | array | 详情数据  |

#### details

| 参数名称                | 参数类型         | 描述             |
|---------------------|--------------|----------------|
| id                  | string       | 资源ID           |
| vendor              | string       | 云厂商            |
| email               | string       | 邮箱             |
| cloud_id            | string       | 云ID            |
| parent_account_name | string       | 父账号名称          |
| parent_account_id   | string       | 父账号ID          |
| site                | string       | 站点             |
| business_type       | string       | 业务类型           |
| status              | string       | 状态             |
| managers            | string array | 主账号管理者列表       |
| bak_managers        | string array | 主账号备份管理者列表     |
| dept_id             | int          | 部门ID           |
| op_product_id       | int          | 运营产品ID         |
| bk_biz_id           | int          | 业务ID           |
| status              | string       | 状态             |
| memo                | string       | 备注             |
| creator             | string       | 创建者            |
| reviser             | string       | 修改者            |
| created_at          | string       | 创建时间           |
| updated_at          | string       | 修改时间           |
