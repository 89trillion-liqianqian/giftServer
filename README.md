# giftServer
礼品码 - redis 

1.目录结构

```
目录：
liqianqian@liqianqian giftServer % pwd
/Users/liqianqian/go/src/giftServer
项目结构分析：
liqianqian@liqianqian baseLogic % tree
.
├── README.md					//技术文档
├── controller				// http api
│   └── base.go				// api 
├── go.mod
├── go.sum
├── locust						//
│   ├── __pycache__
│   │   ├── load.cpython-37.pyc
│   │   └── locust.cpython-37.pyc
│   ├── load.py				//压测脚步
│   └── report_1626931630.271092.html		//压测报告
├── main.go						//入口函数
├── model							//
│   └── requestModel.go		// model 模块
└── test
│   └── test.go				// 单元测试
└── utils
    └── tool.go       // 工具方法
5 directories, 11 files
liqianqian@liqianqian baseLogic % 


```

2。运行

```
go run main.go  
```

3.api 文档

3.1

```
1）管理后台调用 - 创建礼品码：管理员在后台填写礼品描述、可领取次数、有效期和礼包内容（例如金币、钻石），提交后返回一个8位礼品码，例如： GPA8X6TP
http post 
api: ip:port/createGift
请求体：
		"codeType":1,			//	礼品码类型，1-指定用户一次性消耗，2-不指定用户限制兑换次数，3-不限用户不限次数兑换
		"drawCount":19, 	// 可领取次数
		"des": "这是金币的礼品码",		// 描述
		"validTime":16345634354545,//有效期，时间戳
		"content":"{1:1000,2:10000}",	// 内容，json字符串，{id:count}。示例 - 金币goldId 1:数量1
		"createUser":"qq",	// 创建者
		"UserId":"aaa",	// 礼品码类型，1-指定用户一次性消耗的用户
响应体
json
{
	"code":0,
	"msg":"管理后台-创建礼品码 ok",
	"data":"90KKHauh"						// 礼品码
}
状态码
0 ：创建成功
1 ：创建失败
```

3.2

```
1）2）管理后台调用 - 查询礼品码信息： 填写礼品码，查询创建时间、创建人员、礼品描述、礼品内容列表（物品、数量）、可领取次数、有效期、已领取次数、领取列表（领取用户、领取时间）等信息
http get 
api: ip:port/getGift?code=90KKHauh
请求体
code=礼品码
响应体
json
{
	"code": 0,
	"data": {
		"Code": "K2r373SK",							// 礼品码
		"CodeType":"1",	//礼品码类型，1-指定用户一次性消耗，2-不指定用户限制兑换次数，3-不限用户不限次数兑换
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

状态码
0 ：查询成功
1 ：查询失败
```

3.3

```
客户端调用 - 验证礼品码：用户在客户端内输入礼品码并提交，如果礼品码合法且未被领取过，调用下方奖励接口，给用户增加奖励， 加奖励成功后，返回奖励内容供客户端内展示。
http post 
api: ip:port/checkCode
请求体：
		"code":90KKHauh,			//	礼品码
		"uid":111111,			    //	用户ID
响应体
json
{
	"code": 0,
	"data": "{1:1000,2:10000}",  // 礼品内容
	"msg": "验证礼品码成功"
}
状态码
0 ：验证礼品码成功
1 ：验证礼品码失败
```



