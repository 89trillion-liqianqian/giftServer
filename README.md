## 1.整体框架

giftServer 礼品码 redis 

1）管理后台调用 - 创建礼品码：管理员在后台填写礼品描述、可领取次数、有效期和礼包内容（例如金币、钻石），提交后返回一个8位礼品码，例如： GPA8X6TP

2）管理后台调用 - 查询礼品码信息： 填写礼品码，查询创建时间、创建人员、礼品描述、礼品内容列表（物品、数量）、可领取次数、有效期、已领取次数、领取列表（领取用户、领取时间）等信息

3）客户端调用 - 验证礼品码：用户在客户端内输入礼品码并提交，如果礼品码合法且未被领取过，调用下方奖励接口，给用户增加奖励， 加奖励成功后，返回奖励内容供客户端内展示。

## 1.目录结构

```
目录：
liqianqian@liqianqian giftServer % tree
.
├── README.md									#技术文档
├── app									# 项目启动 app
│   ├── http
│   │   └── httpServer.go									#http server
│   └── main.go									#入口
├── config									#配置文件
│   └── app.ini
├── go.mod
├── go.sum
├── internal
│   ├── ctrl									#礼品码控制器，创建、查询、验证方法
│   │   └── giftCtrl.go
│   ├── handler									#礼品码业务逻辑
│   │   ├── giftHandler.go
│   │   └── gifthandler_test.go									#单元测试
│   ├── model									#redis 	model
│   │   ├── redis.go
│   │   └── requestModel.go
│   ├── myerr									#错误返回
│   │   └── err.go
│   └── router									#路由
│       └── router.go
│   └── service
│       └── service.go
├── locust									#压测
│   ├── __pycache__
│   │   ├── load.cpython-37.pyc
│   │   └── locust.cpython-37.pyc
│   ├── load.py
│   └── report_1626931630.271092.html
├── test
│   └── test.go
└── 题三流程图.jpg

```

## 3.逻辑代码分层

|    层     | 文件夹                           | 主要职责                     | 调用关系                  | 其它说明     |
| :-------: | :------------------------------- | ---------------------------- | ------------------------- | ------------ |
|  应用层   | /app/http/httpServer.go          | http 服务器启动              | 调用路由层                | 不可同层调用 |
|  路由层   | /internal/router/router.go       | 路由转发                     | 被应用层调用，调用控制层  | 不可同层调用 |
|  控制层   | /internal/ctrl/giftCtrl,go       | 礼品码管理，创建，查询，验证 | 被路由层调用，调用handler | 不可同层调用 |
| handler层 | /internal/handler/giftHandler.go | 处理具体业务                 | 被控制层调用              | 不可同层调   |
|   model   | /internal/model                  | reids 数据处理               | 被控制层调用              |              |
| 压力测试  | Locust/load.py                   | 进行压力测试                 | 无调用关系                | 不可同层调用 |

## 4.存储设计

礼品码信息：

| 内容         | 数据库 | Key        | 类型 | 说明 |
| ------------ | ------ | ---------- | ---- | ---- |
| 礼品码       | redis  | Code       | Hash |      |
| 礼品码类型   | redis  | CodeType   | Hash |      |
| 可领取次数   | redis  | DrawCount  | Hash |      |
| 有效期时间戳 | redis  | ValidTime  | Hash |      |
| 奖品内容     | redis  | Content    | Hash |      |
| 管理员       | redis  | CreateUser | Hash |      |
| 已领取次数   | redis  | CostCount  | Hash |      |
| 指定玩家     | redis  | UserId     | Hash |      |
|              |        |            |      |      |

礼品码：限制次数类型的存储

| 内容   | 数据库 | Key       | 类型   | 说明 |      |
| ------ | ------ | --------- | ------ | ---- | ---- |
| 礼品码 | redis  | Code+type | String |      |      |
|        |        |           |        |      |      |



## 5.接口设计供客户端调用的接口

5.1管理后台调用 - 创建礼品码：管理员在后台填写礼品描述、可领取次数、有效期和礼包内容（例如金币、钻石），提交后返回一个8位礼品码，例如： GPA8X6TP

请求方法

http post 

接口地址：

127.0.0.1:8000/createGift

请求参数：

```
{
		"codeType":1,			//	礼品码类型，1-指定用户一次性消耗，2-不指定用户限制兑换次数，3-不限用户不限次数兑换
		"drawCount":19, 	// 可领取次数
		"des": "这是金币的礼品码",		// 描述
		"validTime":16345634354545,//有效期，时间戳
		"content":"{"1":1000,"2":10000}",	// 内容，json字符串，{id:count}。示例 - 金币goldId 1:数量1
		"createUser":"qq",	// 创建者
		"userId":"123456",	// 礼品码类型，1-指定用户一次性消耗的用户
}
```

json

请求响应

```
{
	"code":0,
	"msg":"管理后台-创建礼品码 ok",
	"data":"90KKHauh"						// 礼品码
}
```

响应状态码

| 状态码 | 说明           |
| ------ | -------------- |
| 0      | 创建礼品码成功 |
| 1      | 创建礼品码失败 |

5.2管理后台调用 - 查询礼品码信息： 填写礼品码，查询创建时间、创建人员、礼品描述、礼品内容列表（物品、数量）、可领取次数、有效期、已领取次数、领取列表（领取用户、领取时间）等信息

请求方法

http post 

接口地址：

127.0.0.1:8000/getGift

请求参数：

```
{
		"code":90KKHauh
}
```

json

请求响应

```
{
	"code": 0,
	"data": {
		"Code": "90KKHauh",							// 礼品码
		"CodeType":"3",	//礼品码类型，1-指定用户一次性消耗，2-不指定用户限制兑换次数，3-不限用户不限次数兑换
		"Content": "{1:1000,2:10000}",	// 内容，示例 - 金币goldId 1:数量1
		"CostCount": "0",								// 已领取次数
		"CreateUser": "qq",							// 创建着
    "Des": "des 这是金币的礼品码",		 // 描述
		"DrawCount": "19",							// 可领取次数
		"ValidTime": "16345634354545"		// 有效时间
		"historyData":"{11111:1626927791}" // {用户ID:领取时间戳}
	},
	"msg": "查询礼品码成功"
}
```

响应状态码

| 状态码 | 说明           |
| ------ | -------------- |
| 0      | 查询礼品码成功 |
| 1      | 查询礼品码失败 |

5.3验证礼品码：

请求方法

http post 

接口地址：

127.0.0.1:8000/checkCode

请求参数：

```
{
		"code":90KKHauh,			//	礼品码
		"uid":8a601a2f-e101-437a-baa4-af37783c38f7,			    //	用户ID
}
```

json

请求响应

```
{
	"code": 0,
	"msg": "验证礼品码成功",
	"changes": {1:1000,2:10000},  // 变化量，礼品内容 1-金币，2-钻石
	"balance": {1:2,2:4},  // 变化前，礼品内容 1-金币，2-钻石
	"counter": {1:1002,2:10004},  // 变化后，礼品内容 1-金币，2-钻石
}
```

响应状态码

| 状态码 | 说明           |
| ------ | -------------- |
| 0      | 验证礼品码成功 |
| 1      | 验证礼品码失败 |

## 6.第三方库

gin

```
用于api服务，go web 框架
代码： github.com/gin-gonic/gin

```

redis

```
用于礼品码数据存储
包含：hash，string 
代码："github.com/gomodule/redigo/redis"
```

## 7.如何编译执行

```
#切换主目录下
cd ./app/
#编译
go build
```

## 8.todo 

```
后续优化，连接验证
```



