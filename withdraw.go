package xpay

import (
	"context"
)

const (
	uriCreateWithdrawOrder = "/xpay/create_withdraw_order"
	uriQueryWithdrawOrder  = "/xpay/query_withdraw_order"
	uriQueryBizBalance     = "/xpay/query_biz_balance"
)

// ---------- 创建提现单 ----------

// CreateWithdrawOrderReq 创建提现单请求。
type CreateWithdrawOrderReq struct {
	WithdrawNo     string `json:"withdraw_no"`               // 提现单号，8-32 位，仅字母/数字/_/-
	WithdrawAmount string `json:"withdraw_amount,omitempty"` // 提现金额，单位元（如 0.01）；不传表示全额提现
	Env            int    `json:"env"`                       // 由 SDK 自动填充
}

// CreateWithdrawOrderResp 创建提现单响应。
type CreateWithdrawOrderResp struct {
	baseResp
	WithdrawNo   string `json:"withdraw_no"`    // 提现单号
	WxWithdrawNo string `json:"wx_withdraw_no"` // 微信侧单号
}

// CreateWithdrawOrder 创建提现单。
func (c *Client) CreateWithdrawOrder(ctx context.Context, req *CreateWithdrawOrderReq) (*CreateWithdrawOrderResp, error) {
	var out CreateWithdrawOrderResp
	if err := c.post(ctx, uriCreateWithdrawOrder, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// 提现单状态枚举（QueryWithdrawOrderResp.Status）。
const (
	WithdrawStatusProcessing = 1 // 提现中
	WithdrawStatusSuccess    = 2 // 成功
	WithdrawStatusFailed     = 3 // 失败
)

// ---------- 查询提现单 ----------

// QueryWithdrawOrderReq 查询提现单请求。
type QueryWithdrawOrderReq struct {
	WithdrawNo string `json:"withdraw_no"` // 提现单号
	Env        int    `json:"env"`         // 由 SDK 自动填充
}

// QueryWithdrawOrderResp 查询提现单响应。
type QueryWithdrawOrderResp struct {
	baseResp
	WithdrawNo               string `json:"withdraw_no"`                // 提现单号
	WxWithdrawNo             string `json:"wx_withdraw_no"`             // 微信侧单号
	Status                   int    `json:"status"`                     // 见 WithdrawStatus*
	WithdrawAmount           string `json:"withdraw_amount"`            // 提现金额，单位元
	WithdrawSuccessTimestamp string `json:"withdraw_success_timestamp"` // 成功时间，秒级时间戳字符串
	CreateTime               string `json:"create_time"`                // 创建时间
	FailReason               string `json:"fail_reason"`                // 失败原因
}

// QueryWithdrawOrder 查询提现单状态。
func (c *Client) QueryWithdrawOrder(ctx context.Context, req *QueryWithdrawOrderReq) (*QueryWithdrawOrderResp, error) {
	var out QueryWithdrawOrderResp
	if err := c.post(ctx, uriQueryWithdrawOrder, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 查询商家账户可提现余额 ----------

// BizBalance 商家可提现余额。
type BizBalance struct {
	baseResp
	BalanceAvailable struct {
		Amount       string `json:"amount"`        // 可提现余额，单位元
		CurrencyCode string `json:"currency_code"` // 币种，一般为 CNY
	} `json:"balance_available"`
}

// QueryBizBalance 查询商家账户可提现余额。
func (c *Client) QueryBizBalance(ctx context.Context) (*BizBalance, error) {
	req := struct {
		Env int `json:"env"`
	}{}
	var out BizBalance
	if err := c.post(ctx, uriQueryBizBalance, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
