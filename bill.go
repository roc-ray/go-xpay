package xpay

import (
	"context"
)

const (
	uriDownloadBill              = "/xpay/download_bill"
	uriStartDownloadOrder        = "/xpay/start_download_order"
	uriQueryDownloadOrder        = "/xpay/query_download_order"
	uriDownloadIosSettlementBill = "/xpay/download_ios_settlement_bill"
	uriDownloadAdverfundsOrder   = "/xpay/download_adverfunds_order"
)

// ---------- 下载小程序账单 ----------

// DownloadBillReq 下载小程序账单请求。
type DownloadBillReq struct {
	BeginDS int64 `json:"begin_ds"` // 起始日期，如 20230801
	EndDS   int64 `json:"end_ds"`   // 截止日期
}

// DownloadBillResp 下载账单响应。
type DownloadBillResp struct {
	baseResp
	URL string `json:"url"` // 账单下载地址，30 分钟有效；首次调用触发生成，需轮询
}

// DownloadBill 下载小程序账单。首次调用触发生成下载 url，可间隔轮询直至返回最终 url。
func (c *Client) DownloadBill(ctx context.Context, req *DownloadBillReq) (*DownloadBillResp, error) {
	var out DownloadBillResp
	if err := c.post(ctx, uriDownloadBill, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 发起下载订单任务 ----------

// 下载订单类型枚举（StartDownloadOrderReq.OrderType）。
const (
	DownloadOrderTypeCurrency = 1 // 代币交易
	DownloadOrderTypeGoods    = 2 // 道具直购
	DownloadOrderTypeMember   = 3 // 会员订阅
	DownloadOrderTypeRefund   = 4 // 退款订单
)

// 下载订单支付渠道枚举（StartDownloadOrderReq.PayChannel）。
const (
	DownloadPayChannelVpay = 1 // 普通虚拟支付
	DownloadPayChannelIAP  = 2 // 苹果 IAP
)

// StartDownloadOrderReq 发起下载订单任务请求。
type StartDownloadOrderReq struct {
	BeginDS      int64  `json:"begin_ds"`                // 开始日期 YYYYMMDD
	EndDS        int64  `json:"end_ds"`                  // 结束日期，与开始日期间隔 ≤31 天
	OrderType    int    `json:"order_type"`              // 见 DownloadOrderType*
	OrderInfo    string `json:"order_info,omitempty"`    // 搜索关键字（交易单号/商户单号/用户ID 模糊）
	IsProvided   *bool  `json:"is_provided,omitempty"`   // 发货状态，OrderType=2/3 时必传，默认 true
	RefundStatus int    `json:"refund_status,omitempty"` // 仅 OrderType=4 有效：0-全部 2-已退款 4-退款中 5-退款失败
	Env          int    `json:"env"`                     // 由 SDK 自动填充
	PayChannel   int    `json:"pay_channel"`             // 见 DownloadPayChannel*
}

// StartDownloadOrderResp 发起下载订单任务响应。
type StartDownloadOrderResp struct {
	baseResp
	TaskID string `json:"task_id"` // 下载任务 ID
}

// StartDownloadOrder 发起下载订单任务，用 QueryDownloadOrder 查询结果。
func (c *Client) StartDownloadOrder(ctx context.Context, req *StartDownloadOrderReq) (*StartDownloadOrderResp, error) {
	var out StartDownloadOrderResp
	if err := c.post(ctx, uriStartDownloadOrder, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// 下载任务状态枚举（QueryDownloadOrderResp.Status）。
const (
	DownloadTaskInit    = 0 // 初始化
	DownloadTaskRunning = 1 // 运行中
	DownloadTaskSuccess = 2 // 成功
	DownloadTaskFailed  = 3 // 失败
)

// QueryDownloadOrderReq 查询下载订单任务请求。
type QueryDownloadOrderReq struct {
	TaskID string `json:"task_id"` // 下载任务 ID
	Env    int    `json:"env"`     // 由 SDK 自动填充
}

// QueryDownloadOrderResp 查询下载订单任务响应。
type QueryDownloadOrderResp struct {
	baseResp
	TaskID      string `json:"task_id"`      // 下载任务 ID
	Status      int    `json:"status"`       // 见 DownloadTask*
	DownloadURL string `json:"download_url"` // 仅 Status=DownloadTaskSuccess 时有值
	ExpireAt    int64  `json:"expire_at"`    // url 过期 Unix 秒级时间戳
}

// QueryDownloadOrder 查询下载订单任务结果。
func (c *Client) QueryDownloadOrder(ctx context.Context, req *QueryDownloadOrderReq) (*QueryDownloadOrderResp, error) {
	var out QueryDownloadOrderResp
	if err := c.post(ctx, uriQueryDownloadOrder, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 下载虚拟支付 iOS 月结账单 ----------

// DownloadIosSettlementBillReq 下载 iOS 月结账单请求。
type DownloadIosSettlementBillReq struct {
	StartMonth string `json:"start_month"` // 开始月份 YYYYMM
	EndMonth   string `json:"end_month"`   // 结束月份 YYYYMM
}

// IosSettlementBill iOS 月结账单项。
type IosSettlementBill struct {
	Month   string `json:"month"`    // 月份 YYYYMM
	BillURL string `json:"bill_url"` // 账单下载链接，限时有效
}

// DownloadIosSettlementBillResp 下载 iOS 月结账单响应。
type DownloadIosSettlementBillResp struct {
	baseResp
	BillList []IosSettlementBill `json:"bill_list"` // 账单列表
}

// DownloadIosSettlementBill 查询指定月份范围内的 iOS 结算账单列表及其下载链接。
func (c *Client) DownloadIosSettlementBill(ctx context.Context, req *DownloadIosSettlementBillReq) (*DownloadIosSettlementBillResp, error) {
	var out DownloadIosSettlementBillResp
	if err := c.post(ctx, uriDownloadIosSettlementBill, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 下载广告金对应的商户订单信息 ----------

// DownloadAdverfundsOrderReq 下载广告金对应商户订单请求（无支付签名）。
type DownloadAdverfundsOrderReq struct {
	FundID string `json:"fund_id"` // 广告金发放 ID
	Env    int    `json:"env"`     // 由 SDK 自动填充（仅签名校验用）
}

// DownloadAdverfundsOrderResp 下载广告金对应商户订单响应。
type DownloadAdverfundsOrderResp struct {
	baseResp
	URL string `json:"url"` // 订单下载链接；首次调用触发生成，需轮询
}

// DownloadAdverfundsOrder 下载广告金对应的商户订单信息（仅支持通用赠送广告金）。
func (c *Client) DownloadAdverfundsOrder(ctx context.Context, req *DownloadAdverfundsOrderReq) (*DownloadAdverfundsOrderResp, error) {
	var out DownloadAdverfundsOrderResp
	if err := c.post(ctx, uriDownloadAdverfundsOrder, req, requestOptions{}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
