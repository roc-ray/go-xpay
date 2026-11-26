# go-vpay

微信小程序【虚拟支付】（xpay）服务端 API 的 Go SDK，零第三方依赖，仅使用标准库。

> 官方文档：<https://developers.weixin.qq.com/miniprogram/dev/server/API/VirtualPayment/>

## 特性

- 覆盖虚拟支付全部服务端 API：代币、订单/退款、账单、提现、道具、广告金、投诉处理、订阅扣款、iOS 结算等
- access_token 由外部 TokenProvider 传入（如公司内部统一 token 服务），SDK 每次请求实时获取，不做缓存
- 自动计算支付签名 `pay_sig` 与用户态签名 `signature`，请求体 `env` 字段自动填充
- 请求/响应强类型结构体，`context` 超时控制，并发安全

## 安装

```bash
go get github.com/roc-ray/go-vpay
```

## 快速开始

```go
package main

import (
	"context"
	"fmt"

	"github.com/roc-ray/go-vpay"
)

func main() {
	client, err := wxvpay.New(wxvpay.Config{
		AppKey: "your-vpay-app-key",  // 虚拟支付 AppKey（小程序后台「虚拟支付-支付基础配置」）
		Env:    wxvpay.EnvProduction, // 0-正式（默认），1-沙箱需另配 SandboxKey
		// access_token 统一由外部传入（SDK 不负责获取），如内部统一 token 服务：
		TokenProvider: func(ctx context.Context) (string, error) {
			return fetchAccessTokenFromInternalService()
		},
	})
	if err != nil {
		panic(err)
	}

	// 查询代币余额（需要用户 session_key 参与用户态签名）
	balance, err := client.QueryUserBalance(context.Background(), sessionKey, &wxvpay.QueryUserBalanceReq{
		Openid: "user-openid",
		UserIP: "1.2.3.4",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("代币余额:", balance.Balance)
}
```

### 环境说明

签名所用的 AppKey 与请求体 `env` 字段一一对应，SDK 统一由 `Config.Env` 控制：

| Env | 环境 | AppKey 来源 |
|---|---|---|
| `wxvpay.EnvProduction` (0) | 正式 | `Config.AppKey` |
| `wxvpay.EnvSandbox` (1) | 沙箱 | `Config.SandboxKey` |

## API 一览

| 分类 | 方法 | 接口 |
|---|---|---|
| 代币 | `QueryUserBalance` | /xpay/query_user_balance 查询代币余额 * |
| 代币 | `CurrencyPay` | /xpay/currency_pay 扣减代币 * |
| 代币 | `CancelCurrencyPay` | /xpay/cancel_currency_pay 代币支付退款 * |
| 代币 | `PresentCurrency` | /xpay/present_currency 代币赠送 |
| 订单 | `QueryOrder` | /xpay/query_order 查询创建的订单 |
| 订单 | `RefundOrder` | /xpay/refund_order 启动订单退款任务 |
| 订单 | `NotifyProvideGoods` | /xpay/notify_provide_goods 通知已发货完成 |
| 账单 | `DownloadBill` | /xpay/download_bill 下载小程序账单 |
| 账单 | `StartDownloadOrder` / `QueryDownloadOrder` | /xpay/start_download_order、/xpay/query_download_order |
| 账单 | `DownloadIosSettlementBill` | /xpay/download_ios_settlement_bill |
| 账单 | `DownloadAdverfundsOrder` | /xpay/download_adverfunds_order |
| 提现 | `CreateWithdrawOrder` | /xpay/create_withdraw_order 创建提现单 |
| 提现 | `QueryWithdrawOrder` | /xpay/query_withdraw_order 查询提现单 |
| 提现 | `QueryBizBalance` | /xpay/query_biz_balance 查询商家可提现余额 |
| 道具 | `StartUploadGoods` / `QueryUploadGoods` | /xpay/start_upload_goods、/xpay/query_upload_goods |
| 道具 | `StartPublishGoods` / `QueryPublishGoods` | /xpay/start_publish_goods、/xpay/query_publish_goods |
| 广告金 | `QueryTransferAccount` | /xpay/query_transfer_account 查询充值账户 |
| 广告金 | `QueryAdverFunds` | /xpay/query_adver_funds 查询发放记录 |
| 广告金 | `CreateFundsBill` | /xpay/create_funds_bill 充值广告金 |
| 广告金 | `BindTransferAccount` | /xpay/bind_transfer_accout 绑定充值账户 |
| 广告金 | `QueryFundsBill` / `QueryRecoverBill` | 查询充值/回收记录 |
| 投诉 | `GetComplaintList` / `GetComplaintDetail` | 投诉列表、详情 |
| 投诉 | `GetNegotiationHistory` | 协商历史 |
| 投诉 | `ResponseComplaint` / `CompleteComplaint` | 回复用户、完成处理 |
| 投诉 | `UploadVpFile` / `GetUploadFileSign` | 上传媒体文件、图片签名头部 |
| 订阅 | `QuerySubscribeContract` | 查询签约关系 |
| 订阅 | `SendSubscribePrePayment` | 预通知扣款 |
| 订阅 | `SubmitSubscribePayOrder` | 发起订阅扣款（仅正式环境） |
| 订阅 | `CancelSubscribeContract` | 商家解约 |
| 其他 | `QueryPunishmentReasons` | 商户被管控原因查询 |
| 其他 | `QueryOrderByIosOrder` | 苹果订单对应单号查询 |

\* 标记的接口需要用户态签名，方法需传入用户当前的 `session_key`。

## 错误处理

所有接口在微信返回 `errcode != 0` 时返回 `*wxvpay.APIError`，并提供常用错误码判定：

```go
_, err := client.PresentCurrency(ctx, &wxvpay.PresentCurrencyReq{
	Openid:  openid,
	OrderID: "present-1001",
	Amount:  100,
})
if err != nil {
	// 268490004 重复操作：赠送/代币支付/广告金充值接口返回时表示之前已成功，可视为成功
	if wxvpay.IsDuplicateOperation(err) {
		return
	}
	// 268490009 session_key 过期，应引导用户重新登录
	if wxvpay.IsSessionKeyExpired(err) {
		return
	}
	if e, ok := wxvpay.AsAPIError(err); ok {
		fmt.Println("errcode:", e.Code, "errmsg:", e.Msg)
	}
}
```

## 签名算法

SDK 内部已自动完成签名，如需自行校验消息推送（如 `xpay_goods_deliver_notify`）可使用导出的签名函数：

```go
// 支付签名 pay_sig = hex(hmac_sha256(app_key, uri + '&' + post_body))
sig := wxvpay.PaySig(appKey, "/xpay/query_user_balance", postBody)

// 用户态签名 signature = hex(hmac_sha256(session_key, post_body))
sig := wxvpay.UserSignature(sessionKey, postBody)
```

## License

[MIT](LICENSE)
